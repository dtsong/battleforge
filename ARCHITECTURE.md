# Architecture & Data Design

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Browser / Client                         │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Next.js Frontend (Port 3000)                  │
├─────────────────────────────────────────────────────────────────┤
│  • React Components (Showdown, TCG Live analysis pages)         │
│  • TypeScript Type System (aligned with backend)                │
│  • Tailwind CSS UI                                              │
│  • API Client Layer (calls /api/*)                              │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP/JSON
                         │ /api/showdown/analyze
                         │ /api/tcglive/analyze
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│               Go HTTP API Server (Port 8080)                    │
├─────────────────────────────────────────────────────────────────┤
│  httpapi/                                                        │
│    ├─ router.go         (chi router, request routing)           │
│    ├─ showdown_handlers.go (Showdown analysis endpoints)        │
│    └─ tcglive_handlers.go   (TCG Live analysis endpoints)       │
└────────────────────────┬────────────────────────────────────────┘
                         │
                    ┌────┴────┬─────────────┐
                    ▼         ▼             ▼
        ┌──────────────┐ ┌────────┐ ┌──────────────┐
        │   Parsing    │ │Analysis│ │ Observability│
        │   Engine     │ │Engine  │ │   (Logs)     │
        ├──────────────┤ ├────────┤ └──────────────┘
        │ Showdown     │ │Generate│
        │ .log parser  │ │Battle  │
        │ (internal/)  │ │Summary │
        │              │ │        │
        │ TCG Live     │ │Query   │
        │ export       │ │Stats   │
        │ parser (TBD) │ │        │
        └──────┬───────┘ └────┬───┘
               │              │
               └──────┬───────┘
                      ▼
        ┌──────────────────────────┐
        │   PostgreSQL Database    │
        │  (Cloud SQL on GCP)      │
        └──────────────────────────┘
                      ▲
                      │
        ┌─────────────┴──────────────┐
        │                            │
        ▼                            ▼
   ┌─────────┐              ┌──────────────┐
   │  Battles│              │ Analysis     │
   │ (Logs)  │              │ Results      │
   └─────────┘              └──────────────┘

┌─────────────────────────────────────────────────────────────────┐
│              Infrastructure (Terraform on GCP)                   │
├─────────────────────────────────────────────────────────────────┤
│  • Cloud Run (battleforge-api container)                        │
│  • Cloud SQL (PostgreSQL)                                       │
│  • Cloud Storage (state files, logs)                            │
│  • Artifact Registry (Docker images)                            │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Showdown Replay Analysis Flow

```
User uploads Showdown .log file
           │
           ▼
    ┌─────────────┐
    │ HTTP POST   │
    │ /api/       │
    │ showdown/   │
    │ analyze     │
    └──────┬──────┘
           │
           ▼
   ┌──────────────┐
   │ Showdown     │
   │ .log Parser  │
   │ (internal/)  │
   └──────┬───────┘
           │
           ▼ (parsed battle events)
   ┌──────────────┐
   │ Analysis     │
   │ Engine       │
   │ (internal/)  │
   └──────┬───────┘
           │
           ▼ (BattleSummary struct)
   ┌──────────────┐
   │ Save to DB   │
   └──────┬───────┘
           │
           ▼
   ┌──────────────┐
   │ Return JSON  │
   │ BattleSummary│
   └──────┬───────┘
           │
           ▼ (HTTP 200 + JSON)
   Frontend renders analysis dashboard
   (turn timeline, damage charts, key moments)
```

## Component Responsibility Matrix

| Component | Responsibility | Technology |
|-----------|-----------------|------------|
| **Frontend** | UI/UX, user interactions, visualization | Next.js, React, TypeScript, Tailwind |
| **HTTP Router** | Request routing, input validation | go-chi/chi |
| **Showdown Parser** | Parse .log format into structured events | Go stdlib (strings, parsing) |
| **TCG Live Parser** | Parse game exports (future) | Go stdlib |
| **Analysis Engine** | Process events, calculate stats, generate insights | Go stdlib |
| **Database** | Persist battles, results, user data | PostgreSQL |
| **Observability** | Logging, metrics collection | Go stdlib + (future: Datadog/Prometheus) |
| **Infrastructure** | Deployment, scaling, state management | Terraform, GCP |

## API Contract

### POST /api/showdown/analyze

**Request:**
```json
{
  "battleLog": "string (base64 encoded .log content)",
  "metadata": {
    "format": "Regulation H",
    "uploadedAt": "2025-11-22T10:00:00Z"
  }
}
```

**Response:**
```json
{
  "battlId": "uuid",
  "status": "success",
  "data": { BattleSummary object (see types.go / showdown.ts) }
}
```

### POST /api/tcglive/analyze

Similar structure for TCG Live game exports (planned).
