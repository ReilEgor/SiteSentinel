package main

import (
	"context"
	"github.com/ReilEgor/SiteSentinel/checker-service/internal/adapter/broker/rabbitmq"
	rabbitmqConst "github.com/ReilEgor/SiteSentinel/pkg/rabbitmq"
	"log/slog"
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

	rabbitURL := os.Getenv("RABBIT_URL")
	resQueue := rabbitmq.ConsumerQueueName(rabbitmqConst.CheckTasksQueue)
	taskQueue := rabbitmq.PublisherQueueName(rabbitmqConst.CheckResultsQueue)

	app, cleanup, err := InitializeApp(rabbitmq.RabbitURL(rabbitURL), resQueue, taskQueue)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("checker service is starting")

	go func() {
		logger.Info("starting consumer")
		app.Consumer.Start(ctx)
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")
}
