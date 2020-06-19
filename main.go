package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"encoding/json"
	"fmt"
	"github.com/philips-labs/cf-crontab/config"
	"github.com/philips-labs/cf-crontab/crontab"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type CrontabPlugin struct {}

func (c *CrontabPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Version: struct {
			Major int
			Minor int
			Build int
		}{Major: 0, Minor: 0, Build: 1},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 51,
			Build: 0,
		},
		Name: "cf-crontab",
		Commands: []plugin.Command{
			{
				Name:     "list-cron",
				HelpText: "List all tabs",
				UsageDetails: plugin.Usage{
					Usage: "cf list-cron",
				},
			},
			{
				Name:     "add-cron",
				HelpText: "Add a cron job",
				UsageDetails: plugin.Usage{
					Usage: "cf add-cron",
					Options: map[string]string{
						"schedule": "the cron schedule",
						"type": "the job type [http,amqp,iron,etc]",
						"params": "job params",
					},
				},
			},
			{
				Name:     "remove-cron",
				HelpText: "Add a cron job",
				UsageDetails: plugin.Usage{
					Usage: "cf remove-cron",
					Options: map[string]string{
						"index": "the index of the cron entry to remove",
					},
				},
			},
			{
				Name:     "backup-cron",
				HelpText: "Backup cron table",
				UsageDetails: plugin.Usage{
					Usage: "cf backup-cron",
				},
			},
			{
				Name:     "restore-cron",
				HelpText: "Resotre cron table",
				UsageDetails: plugin.Usage{
					Usage: "cf restore-cron",
					Options: map[string]string{
						"file": "the file to restore from",
					},
				},
			},
		},
	}
}

func (c *CrontabPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "list-cron":
		tasks, err := config.LoadFromEnv()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		for i, task := range tasks {
			fmt.Printf("%d: %v\n", i, task)
		}
	case "add-cron":
		if len(args) < 2 {
			fmt.Printf("need json file with tasks\n")
			return
		}
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		var tasks []crontab.Task
		err = json.Unmarshal(data, &tasks)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		for i, task := range tasks {
			fmt.Printf("%d: %v\n", i, task)
		}
	case "remove-cron":
		fmt.Println("TODO")
	case "backup-cron":
		tasks, err := config.LoadFromEnv()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		data, err := json.Marshal(&tasks)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		_ = ioutil.WriteFile(args[1], data, 0644)
	case "restore-cron":
		restore, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		var tasks []crontab.Task
		err = json.Unmarshal(restore, &tasks)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		parts, err := config.EnvParts(tasks)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		for k,v := range parts {
			fmt.Printf("  %s: %s\n", k, v)
		}
	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using crontab")
	}
}

func main() {
	log.SetOutput(ioutil.Discard)

	c := cron.New()

	tasks, err := config.LoadFromEnv()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}
	for _, task := range tasks {
		task.Add(c)
	}

	c.Start()

	plugin.Start(&CrontabPlugin{})
}