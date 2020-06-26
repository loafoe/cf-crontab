package crontab

import "encoding/json"

// Job describes a job
type Job struct {
	Type   string            `json:"type"`
	Command json.RawMessage `json:"command,omitempty"`
}
