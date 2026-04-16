# Plan: DevHealth — GitHub Org Health & Security Dashboard

**TL;DR**: Build a platform tool that scans GitHub repositories for health, security, and best-practice compliance, then surfaces results in a web dashboard. Go backend + TypeScript frontend, deployed on Azure. This directly mirrors what the Sandvik Developer Experience team builds.

**Why this project is the right fit:**
- It IS a developer-experience platform tool — the exact kind of thing the team builds
- Natural Go/TypeScript split (Go backend, TypeScript dashboard)
- Heavy GitHub API usage (the job lists GitHub as a bonus)
- Security-focused scanning (branch protection, secret scanning, vulnerability alerts)
- Azure deployment demonstrates cloud knowledge
- Very discussable in an interview

---

## Architecture

```
GitHub API  ←→  Go Backend (API + Scanner)  ←→  SQLite  ←→  TypeScript Dashboard (React)
                       ↑                                            ↑
               Azure Container App                        Azure Static Web Apps
```

---

## Phase 1: Go Backend — Core Scanner + API (Days 1–4)

1. Set up Go module, project structure, Dockerfile
2. Integrate GitHub REST API (`google/go-github` library) — authenticate via PAT, list repos
3. Implement health checks per repo:
   - Branch protection on default branch
   - Secret scanning enabled
   - Dependabot alerts enabled
   - CI/CD pipeline exists (`.github/workflows/`)
   - README, LICENSE, CODEOWNERS files exist
   - Open vulnerability alert count
   - Stale repo detection (last commit age)
4. Health scoring engine — 0–100 score per repo, weighted by severity
5. REST API: `POST /api/scan`, `GET /api/repos`, `GET /api/repos/:owner/:name`, `GET /api/summary`
6. SQLite persistence for scan history
7. Unit tests for scoring logic and API handlers

## Phase 2: TypeScript Dashboard (Days 5–8)

8. Vite + React + TypeScript + Tailwind CSS setup
9. Org overview dashboard — average health score, score distribution chart (Recharts)
10. Repo list — sortable/filterable table, color-coded health scores (red/yellow/green)
11. Repo detail view — pass/fail per check with actionable recommendations
12. Scan trigger button with progress indication

## Phase 3: Azure Deployment + Polish (Days 9–12)

13. Docker Compose for local dev
14. Deploy Go backend → Azure Container Apps (free tier: 2M req/month)
15. Deploy dashboard → Azure Static Web Apps (free tier)
16. Professional README with architecture diagram, screenshots, setup instructions
17. GitHub Actions CI/CD — auto-build, test, deploy on push (demonstrates GitHub Actions knowledge)

**Stretch goals** (if time): webhook listener for real-time scanning, health trend charts, GitHub Enterprise support (relevant for Sandvik)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, `chi` router, `go-github`, `go-sqlite3`, `testify` |
| Frontend | React 18, Vite, TypeScript, Tailwind CSS, Recharts |
| Cloud | Azure Container Apps, Azure Static Web Apps |
| CI/CD | GitHub Actions |
| Local dev | Docker Compose |

---

## How This Strengthens Your Application

- In your cover letter: *"To prepare for this role, I built a developer-experience platform tool using Go, TypeScript, and Azure that scans GitHub organizations for security and health compliance."*
- In the interview: discuss trade-offs (why Go for backend, why these health checks matter, how you'd extend it for enterprise use at Sandvik)
- Link the GitHub repo directly in your application
- It proves you can build exactly what the Developer Experience team builds

---

## Decisions

- **SQLite over PostgreSQL** — simpler for demo, portable, mention PostgreSQL-readiness in interview
- **Vite + React over Next.js** — lighter for a pure SPA dashboard, simpler static deployment
- **chi router** — idiomatic Go, close to stdlib, shows Go philosophy
- **GitHub REST API over GraphQL** — simpler to start, sufficient for these checks

---

## Project Structure

```
devhealth/
├── backend/                    # Go backend
│   ├── cmd/server/main.go     # Entry point
│   ├── internal/
│   │   ├── api/               # HTTP handlers
│   │   ├── scanner/           # GitHub scanning logic
│   │   ├── scoring/           # Health score calculation
│   │   ├── store/             # SQLite persistence
│   │   └── models/            # Data types
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── dashboard/                  # TypeScript frontend
│   ├── src/
│   │   ├── components/        # React components
│   │   ├── pages/             # Dashboard, RepoDetail
│   │   ├── api/               # API client
│   │   └── types/             # TypeScript types
│   ├── Dockerfile
│   ├── package.json
│   └── tsconfig.json
├── .github/workflows/          # CI/CD
│   ├── backend.yml
│   └── dashboard.yml
├── docker-compose.yml
└── README.md
```

---

## Verification

1. **Go tests pass**: `go test ./...` — scanner, scoring, and API handler tests
2. **TypeScript builds**: `npm run build` with zero errors
3. **Docker Compose works**: `docker-compose up` starts both services, dashboard connects to API
4. **Scan completes**: Scan a real GitHub org (e.g., your own repos), verify scores make sense
5. **Azure deployment**: Both services accessible via public URLs
6. **GitHub Actions green**: CI pipeline runs tests and builds on push
7. **README is clear**: Someone can clone, run locally, and understand the project in 5 minutes
