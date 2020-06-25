package crontab

import (
	"fmt"
)

type Cartel struct {
	Command string
	Tags string
	Instance string
	Action string
}

func (c Cartel) Run() {
	fmt.Printf("not implemented\n")
}
