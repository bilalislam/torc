package consumer

import (
	rabbit "github.com/bilalislam/torc/rabbitmq"
	"time"
)

type Request struct {
	Uri           []string
	UserName      string
	Password      string
	Exchange      string
	ExchangeType  int
	Queue         string
	RoutingKey    string
	RetryCount    int
	PrefetchCount int
}

func AddConsumer(request Request) (*rabbit.MessageBrokerServer, *rabbit.Consumer) {
	var rabbitClient = rabbit.NewRabbitMqClient(request.Uri,
		request.UserName,
		request.Password,
		"",
		rabbit.RetryCount(request.RetryCount, time.Duration(0)),
		rabbit.PrefetchCount(request.PrefetchCount))

	return rabbitClient, rabbitClient.AddConsumer(request.Queue).
		SubscriberExchange(request.RoutingKey, rabbit.ExchangeType(request.ExchangeType), request.Exchange)
}
