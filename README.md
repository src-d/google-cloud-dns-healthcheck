# Google Cloud DNS healthcheck

This program checks the health of the IPs contained in a Google DNS record via http probes and updates the record accordingly.

```
$ google-cloud-dns-healthcheck run --help
Usage:
  google-cloud-dns-healthcheck [OPTIONS] run [run-OPTIONS]

Run an in-cluster watcher for PVs and create the needed paths if needed

Help Options:
  -h, --help                  Show this help message

[run command options]
      -n, --record-name=      Dns record name [$RECORD_NAME]
      -p, --project=          Google project [$PROJECT]
      -z, --managed-zone=     Google DNS managed zone [$MANAGED_ZONE]
      -c, --healthcheck-path= HealthcheckPath [$HEALTHCHECK_PATH]
      -r, --rrdatas=          Expected rrdatas (in comma-separated format from env variable) [$RRDATAS]
      -t, --http-timeout=     Expected rrdatas comma-separated format (default: 5) [$HTTP_TIMEOUT]
      -s, --http-scheme=      Http scheme (default: http) [$HTTP_SCHEME]
      -P, --http-port=        Port for the HTTP connections [$HTTP_PORT]
      -d, --dry-run           Run without performing any modification [$DRY_RUN]
```

Google credentials are taken from the environment. The usual setup is via `GOOGLE_APPLICATION_CREDENTIALS` env variable.

There's also a [Helm chart](https://github.com/src-d/charts/tree/master/google-cloud-dns-healthcheck) to deploy it in kubernetes.
