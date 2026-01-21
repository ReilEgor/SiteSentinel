package handler

import (
	"log/slog"
	"net/http"

	"github.com/ReilEgor/SiteSentinel/monitor-service/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	uc     domain.SiteUsecase
	logger *slog.Logger
}

func NewHandler(uc domain.SiteUsecase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("component", "handler")),
	}
}

func (h *Handler) AddSite(c *gin.Context) {
	var req struct {
		UserID   uuid.UUID `json:"userID" binding:"required"`
		URL      string    `json:"url" binding:"required,url"`
		Interval int       `json:"interval" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	ctx := c.Request.Context()

	site, err := h.uc.AddSite(ctx, req.UserID, req.URL, req.Interval)
	if err != nil {
		h.logger.Error("failed to add site", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add site"})
		return
	}

	c.JSON(http.StatusCreated, site)
}

func (h *Handler) GetUserSites(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	h.logger.Info("GetUserSites called", slog.String("userID", userIDStr))
	if err != nil {
		h.logger.Warn("invalid userID parameter", slog.String("value", userIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID parameter"})
		return
	}

	ctx := c.Request.Context()
	sites, err := h.uc.GetUserSites(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get user sites", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve sites"})
		return
	}

	if sites == nil {
		h.logger.Info("no sites found for user", slog.String("uuid", userID.String()))
	}

	c.JSON(http.StatusOK, gin.H{"sites": sites})
}

func (h *Handler) DeleteSite(c *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"userID" binding:"required"`
		URL    string    `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	ctx := c.Request.Context()

	site, err := h.uc.DeleteSite(ctx, req.UserID, req.URL)
	if err != nil {
		h.logger.Error("failed to delete site", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add site"})
		return
	}

	c.JSON(http.StatusCreated, site)
}
