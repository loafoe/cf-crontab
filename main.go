package main

import (
	"fmt"
	"github.com/philips-labs/cf-crontab/config"
	"github.com/robfig/cron"
	"code.cloudfoundry.org/cli/plugin"
)

type GronPlugin struct {}

func (c *GronPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-crontab",
		Commands: []plugin.Command{
			{
				Name:     "list-grontab",
				HelpText: "List all tabs",
				UsageDetails: plugin.Usage{
					Usage: "cf list-grontab",
				},
			},
		},
	}
}

func (c *GronPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "list-grontab":
		fmt.Println("TODO")
	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using crontab")
	}
}

func main() {
	c := cron.New()

	tasks, err := config.LoadFromEnv()
	if err != nil {
		return
	}
	for _, task := range tasks {
		task.Add(c)
	}

	c.Start()

	plugin.Start(new(GronPlugin))
}