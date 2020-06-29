package crontab

import (
	"fmt"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/philips-software/go-hsdp-api/cartel"
)

type Cartel struct {
	Command  string `json:"command"`
	Name     string `json:"name"`
	Tags     string `json:"tags"`
	Instance string `json:"instance"`
	Action   string `json:"action"`
	Task     *Task  `json:"-"`
}

func (c Cartel) Run() {
	var client *cartel.Client

	err := gautocloud.Inject(&client)
	if err != nil {
		fmt.Printf("no cartel service found. please bind one to cf-crontab\n")
		return
	}
	switch c.Command {
	case "start":
		resp, _, err := client.Start(c.Name)
		if err != nil {
			fmt.Printf("error starting %s: %v\n", c.Name, err)
			return
		}
		if resp != nil && resp.Success() {
			fmt.Printf("cartel instance %s started", c.Name)
		}
	case "stop":
		resp, _, err := client.Stop(c.Name)
		if err != nil {
			fmt.Printf("error stopping %s: %v\n", c.Name, err)
			return
		}
		if resp != nil && resp.Success() {
			fmt.Printf("cartel instance %s stopped", c.Name)
		}
	default:
		fmt.Printf("command `%s` not supported\n", c.Command)
	}
}
