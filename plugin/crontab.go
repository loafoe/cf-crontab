package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/jedib0t/go-pretty/text"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/jedib0t/go-pretty/table"
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
				HelpText: "Remove a cron job",
				UsageDetails: plugin.Usage{
					Usage: "cf remove-cron",
					Options: map[string]string{
						"index": "the index of the cron entry to remove",
					},
				},
			},
			{
				Name:     "save-crontab",
				HelpText: "Save crontab table to the environment",
				UsageDetails: plugin.Usage{
					Usage: "cf save-crontab",
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

func (c *Crontab) ServerEntries(cliConnection plugin.CliConnection) ([]*crontab.Task, error) {
	fmt.Printf("Discovering crontab server ...\n")
	server, err := CrontabServerResolver(cliConnection)
	if err != nil {
		return nil, err
	}
	host, err := server.Host()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Getting crontab from %s ...\n", host)
	entries, err := server.GetEntries()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Crontab) Run(cliConnection plugin.CliConnection, args []string) {
	switch args[0] {
	case "crontab":
		entries, err := c.ServerEntries(cliConnection)
		if err != nil {
			fmt.Printf("error getting server entries: %v\n", err)
			return
		}
		c.RenderEntries(entries)
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
		if len(tasks) == 0 {
			fmt.Printf("no tasks found. check your .json\n")
			return
		}
		fmt.Printf("Adding %d entries ...\n", len(tasks))
		server, err := CrontabServerResolver(cliConnection)
		if err != nil {
			fmt.Printf("error resolving server: %v\n", err)
			return
		}
		entries, err := server.AddEntries(tasks)
		if err != nil {
			fmt.Printf("error adding entries: %v\n", err)
			return
		}
		c.RenderEntries(entries)
		fmt.Printf("OK\n")
	case "remove-cron":
		if len(args) < 2 {
			fmt.Printf("need entryID\n")
			return
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Printf("Discovering crontab server ...\n")
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
		fmt.Printf("Deleting entry %d from %s ...\n", index, host)
		ok, err := server.DeleteEntry(index)
		if err != nil {
			fmt.Printf("error deleting: %v\n", err)
			return
		}
		if !ok {
			fmt.Printf("FAILED\n")
			return
		}
		fmt.Printf("OK\n")
	case "save-crontab":
		fmt.Printf("Saving crontab ...\n")
		server, err := CrontabServerResolver(cliConnection)
		if err != nil {
			fmt.Printf("error resolving server: %v\n", err)
			return
		}
		if err := server.SaveCrontab(); err != nil {
			fmt.Printf("error saving crontab: %v\n", err)
			return
		}
		fmt.Printf("OK\n")
	case "CLI-MESSAGE-UNINSTALL":
		fmt.Println("Thanks for using cf-crontab")
	}
}

func (c *Crontab) RenderEntries(entries []*crontab.Task) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"#", "schedule", "type", "details"})
	for _, e := range entries {
		details := fmt.Sprintf("%v", e.Job)
		t.AppendRow([]interface{}{e.EntryID, e.Schedule, e.Job.Type, details})
	}
	t.Render()
}
