package domain

import (
	"context"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

type MonitorRepository interface {
	SaveCheckResult(ctx context.Context, res pkgModels.CheckResult) error
	UpdateSiteStatus(ctx context.Context, siteID uuid.UUID, isUp bool) error
	GetSitesToSchedule(ctx context.Context) ([]pkgModels.Site, error)
}

type TaskPublisher interface {
	PublishCheckTask(ctx context.Context, site pkgModels.Site) error
}

type SiteUsecase interface {
	AddSite(ctx context.Context, userID uuid.UUID, url string) (*pkgModels.Site, error)
	GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error)
	ProcessCheckResult(ctx context.Context, result pkgModels.CheckResult) error
}
