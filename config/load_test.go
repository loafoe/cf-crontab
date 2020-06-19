package config

import (
	"github.com/philips-labs/cf-crontab/crontab"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnvPartsAndLoad(t *testing.T) {
	task := crontab.Task{
		Schedule: "* * * * *",
		Job: crontab.Job{
			Type: "http",
			Params: map[string]string {
				"method": "POST",
				"foo": "BAR",
			},
		},
	}
	tasks := make([]crontab.Task, 0)
	for i :=0; i < 100; i++ {
		tasks = append(tasks, task)
	}

	parts, err := EnvParts(tasks)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.Equal(t, 3,len(parts)) {
		return
	}
	keys := make([]string, 0)
	for k, v := range parts {
		keys = append(keys, k)
		_ = os.Setenv(k, v)
	}
	assert.Equal(t, "CRONTAB_CONFIG_0", keys[0])
	loadedParts, err := LoadFromEnv()
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, 100, len(loadedParts))
}
