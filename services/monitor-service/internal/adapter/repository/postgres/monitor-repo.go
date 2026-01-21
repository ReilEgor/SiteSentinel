package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *MonitorRepository) AddSite(ctx context.Context, userID uuid.UUID, url string, interval int) (*pkgModels.Site, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	querySite := `
        INSERT INTO sites (url, interval_seconds) 
        VALUES ($1, $2) 
        ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url 
        RETURNING id, url, interval_seconds, is_up, last_check, created_at`

	var created pkgModels.Site
	err = tx.QueryRowContext(ctx, querySite, url, interval).Scan(
		&created.ID, &created.URL, &created.Interval, &created.IsUp, &created.LastCheck, &created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert site error: %w", err)
	}

	queryLink := `
        INSERT INTO user_sites (user_id, site_id) 
        VALUES ($1, $2) 
        ON CONFLICT DO NOTHING`

	if _, err := tx.ExecContext(ctx, queryLink, userID, created.ID); err != nil {
		return nil, fmt.Errorf("link user to site error: %w", err)
	}

	return &created, tx.Commit()
}

func (r *MonitorRepository) GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error) {
	query := `
        SELECT s.id, s.url, s.interval_seconds, s.is_up, s.last_check, s.created_at 
        FROM sites s
        JOIN user_sites us ON s.id = us.site_id
        WHERE us.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*pkgModels.Site
	for rows.Next() {
		s := &pkgModels.Site{}
		if err := rows.Scan(&s.ID, &s.URL, &s.Interval, &s.IsUp, &s.LastCheck, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.UserID = userID
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

func (r *MonitorRepository) DeleteSite(ctx context.Context, userID uuid.UUID, url string) (*pkgModels.Site, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var site pkgModels.Site
	querySelect := `
        SELECT s.id, s.url, s.interval_seconds, s.is_up, s.last_check, s.created_at 
        FROM sites s
        JOIN user_sites us ON s.id = us.site_id
        WHERE us.user_id = $1 AND s.url = $2`

	err = tx.QueryRowContext(ctx, querySelect, userID, url).Scan(
		&site.ID, &site.URL, &site.Interval, &site.IsUp, &site.LastCheck, &site.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("site not found for this user")
		}
		return nil, fmt.Errorf("query site: %w", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM user_sites WHERE user_id = $1 AND site_id = $2`, userID, site.ID)
	if err != nil {
		return nil, fmt.Errorf("delete user_site link: %w", err)
	}

	queryCleanOrphan := `
        DELETE FROM sites 
        WHERE id = $1 AND NOT EXISTS (SELECT 1 FROM user_sites WHERE site_id = $1)`

	_, err = tx.ExecContext(ctx, queryCleanOrphan, site.ID)
	if err != nil {
		return nil, fmt.Errorf("cleanup orphaned site: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &site, nil
}
