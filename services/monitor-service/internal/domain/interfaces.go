package domain

import (
	"context"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

//go:generate mockery --name MonitorRepository --output ../mocks/domain --outpkg domain --case=underscore
type MonitorRepository interface {
	SaveCheckResult(ctx context.Context, res pkgModels.CheckResult) error
	UpdateSiteStatus(ctx context.Context, siteID uuid.UUID, isUp bool) error
	GetSitesToSchedule(ctx context.Context) ([]pkgModels.Site, error)
	GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error)
	AddSite(ctx context.Context, userID uuid.UUID, url string, interval int) (*pkgModels.Site, error)
	DeleteSite(ctx context.Context, userID uuid.UUID, url string) (*pkgModels.Site, error)
}

//go:generate mockery --name TaskPublisher --output ../mocks/domain --outpkg domain --case=underscore
type TaskPublisher interface {
	PublishCheckTask(ctx context.Context, site pkgModels.Site) error
}

//go:generate mockery --name SiteUsecase --output ../mocks/domain --outpkg domain --case=underscore
type SiteUsecase interface {
	AddSite(ctx context.Context, userID uuid.UUID, url string, interval int) (*pkgModels.Site, error)
	DeleteSite(ctx context.Context, userID uuid.UUID, url string) (*pkgModels.Site, error)
	GetUserSites(ctx context.Context, userID uuid.UUID) ([]*pkgModels.Site, error)
	ProcessCheckResult(ctx context.Context, result pkgModels.CheckResult) error
}
