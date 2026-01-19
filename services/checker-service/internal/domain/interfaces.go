package domain

import (
	"context"
	"time"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
)

//go:generate mockery --name Checker --output ../mocks/domain --outpkg domain --case=underscore
type Checker interface {
	Check(ctx context.Context, url string) (bool, int, time.Duration, error)
}

//go:generate mockery --name ResultPublisher --output ../mocks/domain --outpkg domain --case=underscore
type ResultPublisher interface {
	PublishResult(ctx context.Context, result pkgModels.CheckResult) error
}

//go:generate mockery --name CheckerUsecase --output ../mocks/domain --outpkg domain --case=underscore
type CheckerUsecase interface {
	ExecuteCheck(ctx context.Context, site pkgModels.Site) (pkgModels.CheckResult, error)
}
