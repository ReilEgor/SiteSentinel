package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/usecase"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerQueueName string

type Consumer struct {
	ch     *amqp.Channel
	queue  ConsumerQueueName
	uc     domain.SiteUsecase
	logger *slog.Logger
}

func NewConsumer(ch *amqp.Channel, queue ConsumerQueueName, uc *usecase.MonitorInteractor) *Consumer {
	return &Consumer{ch: ch, queue: queue, uc: uc, logger: slog.With(slog.String("component", "consumer"))}
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
	if err != nil {
		return
	}

	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		c.logger.Info("failed to consume",
			slog.String("queue", q.Name),
			slog.Any("error", err.Error()),
		)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			c.processResult(ctx, d)
		}
	}
}

func (w *Consumer) processResult(ctx context.Context, d amqp.Delivery) {
	var res pkgModels.CheckResult
	json.Unmarshal(d.Body, &res)

	if err := w.uc.ProcessCheckResult(ctx, res); err != nil {
		d.Nack(false, true)
		return
	}
	d.Ack(false)
}
