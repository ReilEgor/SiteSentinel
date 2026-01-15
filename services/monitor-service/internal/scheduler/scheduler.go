package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
)

type SiteScheduler struct {
	repo      domain.MonitorRepository
	publisher domain.TaskPublisher
	interval  time.Duration
	logger    *slog.Logger
}

func NewSiteScheduler(repo domain.MonitorRepository, pub domain.TaskPublisher) *SiteScheduler {
	return &SiteScheduler{
		repo:      repo,
		publisher: pub,
		interval:  30 * time.Second,
		logger:    slog.With(slog.String("component", "scheduler")),
	}
}

func (s *SiteScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	s.logger.Info("scheduler started",
		slog.String("interval", s.interval.String()))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopping")
			return
		case <-ticker.C:
			s.execute(ctx)
		}
	}
}

func (s *SiteScheduler) execute(ctx context.Context) {
	start := time.Now()
	sites, err := s.repo.GetSitesToSchedule(ctx)
	if err != nil {
		s.logger.Error("failed to get sites from repository",
			slog.Any("error", err))
		return
	}

	if len(sites) == 0 {
		s.logger.Debug("no sites found to schedule")
		return
	}

	s.logger.Info("found sites to check",
		slog.Int("count", len(sites)))

	for _, site := range sites {
		if err := s.publisher.PublishCheckTask(ctx, site); err != nil {
			s.logger.Error("failed to publish task",
				slog.String("site_url", site.URL),
				slog.Any("site_id", site.ID),
				slog.Any("error", err))
			continue
		}

		s.logger.Info("task published successfully",
			slog.String("site_url", site.URL),
			slog.Duration("took", time.Since(start)))
	}
}
