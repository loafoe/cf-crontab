package crontab

import (
	"bytes"
	"net/http"
)

type Http struct {
	Method string
	URL string
	Headers http.Header
	Body string
}

func (h Http) Run() {
	req, err := http.NewRequest(h.Method, h.URL, bytes.NewBufferString(h.Body))
	if err != nil {
		return
	}
	req.Header = h.Headers
	client := http.DefaultClient
	_, _ = client.Do(req)
}