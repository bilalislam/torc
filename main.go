package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bilalislam/torc/log"
	rabbit "github.com/bilalislam/torc/rabbitmq"
	"github.com/bilalislam/torc/storage/clients"
	"github.com/bilalislam/torc/storage/models"
	"github.com/bilalislam/torc/utils"
	"github.com/olivere/elastic/v7"
	"strings"
	"time"
)

type AnakinCommand struct {
	models.IModel `json:"-"`
	Message       string    `json:"message"`
	Title         string    `json:"title"`
	State         State     `json:"state"`
	CorrelationId string    `json:"correlationId"`
	EventOn       time.Time `json:"eventOn"`
}

const (
	StateOk       State = "ok"
	StatePaused   State = "paused"
	StateAlerting State = "alerting"
	StatePending  State = "pending"
	StateNoData   State = "no_data"
)

type State string

var (
	nodes = flag.String("nodes", "http://10.2.24.85:9200,http://10.2.24.188:9200,http://10.2.24.189:9200", "comma-separated list of ES URLs (e.g. 'http://10.2.24.85:9200,http://10.2.24.188:9200,http://10.2.24.189:9200')")
)

func main() {

	logger := log.GetLogger()
	measurement := utils.NewTimeMeasurement(logger)

	var options []elastic.ClientOptionFunc
	urls := strings.SplitN(*nodes, ",", -1)
	options = append(options, elastic.SetURL(urls...))
	client, err := elastic.NewClient(options...)
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}

	repository := clients.ElasticsearchRepository{
		Client: client,
	}

	onConsumed := func(message rabbit.Message) error {
		defer measurement.TimeTrack(time.Now(), "[Anakin Hub Consumer]", "Consumed")
		ctx := context.Background()

		var consumeMessage AnakinCommand
		var err = json.Unmarshal(message.Payload, &consumeMessage)
		if err != nil {
			return err
		}

		err = repository.Save(ctx, consumeMessage)
		if err != nil {
			return err
		}

		fmt.Println(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), " Message:", consumeMessage)
		return nil
	}

	var r = rabbit.NewRabbitMqClient([]string{"127.0.0.1"},
		"guest",
		"guest",
		"/",
		rabbit.RetryCount(3, time.Duration(0)),
		rabbit.PrefetchCount(50))

	r.AddConsumer("anakin").
		HandleConsumer(onConsumed).
		SubscriberExchange("anakin-development", rabbit.ExchangeType(3), "noc-tools")

	_ = r.RunConsumers()

}
