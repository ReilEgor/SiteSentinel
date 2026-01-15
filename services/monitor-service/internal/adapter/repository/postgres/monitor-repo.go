package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

type MonitorRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{db: db, logger: slog.With(slog.String("component", "monitor-repo"))}
}

func (r *MonitorRepository) SaveCheckResult(ctx context.Context, res pkgModels.CheckResult) error {
	latencyMs := res.Latency.Milliseconds()

	query := `
		INSERT INTO check_results (site_id, status_code, response_time_ms, is_up, checked_at) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		res.SiteID,
		res.StatusCode,
		latencyMs,
		res.IsUp,
		res.CheckedAt,
	)
	if err != nil {
		r.logger.Error("failed to save check result",
			slog.String("site_id", res.SiteID.String()),
			slog.Int64("latency_ms", latencyMs),
			slog.Any("error", err))
		return fmt.Errorf("save result: %w", err)
	}

	r.logger.Debug("check result saved to database",
		slog.String("site_id", res.SiteID.String()),
		slog.Int64("latency_ms", latencyMs))

	return nil
}
func (r *MonitorRepository) UpdateSiteStatus(ctx context.Context, siteID uuid.UUID, isUp bool) error {
	query := `UPDATE sites SET is_up = $1, last_check = NOW() WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, isUp, siteID)
	if err != nil {
		r.logger.Error("failed to update site status",
			slog.String("site_id", siteID.String()),
			slog.Any("error", err))
		return fmt.Errorf("update status: %w", err)
	}

	r.logger.Debug("site status updated",
		slog.String("site_id", siteID.String()),
		slog.Bool("is_up", isUp))

	return nil
}
func (r *MonitorRepository) GetSitesToSchedule(ctx context.Context) ([]pkgModels.Site, error) {
	query := `
       SELECT id, url, interval_seconds 
       FROM sites 
       WHERE last_check IS NULL 
          OR last_check + (interval_seconds * interval '1 second') < NOW()`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying sites to schedule: %w", err)
	}
	defer rows.Close()

	var sites []pkgModels.Site
	for rows.Next() {
		var s pkgModels.Site
		if err := rows.Scan(&s.ID, &s.URL, &s.Interval); err != nil {
			r.logger.Error("failed to scan site row", slog.Any("error", err))
			continue
		}
		sites = append(sites, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return sites, nil
}
