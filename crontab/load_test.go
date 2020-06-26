package crontab

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvPartsAndLoad(t *testing.T) {
	task := Task{
		Schedule: "0 * * * * *",
		Job: Job{
			Type: "http",
			Params: map[string]string{
				"method": "POST",
				"foo":    "BAR",
			},
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
	if !assert.Equal(t, 2, len(parts)) {
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
