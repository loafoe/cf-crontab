package crontab

import (
	"encoding/json"
	"os"
)

const (
	maxPartSize = 16384
	configTag   = "CF_CRONTAB_CONFIG"
)

func LoadFromEnv() ([]Task, error) {
	var entries []Task

	tasks := os.Getenv(configTag)

	err := json.Unmarshal([]byte(tasks), &entries)
	if err != nil {
		return entries, err
	}
	return entries, nil
}
