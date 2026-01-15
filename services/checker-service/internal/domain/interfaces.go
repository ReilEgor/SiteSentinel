package domain

import (
	"context"
	"time"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
)

type Checker interface {
	Check(ctx context.Context, url string) (bool, int, time.Duration, error)
}

type ResultPublisher interface {
	PublishResult(ctx context.Context, result pkgModels.CheckResult) error
}

type CheckerUsecase interface {
	ExecuteCheck(ctx context.Context, site pkgModels.Site) (pkgModels.CheckResult, error)
}
