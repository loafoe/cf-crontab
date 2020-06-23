package main

import (
	"os"

	"github.com/philips-labs/cf-crontab/plugin"
	"github.com/philips-labs/cf-crontab/server"

	cfplugin "code.cloudfoundry.org/cli/plugin"
)

var GitCommit = "deadbeaf"

func main() {
	_ = GitCommit
	if len(os.Args) == 2 && os.Args[1] == "server" {
		serverMode()
		return
	}
	pluginMode()
}

func pluginMode() {
	cfplugin.Start(&plugin.Crontab{})
}

func serverMode() {
	server.Start()
}
