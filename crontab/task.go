package crontab

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
)

type Task struct {
	Schedule string       `json:"schedule"`
	Job      Job          `json:"job"`
	EntryID  cron.EntryID `json:"entryID,omitempty"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%d: %s %v", t.EntryID, t.Schedule, t.Job)
}

func (t *Task) Add(cr *cron.Cron) error {
	var job cron.Job
	switch t.Job.Type {
	case "http":
		var command Http
		err := json.Unmarshal(t.Job.Command, &command)
		if err != nil {
			return err
		}
		job = command
	case "amqp":
		var command Amqp
		err := json.Unmarshal(t.Job.Command, &command)
		if err != nil {
			return err
		}
		job = command
	default:
		return fmt.Errorf("unsupported type: %s", t.Job.Type)
	}
	entryID, err := cr.AddFunc(t.Schedule, job.Run)
	t.EntryID = entryID
	return err
}
