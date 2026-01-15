package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ReilEgor/SiteSentinel/checker-service/internal/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
)

type CheckerInteractor struct {
	checker   domain.Checker
	publisher domain.ResultPublisher
}

func NewCheckerUsecase(c domain.Checker, p domain.ResultPublisher) *CheckerInteractor {
	return &CheckerInteractor{checker: c, publisher: p}
}

func (uc *CheckerInteractor) ExecuteCheck(ctx context.Context, site pkgModels.Site) (pkgModels.CheckResult, error) {
	result := pkgModels.CheckResult{
		SiteID:    site.ID,
		CheckedAt: time.Now(),
		IsUp:      false,
	}

	isUp, statusCode, latency, err := uc.checker.Check(ctx, site.URL)

	if err != nil {
		return result, fmt.Errorf("failed to check site %s: %w", site.URL, err)
	}

	result.IsUp = isUp
	result.StatusCode = statusCode
	result.Latency = latency

	slog.Info("checked site successfully",
		slog.String("url", site.URL),
		slog.Bool("is_up", isUp))

	return result, nil
}
