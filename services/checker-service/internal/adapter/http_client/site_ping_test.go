package http_client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPChecker_Check(t *testing.T) {
	t.Run("success_request", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		checker := NewHTTPChecker()

		isUp, status, latency, err := checker.Check(context.Background(), ts.URL)

		assert.NoError(t, err)
		assert.True(t, isUp)
		assert.Equal(t, http.StatusOK, status)
		assert.Greater(t, latency.Nanoseconds(), int64(0))
	})

	t.Run("server_error_500", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		checker := NewHTTPChecker()
		isUp, status, _, err := checker.Check(context.Background(), ts.URL)

		assert.NoError(t, err)
		assert.False(t, isUp)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
