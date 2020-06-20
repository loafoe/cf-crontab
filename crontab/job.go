package crontab

// Job describes a job
type Job struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}
