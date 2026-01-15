package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherQueueName string
type Publisher struct {
	ch     *amqp.Channel
	queue  PublisherQueueName
	logger *slog.Logger
}

func NewPublisher(ch *amqp.Channel, queue PublisherQueueName) *Publisher {
	return &Publisher{ch: ch, queue: queue, logger: slog.With(slog.String("component", "publisher"))}
}

func (p *Publisher) PublishCheckTask(ctx context.Context, site pkgModels.Site) error {
	_, err := p.ch.QueueDeclare(
		string(p.queue),
		true,
		false,
		false,
		false,
		nil,
	)
	body, err := json.Marshal(site)
	if err != nil {
		p.logger.Error("failed to marshal site",
			slog.Any("error", err),
		)
		return fmt.Errorf("marshal site error: %w", err)
	}

	return p.ch.PublishWithContext(ctx,
		"",
		string(p.queue),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
}
