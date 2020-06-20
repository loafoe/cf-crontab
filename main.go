package main

import (
	"fmt"
	"github.com/philips-labs/cf-crontab/plugin"

	cfplugin "code.cloudfoundry.org/cli/plugin"
	"github.com/philips-labs/cf-crontab/config"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()

	tasks, err := config.LoadFromEnv()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}
	for i, _ := range tasks {
		_ = tasks[i].Add(c)
	}

	c.Start()

	cfplugin.Start(&plugin.Crontab{
		Cron:    c,
		Tasks: &tasks,
	})
}
