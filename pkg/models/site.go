package models

import (
	"time"

	"github.com/google/uuid"
)

type Site struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	URL       string    `json:"url"`
	Interval  int       `json:"interval_seconds"`
	Status    string    `json:"status"`
	LastCheck time.Time `json:"last_check"`
	CreatedAt time.Time `json:"created_at"`
}

type CheckResult struct {
	SiteID     uuid.UUID     `json:"site_id"`
	IsUp       bool          `json:"is_up"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"-"`
	LatencyMs  int64         `json:"latency_ms"`
	CheckedAt  time.Time     `json:"checked_at"`
}

func (s *Site) UpdateStatus(isUp bool, checkedAt time.Time) {
	if isUp {
		s.Status = "up"
	} else {
		s.Status = "down"
	}
	s.LastCheck = checkedAt
}
