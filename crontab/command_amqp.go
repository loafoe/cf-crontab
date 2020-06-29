package crontab

import (
	"github.com/loafoe/go-rabbitmq"
	"github.com/streadway/amqp"
)

type Amqp struct {
	Exchange     string `json:"exchange"`
	ExchangeType string `json:"exchange_type"`
	RoutingKey   string `json:"routing_key"`
	Payload      string `json:"payload"`
	ContentType  string `json:"content_type"`
	Instance     string `json:"instance"`
	Task         *Task  `json:"-"`
}

func (a Amqp) Run() {
	producer, err := rabbitmq.NewProducer(rabbitmq.Config{
		Exchange:     a.Exchange,
		ExchangeType: a.ExchangeType,
		Durable:      false,
	})
	if err != nil {
		return
	}
	_ = producer.Publish(a.Exchange, a.RoutingKey, amqp.Publishing{
		ContentType: a.ContentType,
		Body:        []byte(a.Payload),
	})
}
