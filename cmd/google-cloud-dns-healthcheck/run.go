package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"gopkg.in/src-d/go-cli.v0"
	"gopkg.in/src-d/go-log.v1"
)

const ResourceRecordSetTypeA = "A"

func init() {
	app.AddCommand(&RunCommand{})
}

type RunCommand struct {
	cli.PlainCommand `name:"run" short-description:"run a watcher for PVs" long-description:"Run an in-cluster watcher for PVs and create the needed paths if needed"`
	RecordName       string   `short:"n" long:"record-name" required:"true" env:"RECORD_NAME" description:"Dns record name"`
	Project          string   `short:"p" long:"project" required:"true" env:"PROJECT" description:"Google project"`
	ManagedZone      string   `short:"z" long:"managed-zone" required:"true" env:"MANAGED_ZONE" description:"Google DNS managed zone"`
	HealthcheckPath  string   `short:"c" long:"healthcheck-path" required:"true" env:"HEALTHCHECK_PATH" description:"HealthcheckPath"`
	Rrdatas          []string `short:"r" long:"rrdatas" required:"true" env:"RRDATAS" env-delim:"," description:"Expected rrdatas (in comma-separated format from env variable)"`
	HttpTimeout      int64    `short:"t" long:"http-timeout" env:"HTTP_TIMEOUT" default:"5" description:"Expected rrdatas comma-separated format"`
	HttpScheme       string   `short:"s" long:"http-scheme" env:"HTTP_SCHEME" default:"http" description:"Http scheme"`
	HttpPort         string   `short:"P" long:"http-port" env:"HTTP_PORT" description:"Port for the HTTP connections"`
	DryRun           bool     `short:"d" long:"dry-run" env:"DRY_RUN" description:"Run without performing any modification"`
}

func (r *RunCommand) ExecuteContext(ctx context.Context, args []string) error {
	dnsService, err := r.getDnsService(ctx)
	if err != nil {
		log.Errorf(err, "Error getting google dns service handler")
		return err
	}

	record, err := r.getDnsRecord(dnsService)
	if err != nil {
		log.Errorf(err, "Error getting dns record")
		return err
	}

	if len(r.intersection(record.Rrdatas)) == 0 {
		err := fmt.Errorf("Wrong rrdatas value")
		log.Errorf(err, "No intersection between given rrdatas %v and record's rrdatas %v", r.Rrdatas, record.Rrdatas)
		return err
	}

	healthyRrdatas := r.checkRrdatas()

	switch len(healthyRrdatas) {
	case len(record.Rrdatas):
		log.Infof("All rrdatas are healthy")
		return nil
	case 0:
		log.Warningf("All rrdatas are unhealthy. We won't touch the record")
		return nil
	default:
		log.Infof("Updating record")
		if err := r.updateDnsRecord(dnsService, record, healthyRrdatas); err != nil {
			log.Errorf(err, "Error changing DNS RecordSet")
			return err
		}
		return nil
	}

}

func (r *RunCommand) getDnsService(ctx context.Context) (*dns.Service, error) {
	c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	dnsService, err := dns.New(c)
	if err != nil {
		return nil, err
	}

	return dnsService, nil
}

func (r *RunCommand) getDnsRecord(dnsService *dns.Service) (*dns.ResourceRecordSet, error) {
	resp, err := dnsService.ResourceRecordSets.List(r.Project, r.ManagedZone).Name(r.RecordName).Type(ResourceRecordSetTypeA).Do()

	if err != nil {
		return nil, fmt.Errorf("Issues find record %s: %s", r.RecordName, err)
	}

	if len(resp.Rrsets) > 1 {
		return nil, fmt.Errorf("Only expected 1 record set, got %d", len(resp.Rrsets))
	}

	return resp.Rrsets[0], nil
}

func (r *RunCommand) checkRrdatas() []string {
	var client = &http.Client{
		Timeout: time.Second * time.Duration(r.HttpTimeout),
	}

	healthyRrdatas := []string{}
	for _, rrdata := range r.Rrdatas {
		host := rrdata // we may add a port here
		if len(r.HttpPort) > 0 {
			host = fmt.Sprintf("%s:%s", host, r.HttpPort)
		}

		u := &url.URL{
			Scheme: r.HttpScheme,
			Host:   host,
			Path:   r.HealthcheckPath,
		}

		res, err := client.Get(u.String())
		if err != nil {
			log.Warningf("Error probing for %s", u.String())
			continue
		}

		if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusBadRequest {
			log.Infof("Success probing for %s", u.String())
			healthyRrdatas = append(healthyRrdatas, rrdata)
		} else {
			log.Warningf("Error probing for %s: %d code", u.String(), res.StatusCode)
		}

	}

	return healthyRrdatas
}

func (r *RunCommand) updateDnsRecord(dnsService *dns.Service, record *dns.ResourceRecordSet, rrdatas []string) error {
	change := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{record},
		Additions: []*dns.ResourceRecordSet{
			{
				Name:    record.Name,
				Type:    record.Type,
				Ttl:     record.Ttl,
				Rrdatas: rrdatas,
			},
		},
	}

	log.Infof("DNS Record change request for %s - old: %v new: %v", record.Name, change.Deletions[0].Rrdatas, change.Additions[0].Rrdatas)

	if !r.DryRun {
		_, err := dnsService.Changes.Create(r.Project, r.ManagedZone, change).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RunCommand) intersection(a []string) []string {
	intersection := []string{}
	hash := map[string]bool{}

	for _, rrdata := range r.Rrdatas {
		hash[rrdata] = true
	}

	for _, rrdata := range a {
		if _, found := hash[rrdata]; found {
			intersection = append(intersection, rrdata)
		}
	}

	return intersection
}
