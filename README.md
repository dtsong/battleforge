# BattleForge

BattleForge is a full-stack platform for analyzing competitive Pokémon gameplay.

## Overview

Initial focus:

- **Pokémon Showdown** replays (VGC formats, e.g. Regulation H)
- **Pokémon TCG Live** game exports (later)
- Surfacing insights to help players improve decision-making, sequencing, and overall match performance

The goal is to act as a **replay/intel layer** that can power player coaching, practice review, and performance tracking.

## Tech Stack

### Backend

- Go (1.22+)
- HTTP API (`/api/showdown/analyze`, `/api/tcglive/analyze`)
- Showdown `.log` parsing
- Analysis engine that produces a structured `BattleSummary`

### Frontend

- Next.js (App Router) + React + TypeScript
- Tailwind CSS for rapid UI iteration
- React Compiler enabled
- Pages for:
  - `/` – Landing / overview
  - `/showdown` – Showdown replay analysis
  - `/tcg-live` – TCG Live analysis (planned)

### Infrastructure

- Terraform
- Google Cloud Platform
- Cloud Run for the BattleForge API service (planned)
- Remote state in GCS (per environment)

### Tooling

- `pre-commit` for repo-wide quality gates
- Go formatting & linting
- Next.js lint + TS typecheck
- Terraform fmt / validate / tflint

## Repository Layout

```text
battleforge/
  backend/                # Go API + analysis engine
    cmd/
      battleforge-api/    # main HTTP server
        main.go
    internal/
      httpapi/            # routing + HTTP handlers
      showdown/           # Showdown .log fetch + parse
      analysis/           # Battle analysis & BattleSummary model
      config/             # env/config loading
      observability/      # logging, metrics (later)
    go.mod
    go.sum

  frontend/               # Next.js (React + TS) app
    app/
      layout.tsx
      page.tsx
      showdown/
        page.tsx
      tcg-live/
        page.tsx          # stub for future TCG Live UI
    components/
      showdown/
      ui/
      charts/
    lib/
      api.ts              # calls backend API
      showdownTypes.ts    # TS models for BattleSummary
    package.json
    tsconfig.json
    next.config.mjs

  infra/
    terraform/
      envs/
        dev/
          main.tf
          providers.tf
          backend.tf
          variables.tf
          outputs.tf
      modules/
        cloud_run_service/
          main.tf
          variables.tf
          outputs.tf

  .pre-commit-config.yaml
  README.md
```

## Local Development

### Prerequisites

- Go 1.22+
- Node.js 20 LTS
- npm (or pnpm/yarn if you prefer)
- Terraform 1.7+ (for infra work)
- pre-commit installed globally

Example (macOS):

```bash
brew install go node terraform pre-commit
```

### Backend (Go API)

From repo root:

```bash
cd backend
go test ./...
go run ./cmd/battleforge-api
```

By default the API listens on http://localhost:8080.

Health check:

```bash
curl http://localhost:8080/healthz
# -> ok
```

### Frontend (Next.js)

In another terminal:

```bash
cd frontend
npm install
npm run dev
```

Open:

- http://localhost:3000/ for the landing page
- http://localhost:3000/showdown for the Showdown analysis UI

The Showdown page submits replay info to http://localhost:8080/api/showdown/analyze.

### Pre-commit

This repo uses [pre-commit](https://pre-commit.com/) to run checks before each commit:

**Basic hygiene:** whitespace, EOF, YAML/JSON checks

**Go:**
- gofmt, goimports
- golangci-lint
- go test ./...

**Frontend:**
- npm run lint
- npm run typecheck

**Terraform:**
- terraform fmt
- terraform validate
- tflint

Install hooks after cloning:

```bash
pre-commit install
pre-commit run --all-files
```

If a hook makes changes (e.g., gofmt, terraform fmt), re-add and commit.

### Terraform / Infrastructure (Dev)

The Terraform configuration under infra/terraform/envs/dev is intended to:

- Configure the Google provider
- Store state in a GCS bucket
- Deploy the battleforge-api Docker image to Cloud Run via the cloud_run_service module

Basic flow (once you've created a dev project and GCS bucket):

```bash
cd infra/terraform/envs/dev

# set env vars or use a tfvars file:
export TF_VAR_project_id="your-gcp-project-id"
export TF_VAR_battleforge_api_image="gcr.io/your-project/battleforge-api:latest"

terraform init
terraform plan
terraform apply
``````

Future enhancements:

- Staging/prod environments
- Custom domain + HTTPS
- IAM tightening
- Observability (logging/metrics exports)

## Roadmap

- Implement Showdown .log parser (internal/showdown)
- Implement Showdown analysis engine (internal/analysis)
- Define and expose BattleSummary JSON for the frontend
- Build Showdown dashboard: summary card, damage chart, turn timeline, key moments
- Add AI coaching layer over BattleSummary
- Implement TCG Live export parser + analysis
- Wire Terraform module to deploy battleforge-api to Cloud Run
