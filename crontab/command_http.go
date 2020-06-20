package crontab

import (
	"bytes"
	"net/http"
)

type Http struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (h Http) Run() {
	req, err := http.NewRequest(h.Method, h.URL, bytes.NewBufferString(h.Body))
	if err != nil {
		return
	}
	for h, v := range h.Headers {
		req.Header.Set(h, v)
	}
	client := http.DefaultClient
	_, _ = client.Do(req)
}
