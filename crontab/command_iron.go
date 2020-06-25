package crontab

import "fmt"

type Iron struct {
	Command string
	CodeName string
	Cluster string
	Timeout int
	Payload string
}

func (i Iron) Run() {
	fmt.Printf("not implemented\n")
}
