package postgres

import (
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pkgModels "github.com/ReilEgor/SiteSentinel/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestMonitorRepository_SaveCheckResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub db: %s", err)
	}
	defer db.Close()

	repo := NewMonitorRepository(db)

	now := time.Now()
	res := pkgModels.CheckResult{
		SiteID:     uuid.New(),
		StatusCode: 200,
		Latency:    500 * time.Millisecond,
		IsUp:       true,
		CheckedAt:  now,
	}
	mock.ExpectExec("INSERT INTO check_results").
		WithArgs(
			res.SiteID,
			res.StatusCode,
			res.Latency.Milliseconds(),
			res.IsUp,
			res.CheckedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveCheckResult(context.Background(), res)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMonitorRepository_UpdateSiteStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewMonitorRepository(db)
	siteID := pkgModels.IntToUUID(1)

	mock.ExpectExec("UPDATE sites SET is_up = \\$1, last_check = NOW\\(\\) WHERE id = \\$2").
		WithArgs(true, siteID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateSiteStatus(context.Background(), siteID, true)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMonitorRepository_GetSitesToSchedule(t *testing.T) {

}

func TestMonitorRepository_AddSite(t *testing.T) {

}

func TestMonitorRepository_GetUserSites(t *testing.T) {

}
