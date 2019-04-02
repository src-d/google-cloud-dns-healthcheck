package main

import "gopkg.in/src-d/go-cli.v0"

var (
	version string
	build   string
)

var app = cli.New("google-cloud-dns-healthcheck", version, build,
	"Check health of endpoints related to a Google Cloud DNS record and update it accordingly")

func main() {
	app.RunMain()
}
