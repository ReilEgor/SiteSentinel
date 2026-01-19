package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	mocks "github.com/ReilEgor/SiteSentinel/monitor-service/internal/mocks/domain"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"context"
	"testing"
)

func TestHandler_AddSite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		reqBody := struct {
			UserID   uuid.UUID `json:"userID"`
			URL      string    `json:"url"`
			Interval int       `json:"interval"`
		}{
			UserID:   domain.IntToUUID(1),
			URL:      "https://google.com",
			Interval: domain.MinCheckInterval,
		}

		expectedSite := &pkgModels.Site{
			UserID: reqBody.UserID,
			URL:    reqBody.URL,
		}
		body, _ := json.Marshal(reqBody)

		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)

		mockUC.On("AddSite", ctx, reqBody.UserID, reqBody.URL, reqBody.Interval).
			Return(expectedSite, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/addSite", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddSite(c)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response *pkgModels.Site
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, reqBody.UserID, response.UserID)
	})
	t.Run("bad req", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"userID":   domain.IntToUUID(1),
			"badurl":   "https://google.com",
			"interval": domain.MinCheckInterval,
		}
		body, _ := json.Marshal(reqBody)

		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/addSite", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddSite(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("uc error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"userID":   domain.IntToUUID(1),
			"url":      "https://google.com",
			"interval": domain.MinCheckInterval,
		}
		body, _ := json.Marshal(reqBody)

		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)
		mockUC.On("AddSite", ctx, reqBody["userID"], reqBody["url"], reqBody["interval"]).
			Return(nil, domain.ErrAddSiteFailed).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/addSite", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.AddSite(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetUserSites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ans := &pkgModels.Site{
			ID:        domain.IntToUUID(1),
			UserID:    domain.IntToUUID(1),
			URL:       "https://google.com",
			Interval:  60,
			Status:    "up",
			LastCheck: &time.Time{},
			CreatedAt: time.Time{},
			IsUp:      true,
		}

		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)

		expectedSites := []*pkgModels.Site{ans}
		mockUC.On("GetUserSites", ctx, domain.IntToUUID(1)).
			Return(expectedSites, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/getUserSites/:userID", nil)
		c.Params = []gin.Param{{Key: "userID", Value: domain.IntToUUID(1).String()}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.GetUserSites(c)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string][]*pkgModels.Site
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response["sites"], 1)
		assert.Equal(t, ans.URL, response["sites"][0].URL)
	})
	t.Run("invalid userID", func(t *testing.T) {
		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/getUserSites/:userID", nil)
		c.Params = []gin.Param{{Key: "userID", Value: "bad_user_id"}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.GetUserSites(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("uc error", func(t *testing.T) {
		mockUC := mocks.NewSiteUsecase(t)
		h := NewHandler(mockUC)

		ucErr := errors.New("usecase-error")
		mockUC.On("GetUserSites", ctx, domain.IntToUUID(1)).
			Return(nil, ucErr).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/getUserSites/:userID", nil)
		c.Params = []gin.Param{{Key: "userID", Value: domain.IntToUUID(1).String()}}
		c.Request.Header.Set("Content-Type", "application/json")

		h.GetUserSites(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
