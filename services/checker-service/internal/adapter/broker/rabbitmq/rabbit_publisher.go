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

func (p *Publisher) PublishResult(ctx context.Context, result pkgModels.CheckResult) error {
	_, err := p.ch.QueueDeclare(
		string(p.queue),
		true,
		false,
		false,
		false,
		nil,
	)
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal check result: %w", err)
	}
	slog.Info("Publishing check result")
	return p.ch.PublishWithContext(ctx,
		"",
		string(p.queue),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
