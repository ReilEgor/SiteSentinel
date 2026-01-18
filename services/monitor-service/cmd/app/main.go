package main

import (
	"context"
	"log/slog"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/adapter/broker/rabbitmq"
	rabbitmqConst "github.com/ReilEgor/SiteSentinel/pkg/rabbitmq"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	//TODO: configure logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "monitor-service"))

	dsn := os.Getenv("DB_SOURCE")
	rabbitURL := os.Getenv("RABBIT_URL")
	resQueue := rabbitmq.ConsumerQueueName(rabbitmqConst.CheckResultsQueue)
	taskQueue := rabbitmq.PublisherQueueName(rabbitmqConst.CheckTasksQueue)

	app, cleanup, err := InitializeApp(dsn, rabbitmq.RabbitURL(rabbitURL), resQueue, taskQueue)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("monitor service is starting")

	go func() {
		logger.Info("starting scheduler")
		app.Scheduler.Run(ctx)
	}()

	go func() {
		logger.Info("starting consumer")
		app.Consumer.Start(ctx)
	}()

	go func() {
		port := os.Getenv("HTTP_PORT")
		if port == "" {
			port = "8080"
		}

		logger.Info("starting http server", slog.String("port", port))
		if err := app.Server.Run(":" + port); err != nil {
			logger.Error("http server failed", slog.Any("error", err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
