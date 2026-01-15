package domain

import (
	"errors"
	"net/url"
	"time"

	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
)

var (
	ErrInvalidURL = errors.New("invalid website URL")
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
