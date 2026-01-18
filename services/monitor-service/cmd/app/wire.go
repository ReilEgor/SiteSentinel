//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/broker/rabbitmq"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/repository/postgres"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/transport/rest"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/transport/rest/handler"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/scheduler"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/usecase"
	"github.com/google/wire"
)

var SchedulerSet = wire.NewSet(
	scheduler.NewSiteScheduler,
)

var BrokerSet = wire.NewSet(
	rabbitmq.NewRabbitMQConn,
	rabbitmq.NewRabbitMQChannel,
	rabbitmq.NewPublisher,
	rabbitmq.NewConsumer,
	wire.Bind(new(domain.TaskPublisher), new(*rabbitmq.Publisher)),
)

var RepositorySet = wire.NewSet(
	postgres.NewPostgresDB,
	postgres.NewMonitorRepository,
	wire.Bind(new(domain.MonitorRepository), new(*postgres.MonitorRepository)),
)

var UsecaseSet = wire.NewSet(
	usecase.NewMonitorInteractor,
	wire.Bind(new(domain.SiteUsecase), new(*usecase.MonitorInteractor)),
)

var HTTPSet = wire.NewSet(
	rest.NewGinServer,
	handler.NewHandler,
)

type App struct {
	Logic     domain.SiteUsecase
	Scheduler *scheduler.SiteScheduler
	Consumer  *rabbitmq.Consumer
	Server    *rest.GinServer
}

func InitializeApp(
	dsn string,
	rabbitURL rabbitmq.RabbitURL,
	resultsQueue rabbitmq.ConsumerQueueName,
	taskQueue rabbitmq.PublisherQueueName,
) (*App, func(), error) {
	wire.Build(
		RepositorySet,
		UsecaseSet,
		BrokerSet,
		SchedulerSet,
		HTTPSet,
		wire.Struct(new(App), "*"),
	)
	return nil, nil, nil
}
