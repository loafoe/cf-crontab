package main

import (
	"github.com/philips-labs/cf-crontab/server"
)

var GitCommit = "deadbeaf"

func main() {
	_ = GitCommit
	serverMode()
}

func serverMode() {
	server.Start()
}
