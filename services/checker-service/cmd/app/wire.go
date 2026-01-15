//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ReilEgor/SiteSentinel/checker-service/internal/adapter/broker/rabbitmq"
	"github.com/ReilEgor/SiteSentinel/checker-service/internal/adapter/http_client"
	"github.com/ReilEgor/SiteSentinel/checker-service/internal/domain"
	"github.com/ReilEgor/SiteSentinel/checker-service/internal/usecase"
	"github.com/google/wire"
)

var HTTPSet = wire.NewSet(
	http_client.NewHTTPChecker,
	wire.Bind(new(domain.Checker), new(*http_client.HTTPChecker)),
)

var BrokerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewPublisher,
	rabbitmq.NewConsumer,

	wire.Bind(new(domain.ResultPublisher), new(*rabbitmq.Publisher)),
)

var UsecaseSet = wire.NewSet(
	usecase.NewCheckerUsecase,
	wire.Bind(new(domain.CheckerUsecase), new(*usecase.CheckerInteractor)),
)

type App struct {
	Logic    domain.CheckerUsecase
	Consumer *rabbitmq.Consumer
}

func InitializeApp(
	rabbitURL rabbitmq.RabbitURL,
	resultsQueue rabbitmq.ConsumerQueueName,
	taskQueue rabbitmq.PublisherQueueName,
) (*App, func(), error) {
	wire.Build(
		HTTPSet,
		UsecaseSet,
		BrokerSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
