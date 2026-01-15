package rabbitmq

import (
	"log/slog"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitURL string

func NewRabbitMQConn(url RabbitURL) (*amqp.Connection, func(), error) {
	logger := slog.With(slog.String("component", "rabbitmqConnector"))
	conn, err := amqp.Dial(string(url))
	if err != nil {
		logger.Error("failed to connect to rabbitmq",
			slog.Any("error", err),
			slog.String("url", maskRabbitURL(string(url))),
		)
		return nil, func() {}, err
	}

	slog.Info("successful connection to RabbitMQ",
		slog.String("component", "rabbitmq"),
		slog.String("url", maskRabbitURL(string(url))))

	cleanup := func() {
		slog.Info("closing RabbitMQ connection")

		if err := conn.Close(); err != nil {
			slog.Error("failed to close RabbitMQ connection",
				slog.String("component", "rabbitmq"),
				slog.Any("error", err))
		}
	}

	return conn, cleanup, nil
}
func NewRabbitMQChannel(conn *amqp.Connection) (*amqp.Channel, func(), error) {
	logger := slog.With(slog.String("component", "rabbitmqChannel"))
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("failed to open rabbitmq channel", slog.Any("error", err))
		return nil, func() {}, err
	}

	logger.Debug("rabbitmq channel opened")

	cleanup := func() {
		logger.Debug("closing rabbitmq channel")
		if err := ch.Close(); err != nil {
			logger.Error("failed to close rabbitmq channel", slog.Any("error", err))
		}
	}
	return ch, cleanup, nil
}
func maskRabbitURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "invalid-rabbit-url"
	}

	return u.Redacted()
}
