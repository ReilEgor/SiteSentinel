package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	MAX_OPEN_CONNS = 25
	MAX_IDLE_CONNS = 25
)

func NewPostgresDB(dsn string) (*sql.DB, func(), error) {
	slog.With(slog.String("component", "postgres"))
	slog.Info("connecting to database",
		slog.String("dsn", maskDSN(dsn)),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open db: %w", err)
	}

	db.SetMaxOpenConns(MAX_OPEN_CONNS)
	db.SetMaxIdleConns(MAX_IDLE_CONNS)
	db.SetConnMaxLifetime(5 * time.Minute)

	start := time.Now()
	if err := db.Ping(); err != nil {
		slog.Error("database ping failed",
			slog.Any("error", err),
			slog.Duration("duration", time.Since(start)))
		return nil, nil, err
	}

	slog.Info("successful connection to PostgreSQL",
		slog.Duration("latency", time.Since(start)),
		slog.Int("max_open_conns", MAX_OPEN_CONNS))

	cleanup := func() {
		slog.Info("closing database connections")
		if err := db.Close(); err != nil {
			slog.Error("error closing database",
				slog.Any("error", err))
		}
	}

	return db, cleanup, nil
}
func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil {
		return "invalid-dsn"
	}

	return u.Redacted()
}
