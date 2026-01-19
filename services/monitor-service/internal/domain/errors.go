package domain

import "errors"

var (
	// Global Errors
	ErrInvalidURL     = errors.New("invalid site URL format")
	ErrIntervalTooLow = errors.New("check interval is below minimum allowed")
	ErrInvalidUserID  = errors.New("invalid user identifier")

	// Scheduling Errors
	ErrFetchSitesFailed  = errors.New("failed to fetch sites for scheduling")
	ErrPublishTaskFailed = errors.New("failed to publish check task")

	// UserCase Errors
	ErrAddSiteFailed    = errors.New("failed to add new site")
	ErrFetchUserSites   = errors.New("failed to fetch user sites")
	ErrSaveCheckResult  = errors.New("failed to save check result")
	ErrUpdateSiteStatus = errors.New("failed to update site status")
)
