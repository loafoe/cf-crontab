package crontab

import (
	"fmt"
)

type Cartel struct {
	Command  string `json:"command""`
	Tags     string `json:"tags"`
	Instance string `json:"instance"`
	Action   string `json:"action"`
	Task     *Task  `json:"-"`
}

func (c Cartel) Run() {
	fmt.Printf("not implemented\n")
}
