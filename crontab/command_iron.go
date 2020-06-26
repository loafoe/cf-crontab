package crontab

import (
	"fmt"
)

type Iron struct {
	Command  string `json:"command"`
	CodeName string `json:"code_name"`
	Cluster  string `json:"cluster"`
	Timeout  int `json:"timeout"`
	Payload  string `json:"payload"`
}

func (i Iron) Run() {

	fmt.Printf("not implemented\n")
}
