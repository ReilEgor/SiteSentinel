package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	mocks "github.com/ReilEgor/SiteSentinel/monitor-service/internal/mocks/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
)

func TestSiteScheduler_Execute(t *testing.T) {
	ctx := context.Background()

	testSites := []pkgModels.Site{
		{ID: domain.IntToUUID(1), URL: "https://test.com"},
	}
	t.Run("success execution", func(t *testing.T) {
		mockPub := mocks.NewTaskPublisher(t)
		mockRepo := mocks.NewMonitorRepository(t)
		mockRepo.On("GetSitesToSchedule", ctx).Return(testSites, nil).Once()
		mockPub.On("PublishCheckTask", ctx, testSites[0]).Return(nil).Once()

		s := NewSiteScheduler(mockRepo, mockPub)
		s.execute(ctx)
	})
	t.Run("no sites to schedule", func(t *testing.T) {
		mockPub := mocks.NewTaskPublisher(t)
		mockRepo := mocks.NewMonitorRepository(t)
		mockRepo.On("GetSitesToSchedule", ctx).Return([]pkgModels.Site{}, nil).Once()

		s := NewSiteScheduler(mockRepo, mockPub)
		s.execute(ctx)

		mockPub.AssertNotCalled(t, "PublishCheckTask", mock.Anything, mock.Anything)
	})
	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMonitorRepository(t)
		mockPub := mocks.NewTaskPublisher(t)

		dbErr := errors.New("connection refused")
		mockRepo.On("GetSitesToSchedule", ctx).Return(nil, dbErr).Once()

		s := NewSiteScheduler(mockRepo, mockPub)

		err := s.execute(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFetchSitesFailed)

		mockPub.AssertNotCalled(t, "PublishCheckTask", mock.Anything, mock.Anything)
	})
}
