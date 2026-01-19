package domain

import (
	"net/url"
	"time"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

const (
	MinCheckInterval = 60 // seconds
)

func NewSite(userID uuid.UUID, rawURL string) (*pkgModels.Site, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	return &pkgModels.Site{
		ID:        uuid.New(),
		UserID:    userID,
		URL:       rawURL,
		Status:    "pending",
		CreatedAt: time.Now(),
	}, nil
}
