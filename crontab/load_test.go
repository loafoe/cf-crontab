package crontab

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvPartsAndLoad(t *testing.T) {
	var cmd = Http{
		Method: "POST",
		URL: "https://foo.com",
	}
	raw, err := json.Marshal(&cmd)
	if !assert.Nil(t, err) {
		return
	}
	task := Task{
		Schedule: "0 * * * * *",
		Job: Job{
			Type: "http",
			Command: raw,
		},
	}
	entries := make([]*Task, 0)
	for i := 0; i < 100; i++ {
		entries = append(entries, &task)
	}

	parts, err := EnvParts(entries)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.Equal(t, 3, len(parts)) {
		return
	}
	keys := make([]string, 0)
	for k, v := range parts {
		keys = append(keys, k)
		_ = os.Setenv(k, v)
	}
	assert.Equal(t,  configTag+"_0", keys[0])
	loadedParts, err := LoadFromEnv()
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 100, len(loadedParts))
}
