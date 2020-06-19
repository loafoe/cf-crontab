package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

type Task struct {
	Schedule string `json:"schedule"`
	Job Job `json:"job"`
	EntryID cron.EntryID `json:"-"`
}

func (t *Task)String() string {
	return fmt.Sprintf("%d: %s %v %v %v", t.EntryID, t.Schedule, t.Job.Type, t.Job.Params["method"], t.Job.Params["url"])
}

func (t *Task)Add(cr *cron.Cron) error {
	var command cron.Job
	switch  t.Job.Type {
	case "http":
		command = Http{
			Method: t.Job.Params["method"],
			URL: t.Job.Params["url"],
			Body: t.Job.Params["body"],
		}
	case "amqp":
		command = Amqp{
			Exchange: t.Job.Params["exchange"],
			ExchangeType: t.Job.Params["exchange_type"],
			RoutingKey:  t.Job.Params["routing_key"],
			ContentType: t.Job.Params["content_type"],
			Payload: t.Job.Params["payload"],
		}
	default:
		return fmt.Errorf("unsupported type: %s", t.Job.Type)
	}
	entryID, err := cr.AddFunc(t.Schedule, command.Run)
	t.EntryID = entryID
	return err
}