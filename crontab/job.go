package crontab

import (
	"encoding/json"
	"fmt"
)

// Job describes a job
type Job struct {
	Type    string          `json:"type"`
	Command json.RawMessage `json:"command,omitempty"`
}

func (j Job) String() string {
	switch j.Type {
	case "http":
		var cmd Http
		err := json.Unmarshal(j.Command, &cmd)
		if err != nil {
			return "http: unknown or invalid commands"
		}
		return fmt.Sprintf("%s %s", cmd.Method, cmd.URL)
	case "iron":
		var cmd Iron
		err := json.Unmarshal(j.Command, &cmd)
		if err != nil {
			return "iron: unknown or invalid commands"
		}
		return fmt.Sprintf("%s %s with timeout %d", cmd.Command, cmd.CodeName, cmd.Timeout)
	default:
		return "..."
	}
}
