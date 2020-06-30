package crontab

import (
	"fmt"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/philips-software/gautocloud-connectors/hsdp"
	"github.com/philips-software/go-hsdp-api/iron"
)

type Iron struct {
	Command  string `json:"command"`
	CodeName string `json:"code_name"`
	Cluster  string `json:"cluster"`
	Timeout  int    `json:"timeout"`
	Payload  string `json:"payload"`
	Task     *Task  `json:"-"`
}

func (i Iron) Run() {
	var client *hsdp.IronClient
	err := gautocloud.Inject(&client)
	if err != nil {
		fmt.Printf("no iron service found. please bind one to cf-crontab\n")
		return
	}
	payload := ""
	for _, cluster := range client.Config.ClusterInfo {
		if cluster.ClusterID == i.Cluster {
			var err error
			payload, err = cluster.Encrypt([]byte(i.Payload))
			if err != nil {
				fmt.Printf("failed to encrypt payload for cluster %s: %v\n", i.Cluster, err)
				return
			}
			break
		}
	}
	if payload == "" {
		fmt.Printf("cluster not found: %s", i.Cluster)
		return
	}
	switch i.Command {
	case "queue":
		task, _, err := client.Tasks.QueueTask(iron.Task{
			CodeName: i.CodeName,
			Cluster:  i.Cluster,
			Payload:  payload,
			Timeout:  i.Timeout,
		})
		if err != nil {
			fmt.Printf("error queuing iron command `%v`: %v\n", i.Command, err)
			return
		}
		if task != nil {
			fmt.Printf("queued iron task %v\n", task.ID)
		}
	default:
		fmt.Printf("command `%v` is not supported\n", i.Command)
	}
	fmt.Printf("not implemented\n")
}
