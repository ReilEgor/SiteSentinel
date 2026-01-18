package usecase

import (
	"context"
	"fmt"
	"log/slog"

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
	return i.repo.AddSite(ctx, userID, url, interval)
}

func (i *MonitorInteractor) GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error) {
	return i.repo.GetUserSites(ctx, userID)
}

func (i *MonitorInteractor) ProcessCheckResult(ctx context.Context, result pkgModels.CheckResult) error {
	i.logger.Info("processing check result from worker")

	if err := i.repo.SaveCheckResult(ctx, result); err != nil {
		i.logger.Error("failed to save check result to history", slog.Any("error", err))
		return fmt.Errorf("save check result: %w", err)
	}

	if err := i.repo.UpdateSiteStatus(ctx, result.SiteID, result.IsUp); err != nil {
		i.logger.Error("failed to update site current status", slog.Any("error", err))
		return fmt.Errorf("update site status: %w", err)
	}

	return nil
}
