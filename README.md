# 🛡️ DevHealth — GitHub Org Health & Security Dashboard

## 🚀 Live Demo

- Dashboard: [https://devhealth.vercel.app](https://devhealth.vercel.app)
- API Health: [https://devhealth.onrender.com/api/health](https://devhealth.onrender.com/api/health)

A developer-experience platform tool that scans GitHub repositories for health, security, and best-practice compliance — then surfaces results in a clean web dashboard.

Built with **Go** (backend API + scanner), **TypeScript/React** (dashboard), and designed for **Azure** deployment.

---

## Architecture

```
┌─────────────┐     ┌──────────────────────────────┐     ┌──────────────────────┐
│  GitHub API  │◄───►│  Go Backend (API + Scanner)  │◄───►│  TypeScript Dashboard │
└─────────────┘     │  • Health checks (9 checks)  │     │  • Score overview     │
                    │  • Weighted scoring (0-100)   │     │  • Charts & tables    │
                    │  • SQLite persistence         │     │  • Repo detail view   │
                    │  • REST API (chi router)      │     │  • Vite + React       │
                    └──────────────────────────────┘     └──────────────────────┘
                              ↕                                    ↕
                    Azure Container Apps              Azure Static Web Apps
```

## Health Checks

| Check | Severity | Weight | What it verifies |
|-------|----------|--------|-----------------|
| Branch Protection | Critical | 20 | Default branch has protection rules |
| Secret Scanning | Critical | 15 | Secret scanning is enabled |
| Dependabot Alerts | High | 15 | Dependabot is enabled for vulnerability alerts |
| CI/CD Pipeline | High | 15 | GitHub Actions workflows exist |
| Vulnerabilities | High | 10 | No open vulnerability alerts |
| README | Medium | 10 | README.md exists in repo root |
| License | Low | 5 | License file is present |
| CODEOWNERS | Low | 5 | CODEOWNERS file exists |
| Repository Activity | Low | 5 | Repo has been pushed to within 180 days |

**Total: 100 points**

## Quick Start

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 22+](https://nodejs.org/)
- A [GitHub Personal Access Token](https://github.com/settings/tokens) with `repo` and `read:org` scopes

### Run Locally

**1. Start the backend:**

```bash
cd backend
export GITHUB_TOKEN="ghp_your_token_here"
go run ./cmd/server
# Server starts on http://localhost:8080
```

**2. Start the dashboard:**

```bash
cd dashboard
npm install
npm run dev
# Dashboard starts on http://localhost:5173
```

**3. Open** http://localhost:5173, enter a GitHub username or org, and click **Scan**.

### Run with Docker Compose

```bash
export GITHUB_TOKEN="ghp_your_token_here"
docker compose up --build
# Backend: http://localhost:8080
# Dashboard: http://localhost:3000
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/health` | API health check |
| `POST` | `/api/scan` | Trigger a scan for an owner |
| `GET` | `/api/repos?owner=` | List scanned repos with scores |
| `GET` | `/api/repos/:owner/:repo` | Detailed health report for a repo |
| `GET` | `/api/summary?owner=` | Aggregate stats for an owner |

### Example: Trigger a scan

```bash
curl -X POST http://localhost:8080/api/scan \
  -H "Content-Type: application/json" \
  -d '{"owner": "facebook", "type": "org"}'
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, chi router, go-github, modernc.org/sqlite |
| Frontend | React 18, TypeScript, Vite, Tailwind CSS, Recharts |
| CI/CD | GitHub Actions |
| Containerization | Docker, Docker Compose |
| Cloud | Render (backend), Vercel (dashboard) |

## Project Structure

```
devhealth/
├── backend/
│   ├── cmd/server/main.go          # Entry point, graceful shutdown
│   ├── internal/
│   │   ├── api/                     # REST API handlers (chi router)
│   │   ├── scanner/                 # GitHub health check scanner
│   │   ├── scoring/                 # Weighted scoring engine
│   │   ├── store/                   # SQLite persistence layer
│   │   └── models/                  # Shared data types
│   └── Dockerfile
├── dashboard/
│   ├── src/
│   │   ├── api/                     # Backend API client
│   │   ├── components/              # Reusable UI components
│   │   ├── pages/                   # Dashboard + RepoDetail pages
│   │   └── types/                   # TypeScript type definitions
│   └── Dockerfile
├── .github/workflows/               # CI pipelines
├── docker-compose.yml
└── README.md
```

## Running Tests

```bash
# Backend unit tests
cd backend
go test ./... -v

# Dashboard type check + build
cd dashboard
npx tsc --noEmit
npm run build
```

## License

MIT
