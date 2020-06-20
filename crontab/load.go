package crontab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

const (
	maxPartSize = 4096
	configTag   = "CRONTAB_CONFIG"
)

func LoadFromEnv() ([]Task, error) {
	entries := make([]Task, 0)

	parts := make(chan string)
	go func(p chan string) {
		count := 0
		for {
			nextPart := fmt.Sprintf("%s_%d", configTag, count)
			if part := os.Getenv(nextPart); part != "" {
				p <- part
			} else {
				close(p)
				return
			}
			count++
		}
	}(parts)
	base64FullConfig := ""
	for p := range parts {
		base64FullConfig = base64FullConfig + p
	}
	if base64FullConfig == "" {
		return entries, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(base64FullConfig)
	if err != nil {
		return entries, err
	}
	err = json.Unmarshal(decoded, &entries)
	if err != nil {
		return entries, err
	}
	return entries, nil
}

func EnvParts(tasks []Task) (map[string]string, error) {
	parts := make(map[string]string)
	data, err := json.Marshal(&tasks)
	if err != nil {
		return parts, err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	part := 0
	l := len(encoded)
	for c := 0; c < l; c++ {
		if c > 0 && c%maxPartSize == 0 {
			partName := fmt.Sprintf("%s_%d", configTag, part)
			part++
			parts[partName] = encoded[c-maxPartSize : c]
		}
	}
	if r := l % maxPartSize; r > 0 {
		partName := fmt.Sprintf("%s_%d", configTag, part)
		parts[partName] = encoded[l-r : l]
	}
	return parts, nil
}
