package crontab

import (
	"fmt"

	"net/http"
)

type Http struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Task    *Task             `json:"-"`
}

func (h Http) Run() {
	req, err := http.NewRequest(h.Method, h.URL, nil)
	if err != nil {
		return
	}
	for h, v := range h.Headers {
		req.Header.Set(h, v)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%d: %v\n", h.Task.EntryID, err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		fmt.Printf("%d: HTTP %d\n", h.Task.EntryID, resp.StatusCode)
	}
}
