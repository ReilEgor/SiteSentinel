package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	mocks "github.com/ReilEgor/SiteSentinel/checker-service/internal/mocks/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckerInteractor_ExecuteCheck(t *testing.T) {

	mockChecker := mocks.NewChecker(t)

	uc := &CheckerInteractor{checker: mockChecker}

	site := pkgModels.Site{
		ID:  uuid.New(),
		URL: "https://mysite.com",
	}

	tests := []struct {
		name       string
		mockReturn struct {
			isUp       bool
			statusCode int
			latency    time.Duration
			err        error
		}
		wantErr      bool
		expectIsUp   bool
		expectStatus int
	}{
		{
			name: "Successful check (200 OK)",
			mockReturn: struct {
				isUp       bool
				statusCode int
				latency    time.Duration
				err        error
			}{true, 200, 100 * time.Millisecond, nil},
			wantErr:      false,
			expectIsUp:   true,
			expectStatus: 200,
		},
		{
			name: "Site is down (500 Internal Error)",
			mockReturn: struct {
				isUp       bool
				statusCode int
				latency    time.Duration
				err        error
			}{false, 500, 50 * time.Millisecond, nil},
			wantErr:      false,
			expectIsUp:   false,
			expectStatus: 500,
		},
		{
			name: "Network error (DNS failure)",
			mockReturn: struct {
				isUp       bool
				statusCode int
				latency    time.Duration
				err        error
			}{false, 0, 0, errors.New("cannot resolve host")},
			wantErr:    true,
			expectIsUp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChecker.On("Check", mock.Anything, site.URL).
				Return(tt.mockReturn.isUp, tt.mockReturn.statusCode, tt.mockReturn.latency, tt.mockReturn.err).
				Once()

			res, err := uc.ExecuteCheck(context.Background(), site)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectIsUp, res.IsUp)
				assert.Equal(t, tt.expectStatus, res.StatusCode)
			}
			assert.Equal(t, site.ID, res.SiteID)
		})
	}
}
