package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/ReilEgor/SiteSentinel/checker-service/internal/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerQueueName string
type Consumer struct {
	ch     *amqp.Channel
	queue  ConsumerQueueName
	uc     domain.CheckerUsecase
	pub    domain.ResultPublisher
	logger *slog.Logger
}

func NewConsumer(ch *amqp.Channel, queue ConsumerQueueName, uc domain.CheckerUsecase, pub domain.ResultPublisher) *Consumer {
	return &Consumer{ch: ch, queue: queue, uc: uc, pub: pub, logger: slog.With(slog.String("component", "consumer"))}
}

func (c *Consumer) Start(ctx context.Context) {
	q, err := c.ch.QueueDeclare(
		string(c.queue),
		true,
		false,
		false,
		false,
		nil,
	)
	slog.Info("Consume check result")
	if err != nil {
		return
	}

	msgs, err := c.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consuming: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker received shutdown signal, stopping...")
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("RabbitMQ channel closed")
				return
			}

			c.processTask(ctx, d)
		}
	}
}

func (c *Consumer) processTask(ctx context.Context, d amqp.Delivery) {
	taskCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var site pkgModels.Site
	if err := json.Unmarshal(d.Body, &site); err != nil {
		c.logger.Error("failed to unmarshal task", slog.Any("error", err))
		d.Ack(false)
		return
	}

	result, err := c.uc.ExecuteCheck(taskCtx, site)
	if err != nil {
		c.logger.Warn("check failed, but sending result anyway",
			slog.String("url", site.URL),
			slog.Any("error", err))
	}

	err = c.pub.PublishResult(taskCtx, result)
	if err == nil {
		d.Ack(false)
	} else {
		c.logger.Error("failed to publish result", slog.Any("error", err))
		d.Nack(false, true)
	}
}
