package worker

import (
	"fmt"
	"os"

	"github.com/dwnGnL/testWork/internal/application"
	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/lib/goerrors"

	amqp "github.com/rabbitmq/amqp091-go"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

const (
	consumerNameTemplate = "pharmacies-%s"
	exchangeKind         = "topic"

	transactionTypeWorker = "worker"
)

type Worker interface {
	StopWorker()
	StartWorker(service application.Core) error
}

type workerImpl struct {
	cfg          *config.Consumer
	consumer     *rabbitmq.Consumer
	consumerName string
}

func getConsumerName() string {
	hostname, err := os.Hostname()
	if err != nil {
		goerrors.Log().WithError(err).Warn("cannot get hostname")
		return ""
	}

	return fmt.Sprintf(consumerNameTemplate, hostname)
}

func (w *workerImpl) StopWorker() {
	noWait := false
	w.consumer.StopConsuming(w.consumerName, noWait)
	w.consumer.Disconnect()
}

func (w *workerImpl) StartWorker(svc application.Core) error {
	prefetchSize := w.cfg.Concurent * 2

	routingKeys := make([]string, len(w.cfg.RoutingKeys))

	for i, item := range w.cfg.RoutingKeys {
		routingKeys[i] = string(item)
	}

	return w.consumer.StartConsuming(getConsumer(svc), w.cfg.QueueName,
		routingKeys,
		rabbitmq.WithConsumeOptionsConcurrency(w.cfg.Concurent),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName(w.cfg.Exchange),
		rabbitmq.WithConsumeOptionsQOSPrefetch(prefetchSize),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(exchangeKind),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName(w.consumerName))
}

func New(cfg *config.Consumer) (Worker, error) {
	consumer, err := rabbitmq.NewConsumer(cfg.Address, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogger(goerrors.Log()))
	if err != nil {
		return nil, err
	}

	return &workerImpl{
		cfg:          cfg,
		consumer:     &consumer,
		consumerName: getConsumerName(),
	}, nil
}
