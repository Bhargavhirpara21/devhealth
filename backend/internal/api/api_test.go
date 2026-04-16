package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BhargavHirpara/devhealth/internal/models"
)

func TestHandleHealth(t *testing.T) {
	srv := &Server{}
	srv.setupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body["status"])
	}
}

func TestHandleScan_MissingOwner(t *testing.T) {
	srv := &Server{}
	srv.setupRoutes()

	body, _ := json.Marshal(models.ScanRequest{Owner: "", Type: "user"})
	req := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleScan_InvalidType(t *testing.T) {
	srv := &Server{}
	srv.setupRoutes()

	body, _ := json.Marshal(models.ScanRequest{Owner: "test", Type: "invalid"})
	req := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleListRepos_MissingOwner(t *testing.T) {
	srv := &Server{}
	srv.setupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/repos", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleSummary_MissingOwner(t *testing.T) {
	srv := &Server{}
	srv.setupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/summary", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
