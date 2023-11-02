package worker

import (
	"context"

	"github.com/dwnGnL/testWork/internal/application"
	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/lib/goerrors"
	"github.com/sirupsen/logrus"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func getConsumer(svc application.Core) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) (action rabbitmq.Action) {

		switch config.RoutingKey(d.RoutingKey) {
		case config.RoutingTest:
			return consumePharmacy(
				context.Background(),
				svc,
				d.Body,
			)

		default:
			goerrors.Log().WithFields(logrus.Fields{
				"body":        d.Body,
				"routing_key": d.RoutingKey,
			}).Warn("unknown routing")
			return rabbitmq.NackDiscard
		}
	}
}

func consumePharmacy(ctx context.Context, svc application.Core, body []byte) rabbitmq.Action {

	return rabbitmq.Ack
}
