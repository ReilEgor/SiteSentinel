package usecase

import (
	"context"
	"fmt"
	"log/slog"

	netURL "net/url"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

type MonitorInteractor struct {
	repo      domain.MonitorRepository
	publisher domain.TaskPublisher
	logger    *slog.Logger
}

func NewMonitorInteractor(repo domain.MonitorRepository, pub domain.TaskPublisher) *MonitorInteractor {
	return &MonitorInteractor{
		repo:      repo,
		publisher: pub,
		logger:    slog.With(slog.String("component", "scheduler")),
	}
}
func (i *MonitorInteractor) AddSite(ctx context.Context, userID uuid.UUID, url string, interval int) (*pkgModels.Site, error) {
	i.logger.Debug("adding new site",
		slog.String("url", url),
		slog.Int("interval", interval))

	parsedURL, err := netURL.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, domain.ErrInvalidURL
	}
	if interval < domain.MinCheckInterval {
		return nil, domain.ErrIntervalTooLow
	}

	site, err := i.repo.AddSite(ctx, userID, url, interval)
	if err != nil {
		i.logger.Error("failed to add site",
			slog.String("user_id", userID.String()),
			slog.String("url", url),
			slog.Int("interval", interval),
			slog.Any("error", err.Error()),
		)
		return nil, fmt.Errorf("%w: %v", domain.ErrAddSiteFailed, err)
	}

	i.logger.Debug("site added successfully",
		slog.String("user_id", userID.String()),
		slog.String("site_id", site.ID.String()),
		slog.String("url", url),
		slog.Int("interval", interval),
	)

	return site, nil
}

func (i *MonitorInteractor) GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error) {
	i.logger.Debug("getting user sites")
	if userID == uuid.Nil {
		i.logger.Warn("userID is required")
		return nil, domain.ErrInvalidUserID
	}
	sites, err := i.repo.GetUserSites(ctx, userID)
	if err != nil {
		i.logger.Error("failed to get user sites",
			slog.String("user_id", userID.String()),
			slog.Any("error", fmt.Errorf("%w: %v", domain.ErrFetchUserSites, err)),
		)
		return nil, fmt.Errorf("%w: %v", domain.ErrFetchUserSites, err)
	}

	i.logger.Debug("fetched user sites",
		slog.String("user_id", userID.String()),
		slog.Int("site_count", len(sites)),
	)

	return sites, nil
}

func (i *MonitorInteractor) ProcessCheckResult(ctx context.Context, result pkgModels.CheckResult) error {
	i.logger.Debug("processing check result from worker")

	if err := i.repo.SaveCheckResult(ctx, result); err != nil {
		i.logger.Error("failed to save check result to history", slog.Any("error", fmt.Errorf("%w: %v", domain.ErrSaveCheckResult, err)))
		return fmt.Errorf("%w: %v", domain.ErrSaveCheckResult, err)
	}

	if err := i.repo.UpdateSiteStatus(ctx, result.SiteID, result.IsUp); err != nil {
		i.logger.Error("failed to update site current status", slog.Any("error", fmt.Errorf("%w: %v", domain.ErrUpdateSiteStatus, err)))
		return fmt.Errorf("%w: %v", domain.ErrUpdateSiteStatus, err)
	}

	i.logger.Debug("site updated successfully",
		slog.String("site_id", result.SiteID.String()),
		slog.Bool("is_up", result.IsUp),
	)

	return nil
}
