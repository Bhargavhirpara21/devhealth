package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/go-github/v62/github"

	"github.com/BhargavHirpara/devhealth/internal/models"
	"github.com/BhargavHirpara/devhealth/internal/scanner"
	"github.com/BhargavHirpara/devhealth/internal/store"
)

// Server holds the API dependencies.
type Server struct {
	router  chi.Router
	store   *store.Store
	scanner *scanner.Scanner
}

// New creates a new API server.
func New(st *store.Store, sc *scanner.Scanner) *Server {
	s := &Server{
		store:   st,
		scanner: sc,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/health", s.handleHealth)
	r.Post("/api/scan", s.handleScan)
	r.Get("/api/repos", s.handleListRepos)
	r.Get("/api/repos/{owner}/{repo}", s.handleGetRepo)
	r.Get("/api/summary", s.handleSummary)

	s.router = r
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req models.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	req.Owner = strings.TrimSpace(req.Owner)
	if req.Owner == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "owner is required",
		})
		return
	}

	if req.Type == "" {
		req.Type = "user"
	}
	if req.Type != "user" && req.Type != "org" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "type must be 'user' or 'org'",
		})
		return
	}

	ctx := r.Context()
	results, err := s.scanner.ScanOwner(ctx, req.Owner, req.Type)
	if err != nil {
		log.Printf("scan error: %v", err)

		if isGitHubNotFound(err) {
			writeJSON(w, http.StatusNotFound, models.ErrorResponse{
				Error:   "Owner not found",
				Details: "The specified GitHub user or organization does not exist.",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Scan failed",
		})
		return
	}

	for i := range results {
		if err := s.store.SaveRepoHealth(&results[i]); err != nil {
			log.Printf("failed to save result for %s: %v", results[i].FullName, err)
		}
	}

	writeJSON(w, http.StatusOK, models.ScanResponse{
		Message:    "Scan completed successfully",
		ReposFound: len(results),
		ScanID:     time.Now().UTC().Format("20060102T150405"),
	})
}

func (s *Server) handleListRepos(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	if owner == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "owner query parameter is required",
		})
		return
	}

	repos, err := s.store.GetReposByOwner(owner)
	if err != nil {
		log.Printf("error fetching repos: %v", err)
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch repos",
		})
		return
	}

	if repos == nil {
		repos = []models.RepoHealth{}
	}

	writeJSON(w, http.StatusOK, repos)
}

func (s *Server) handleGetRepo(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	rh, err := s.store.GetRepo(owner, repo)
	if err != nil {
		log.Printf("error fetching repo: %v", err)
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch repo",
		})
		return
	}

	if rh == nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{
			Error: "Repo not found. Run a scan first.",
		})
		return
	}

	writeJSON(w, http.StatusOK, rh)
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	if owner == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{
			Error: "owner query parameter is required",
		})
		return
	}

	summary, err := s.store.GetSummary(owner)
	if err != nil {
		log.Printf("error fetching summary: %v", err)
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch summary",
		})
		return
	}

	if summary == nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{
			Error: "No scan data found for this owner. Run a scan first.",
		})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func isGitHubNotFound(err error) bool {
	var ghErr *github.ErrorResponse
	if ok := false; !ok {
		_ = ghErr
	}
	return strings.Contains(err.Error(), "404")
}

// NewGitHubClient creates an authenticated GitHub client from a token.
func NewGitHubClient(ctx context.Context, token string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)
	return client
}
