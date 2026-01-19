package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	mocks "github.com/ReilEgor/SiteSentinel/monitor-service/internal/mocks/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMonitorInteractor_AddSite(t *testing.T) {
	ctx := context.Background()
	t.Run("invalid url", func(t *testing.T) {
		userID := domain.IntToUUID(1)
		url := "not-a-url"
		interval := domain.MinCheckInterval

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)
		site, err := interactor.AddSite(ctx, userID, url, interval)
		assert.Nil(t, site)
		assert.ErrorIs(t, err, domain.ErrInvalidURL)
	})
	t.Run("out of bound", func(t *testing.T) {
		userID := domain.IntToUUID(1)
		url := "https://www.test.com/"
		interval := 20

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)
		site, err := interactor.AddSite(ctx, userID, url, interval)
		assert.Nil(t, site)
		assert.ErrorIs(t, err, domain.ErrIntervalTooLow)
	})
	t.Run("repository error", func(t *testing.T) {
		userID := domain.IntToUUID(1)
		url := "https://www.test.com/"
		interval := 70

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		dbErr := errors.New("db connection lost")
		mockRepo.On("AddSite", ctx, userID, url, interval).
			Return(nil, dbErr).Once()

		site, err := interactor.AddSite(ctx, userID, url, interval)

		assert.Nil(t, site)
		assert.ErrorIs(t, err, domain.ErrAddSiteFailed)
		assert.Contains(t, err.Error(), dbErr.Error())
	})
	t.Run("no error", func(t *testing.T) {
		userID := domain.IntToUUID(1)
		url := "https://www.test.com/"
		interval := domain.MinCheckInterval

		expectedSite := &pkgModels.Site{
			ID:       userID,
			URL:      url,
			Interval: interval,
		}

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		mockRepo.On("AddSite", ctx, userID, url, interval).Return(expectedSite, nil).Once()

		site, err := interactor.AddSite(ctx, userID, url, interval)
		assert.NotNil(t, site)
		assert.Nil(t, err)
		assert.Equal(t, expectedSite, site)
	})
}

func TestMonitorInteractor_GetUserSites(t *testing.T) {
	ctx := context.Background()
	t.Run("nil userID", func(t *testing.T) {
		userID := uuid.Nil

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		sites, err := interactor.GetUserSites(ctx, userID)
		assert.Nil(t, sites)
		assert.ErrorIs(t, err, domain.ErrInvalidUserID)
	})
	t.Run("repository error", func(t *testing.T) {
		userID := domain.IntToUUID(1)

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		dbErr := errors.New("db connection lost")
		mockRepo.On("GetUserSites", ctx, userID).
			Return(nil, dbErr).Once()

		sites, err := interactor.GetUserSites(ctx, userID)
		assert.Nil(t, sites)
		assert.ErrorIs(t, err, domain.ErrFetchUserSites)
		assert.Contains(t, err.Error(), dbErr.Error())
	})
	t.Run("no error", func(t *testing.T) {
		userID := domain.IntToUUID(1)
		expectedSites := []*pkgModels.Site{
			{ID: userID, URL: "https://site1.com", Interval: 60},
		}
		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		mockRepo.On("GetUserSites", ctx, userID).
			Return(expectedSites, nil).Once()

		sites, err := interactor.GetUserSites(ctx, userID)

		assert.Nil(t, err)
		assert.Equal(t, expectedSites, sites)
	})
}

func TestMonitorInteractor_ProcessCheckResult(t *testing.T) {
	ctx := context.Background()

	t.Run("filed save check result(repository)", func(t *testing.T) {
		result := pkgModels.CheckResult{
			SiteID:  domain.IntToUUID(1),
			IsUp:    true,
			Latency: 120,
		}

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		mockRepo.On("SaveCheckResult", ctx, result).
			Return(domain.ErrSaveCheckResult).Once()

		err := interactor.ProcessCheckResult(ctx, result)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrSaveCheckResult)
		assert.Contains(t, err.Error(), domain.ErrSaveCheckResult.Error())
	})
	t.Run("filed update site status(repository)", func(t *testing.T) {
		result := pkgModels.CheckResult{
			SiteID:  domain.IntToUUID(1),
			IsUp:    true,
			Latency: 120,
		}

		mockRepo := mocks.NewMonitorRepository(t)
		interactor := NewMonitorInteractor(mockRepo, nil)

		mockRepo.On("SaveCheckResult", ctx, result).
			Return(nil).Once()

		mockRepo.On("UpdateSiteStatus", ctx, result.SiteID, result.IsUp).
			Return(domain.ErrUpdateSiteStatus).Once()

		err := interactor.ProcessCheckResult(ctx, result)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, domain.ErrUpdateSiteStatus)
		assert.Contains(t, err.Error(), domain.ErrUpdateSiteStatus.Error())
	})
}
