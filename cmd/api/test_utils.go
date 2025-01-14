package main

import (
	"github.com/fouched/social/internal/auth"
	"github.com/fouched/social/internal/repo"
	"github.com/fouched/social/internal/repo/cache"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	//logger := zap.Must(zap.NewProduction()).Sugar()
	mockRepo := repo.NewMockRepo()
	mockCache := cache.NewMockCache()
	mockAuthenticator := auth.NewMockAuthenticator()

	return &application{
		logger:        logger,
		repo:          mockRepo,
		cache:         mockCache,
		authenticator: mockAuthenticator,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d but got %d", expected, actual)
	}
}
