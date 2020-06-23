package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/philips-labs/cf-crontab/crontab"
)

// Crontab implements the CF plugin interface
type Crontab struct {
	Entries *[]crontab.Task
}

func (c *Crontab) GetMetadata() plugin.PluginMetadata {
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
				Name:     "crontab",
				HelpText: "List all crontab entries",
				UsageDetails: plugin.Usage{
					Usage: "cf crontab",
				},
			},
			{
				Name:     "add-cron",
				HelpText: "Add a cron job",
				UsageDetails: plugin.Usage{
					Usage: "cf add-cron",
					Options: map[string]string{
						"schedule": "the cron schedule",
						"type":     "the job type [http,amqp,iron,etc]",
						"params":   "job params",
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
				HelpText: "Restore cron table",
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

func (c *Crontab) CrontabEntries() []crontab.Task {
	if c.Entries == nil {
		return []crontab.Task{}
	}
	return *c.Entries
}

func (c *Crontab) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "crontab":
		server, err := CrontabServerResolver(cliConnection)
		if err != nil {
			fmt.Printf("error resolving server: %v\n", err)
			return
		}
		host, err := server.Host()
		if err != nil {
			fmt.Printf("error resolving host: %v\n", err)
			return
		}
		fmt.Printf("Host: %v\n", host)
		entries, err := server.GetEntries()
		if err != nil {
			fmt.Printf("error resolving host: %v\n", err)
			return
		}
		fmt.Printf("entries: %v\n", entries)
	case "add-cron":
		if len(args) < 2 {
			fmt.Printf("need json file with Entries\n")
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
		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		for i, t := range c.CrontabEntries() {
			if int(t.EntryID) == index {
				fmt.Printf("Removing %d\n", index)
				*c.Entries = append((*c.Entries)[:i], (*c.Entries)[i+1:]...)
				// TODO: remove on server
				// Local: (*c.Cron).Remove(cron.EntryID(index))
			}
		}
	case "backup-cron":
		data, err := json.Marshal(c.Entries)
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
		parts, err := crontab.EnvParts(tasks)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		for k, v := range parts {
			fmt.Printf("  %s: %s\n", k, v)
		}
	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using cf-crontab")
	}
}
