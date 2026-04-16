package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"github.com/BhargavHirpara/devhealth/internal/models"
)

// Store handles SQLite persistence.
type Store struct {
	db *sql.DB
}

// New creates a new Store and initializes the schema.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return s, nil
}

func (s *Store) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS repo_health (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		owner TEXT NOT NULL,
		repo TEXT NOT NULL,
		full_name TEXT NOT NULL,
		url TEXT NOT NULL,
		score INTEGER NOT NULL,
		checks_json TEXT NOT NULL,
		scanned_at DATETIME NOT NULL,
		default_branch TEXT NOT NULL,
		last_commit_at DATETIME,
		is_archived BOOLEAN NOT NULL DEFAULT FALSE,
		is_fork BOOLEAN NOT NULL DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_repo_health_owner ON repo_health(owner);
	CREATE INDEX IF NOT EXISTS idx_repo_health_full_name ON repo_health(full_name);
	CREATE INDEX IF NOT EXISTS idx_repo_health_scanned_at ON repo_health(scanned_at);
	`
	_, err := s.db.Exec(query)
	return err
}

// SaveRepoHealth persists a scan result. If a result already exists for the same
// owner/repo, it updates it; otherwise it inserts a new row.
func (s *Store) SaveRepoHealth(rh *models.RepoHealth) error {
	checksJSON, err := json.Marshal(rh.Checks)
	if err != nil {
		return fmt.Errorf("marshaling checks: %w", err)
	}

	query := `
	INSERT INTO repo_health (owner, repo, full_name, url, score, checks_json, scanned_at, default_branch, last_commit_at, is_archived, is_fork)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		score = excluded.score,
		checks_json = excluded.checks_json,
		scanned_at = excluded.scanned_at,
		last_commit_at = excluded.last_commit_at
	`

	result, err := s.db.Exec(query,
		rh.Owner, rh.Repo, rh.FullName, rh.URL, rh.Score,
		string(checksJSON), rh.ScannedAt, rh.DefaultBranch,
		rh.LastCommitAt, rh.IsArchived, rh.IsFork,
	)
	if err != nil {
		return fmt.Errorf("saving repo health: %w", err)
	}

	id, _ := result.LastInsertId()
	rh.ID = id
	return nil
}

// GetReposByOwner returns the latest scan results for all repos belonging to an owner.
func (s *Store) GetReposByOwner(owner string) ([]models.RepoHealth, error) {
	query := `
	SELECT id, owner, repo, full_name, url, score, checks_json, scanned_at, default_branch, last_commit_at, is_archived, is_fork
	FROM repo_health
	WHERE owner = ?
	AND id IN (
		SELECT MAX(id) FROM repo_health WHERE owner = ? GROUP BY full_name
	)
	ORDER BY score ASC
	`
	rows, err := s.db.Query(query, owner, owner)
	if err != nil {
		return nil, fmt.Errorf("querying repos: %w", err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// GetRepo returns the latest scan for a specific repo.
func (s *Store) GetRepo(owner, repo string) (*models.RepoHealth, error) {
	query := `
	SELECT id, owner, repo, full_name, url, score, checks_json, scanned_at, default_branch, last_commit_at, is_archived, is_fork
	FROM repo_health
	WHERE owner = ? AND repo = ?
	ORDER BY scanned_at DESC
	LIMIT 1
	`
	row := s.db.QueryRow(query, owner, repo)

	rh, err := scanRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying repo: %w", err)
	}
	return rh, nil
}

// GetSummary returns aggregate statistics for an owner.
func (s *Store) GetSummary(owner string) (*models.Summary, error) {
	repos, err := s.GetReposByOwner(owner)
	if err != nil {
		return nil, err
	}

	if len(repos) == 0 {
		return nil, nil
	}

	summary := &models.Summary{
		Owner:      owner,
		TotalRepos: len(repos),
		CheckPassRates: make(map[models.CheckName]float64),
	}

	totalScore := 0
	checkPasses := make(map[models.CheckName]int)
	checkTotals := make(map[models.CheckName]int)
	var latestScan time.Time

	for _, r := range repos {
		totalScore += r.Score
		if r.ScannedAt.After(latestScan) {
			latestScan = r.ScannedAt
		}

		switch {
		case r.Score < 40:
			summary.ScoreDistribution.Critical++
		case r.Score < 70:
			summary.ScoreDistribution.Warning++
		case r.Score < 90:
			summary.ScoreDistribution.Good++
		default:
			summary.ScoreDistribution.Excellent++
		}

		for _, c := range r.Checks {
			checkTotals[c.Name]++
			if c.Passed {
				checkPasses[c.Name]++
			}
		}
	}

	summary.AverageScore = float64(totalScore) / float64(len(repos))
	summary.LastScanAt = latestScan

	for name, total := range checkTotals {
		if total > 0 {
			summary.CheckPassRates[name] = float64(checkPasses[name]) / float64(total) * 100
		}
	}

	return summary, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

func scanRows(rows *sql.Rows) ([]models.RepoHealth, error) {
	var results []models.RepoHealth
	for rows.Next() {
		rh, err := scanFromRows(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *rh)
	}
	return results, rows.Err()
}

func scanFromRows(rows *sql.Rows) (*models.RepoHealth, error) {
	var rh models.RepoHealth
	var checksJSON string
	var lastCommitAt sql.NullTime

	err := rows.Scan(
		&rh.ID, &rh.Owner, &rh.Repo, &rh.FullName, &rh.URL,
		&rh.Score, &checksJSON, &rh.ScannedAt, &rh.DefaultBranch,
		&lastCommitAt, &rh.IsArchived, &rh.IsFork,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}

	if lastCommitAt.Valid {
		rh.LastCommitAt = lastCommitAt.Time
	}

	if err := json.Unmarshal([]byte(checksJSON), &rh.Checks); err != nil {
		return nil, fmt.Errorf("unmarshaling checks: %w", err)
	}

	return &rh, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRow(row scannable) (*models.RepoHealth, error) {
	var rh models.RepoHealth
	var checksJSON string
	var lastCommitAt sql.NullTime

	err := row.Scan(
		&rh.ID, &rh.Owner, &rh.Repo, &rh.FullName, &rh.URL,
		&rh.Score, &checksJSON, &rh.ScannedAt, &rh.DefaultBranch,
		&lastCommitAt, &rh.IsArchived, &rh.IsFork,
	)
	if err != nil {
		return nil, err
	}

	if lastCommitAt.Valid {
		rh.LastCommitAt = lastCommitAt.Time
	}

	if err := json.Unmarshal([]byte(checksJSON), &rh.Checks); err != nil {
		return nil, fmt.Errorf("unmarshaling checks: %w", err)
	}

	return &rh, nil
}
