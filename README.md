# FluentFlow — Real-time AI language tutor

[![Go tests](https://github.com/mehdiShariati/fluentflow/actions/workflows/go-test.yml/badge.svg)](https://github.com/mehdiShariati/fluentflow/actions/workflows/go-test.yml)
[![Web](https://github.com/mehdiShariati/fluentflow/actions/workflows/web.yml/badge.svg)](https://github.com/mehdiShariati/fluentflow/actions/workflows/web.yml)
[![Docs](https://github.com/mehdiShariati/fluentflow/actions/workflows/docs.yml/badge.svg)](https://github.com/mehdiShariati/fluentflow/actions/workflows/docs.yml)

## Why FluentFlow

Most language apps optimize for **passive drills**. Real fluency needs **high-frequency spoken production**—safe practice, quick feedback, and a system that does not fall over when you add voice, agents, and analytics.

**FluentFlow** is an end-to-end stack for **speaking-first** learning: learners pick scenarios, join a **LiveKit** room, talk with a **Python voice agent** (OpenAI Realtime when configured), and get **post-session feedback** (OpenAI `gpt-4o-mini` or a deterministic stub). A **Go** API owns auth, profiles, sessions, experiments, events, and **Prometheus** metrics; **Next.js** is the learner UI.

### Use it as your own AI language tutor (starter)

**Anyone** can use this repository to **start** a language-learning product or classroom tool: fork the repo, follow **[docs/getting-started.md](docs/getting-started.md)** to run locally, then **[docs/build-your-own.md](docs/build-your-own.md)** to customize **tutor prompts** (`agent/tutor_agent.py`), **scenarios** (`internal/api/scenarios.go`), UI copy, and deployment. You do not need special permission—add a **license** file if you redistribute your fork.

**System design is intentional:** the API is **stateless** and **horizontally scalable**; **Postgres** holds durable state; the **realtime** path (LiveKit + agent workers) scales separately; **migrations** are embedded; **metrics and health** are built for production operation—not a demo glued to a single machine.

The **problem framing and narrative**—why this exists beyond a feature list—are in **[docs/vision.md](docs/vision.md)** (and [VISION.md](VISION.md) points there). Deep dives: **[docs/scaling.md](docs/scaling.md)**, **[docs/monitoring.md](docs/monitoring.md)**.

---

## At a glance

| Layer | Stack |
|-------|--------|
| Web | Next.js (App Router), LiveKit client |
| API | Go (chi), JWT, Postgres (pgx), `/healthz`, `/metrics` |
| Realtime | LiveKit server, WebRTC |
| Agent | Python (`livekit-agents`), Silero VAD, OpenAI Realtime |
| Data | PostgreSQL — users, sessions, transcripts, feedback, experiments |

| Resource | Link |
|----------|------|
| **Article (Medium)** | [FluentFlow on Medium](https://medium.com/p/4d894c404772?postPublishedType=initial) |
| **Vision & story** | [docs/vision.md](docs/vision.md) |
| **Build your own tutor** (fork, customize, ship) | [docs/build-your-own.md](docs/build-your-own.md) |
| **Scaling & system design** | [docs/scaling.md](docs/scaling.md) |
| **Monitoring & durability signals** | [docs/monitoring.md](docs/monitoring.md) |
| **Online docs** (GitHub Pages) | After you enable Pages: `https://<your-username>.github.io/<repository-name>/` — see [Documentation site (MkDocs)](#documentation-site-mkdocs) |
| **Product / systems PRD** | [docs/prd.md](docs/prd.md) |
| **PRD → implementation map** | [docs/IMPLEMENTATION_MATRIX.md](docs/IMPLEMENTATION_MATRIX.md) |

---

## Table of contents

1. [What this system includes](#what-this-system-includes)
2. [Architecture](#architecture-systems)
3. [Quick start (Docker)](#quick-start-docker)
4. [Tests and verification](#tests-and-verification)
5. [Documentation site (MkDocs)](#documentation-site-mkdocs)
6. [Deployment (production)](#deployment-production)
7. [Scaling](#scaling)
8. [Monitoring](#monitoring)
9. [API surface (v1)](#api-surface-v1)
10. [Key technical decisions](#key-technical-decisions)
11. [Repository layout](#repository-layout)
12. [Open to opportunities](#open-to-opportunities)
13. [License](#license)

*(See also [Build your own AI language tutor](#use-it-as-your-own-ai-language-tutor-starter) above.)*

---

## What this system includes

- **LiveKit end-to-end:** self-hosted `livekit-server` (dev), browser `livekit-client`, join tokens from your API with **`roomConfig` agent dispatch** (see [LiveKit agent dispatch](https://docs.livekit.io/agents/server/agent-dispatch/)).
- **Voice agent:** `livekit-agents` worker with **Silero VAD** + **OpenAI Realtime** (`OPENAI_API_KEY` required for real speech).
- **Session UX:** chat-style transcript, de-duplication, per-message **Translate** / **Analyze** (API-backed with stub fallbacks).
- **Backend:** modular Go service (Chi, JWT, bcrypt, pgx, CORS, `/metrics`, `/healthz`).
- **Data & durability:** Postgres for users, profiles, sessions, events, experiments, flags, feedback, learning snapshots — **durable** lifecycle and audit-friendly event streams, not ephemeral process state.
- **Scale-ready shape:** stateless API tier (add replicas behind a load balancer), separable LiveKit and **horizontal agent workers**; see [Scaling](docs/scaling.md).
- **Product instrumentation:** session event taxonomy, experiment snapshots, feature flags — aligned with **metrics** and future analytics pipelines.
- **CI:** GitHub Actions for Go, web, and documentation builds.

---

## Architecture (systems)

```mermaid
flowchart TB
  subgraph Learner
    WEB[Next.js web]
  end

  subgraph Realtime
    LK[LiveKit server]
    AG[Python agent worker]
  end

  subgraph Platform
    API[Go API / BFF]
    DB[(PostgreSQL)]
  end

  WEB -->|REST + JWT| API
  WEB -->|WebRTC ws://| LK
  API -->|SQL| DB
  AG -->|worker + media| LK
  API -.->|issues join URL + token| WEB
```

### How the LiveKit agent joins

When the learner starts a session, the API returns a **short-lived JWT** whose claims include **room join** permissions and **`roomConfig.agents`**: a `RoomAgentDispatch` for `fluentflow-tutor` plus **JSON metadata** (scenario, level, goals). The browser connects; LiveKit **dispatches a job** to the Python worker, which runs the voice session in that room.

```mermaid
sequenceDiagram
  participant U as Learner browser
  participant W as Next.js
  participant A as Go API
  participant D as Postgres
  participant L as LiveKit server
  participant P as Python agent

  U->>W: Start scenario
  W->>A: POST /v1/sessions (Bearer)
  A->>D: INSERT session, experiments snapshot
  A->>A: Mint JWT video.roomJoin + roomConfig.agents[]
  A-->>W: room name, ws URL, token
  W->>L: Connect WebRTC (learner token)
  L->>P: Dispatch job (agent_name + metadata)
  P->>L: Agent participant joins + publishes audio
  Note over U,P: Mic uplink + tutor downlink (OpenAI Realtime)
```

---

## Quick start (Docker)

**Prerequisites:** Docker Desktop (or compatible), and an **OpenAI API key** for real voice and LLM feedback.

1. Copy env and set your key (recommended):

   ```bash
   cp .env.example .env
   # OPENAI_API_KEY=sk-...
   ```

2. From the repo root:

   ```bash
   docker compose up --build
   ```

   Detached:

   ```bash
   docker compose up -d --build
   ```

3. Open **http://localhost:3000** — register, profile, scenario, **Connect to room**.

| Port | Service |
|------|---------|
| 3000 | Next.js |
| 8080 | Go API |
| 7880 | LiveKit (WS) |
| 5432 | Postgres |

**LiveKit dev defaults:** API key `devkey`, secret `secret` (see [LiveKit local dev](https://docs.livekit.io/home/self-hosting/local/)).

**Disable tutor dispatch (client-only debugging):** set `LIVEKIT_AGENT_NAME=` empty in `.env`.

**Agent Realtime knobs:** `OPENAI_REALTIME_MODEL`, `OPENAI_TRANSCRIPTION_MODEL`, `OPENAI_TTS_VOICE` (see `.env.example`).

**Windows / Docker UDP issues:** Hyper-V can reserve `50000–502xx`. This repo uses [`livekit-docker.yaml`](livekit-docker.yaml) (UDP mux on **7882** only).

**Local dev without Docker:** see [docs/getting-started.md](docs/getting-started.md).

---

## Tests and verification

```bash
make test
```

Full check (Go vet + tests + Next.js lint + production build; Node 20+):

```bash
make verify
```

Avoid `go test ./...` from the repo root if `web/node_modules` contains nested Go packages.

**Learning the codebase:** [teach.md](teach.md).

---

## Documentation site (MkDocs)

The **FluentFlow** technical docs (vision, getting started, **production deployment**, scaling, monitoring) live under [`docs/`](docs/) and are built with **MkDocs Material** ([`mkdocs.yml`](mkdocs.yml)). A successful GitHub Actions run ([`.github/workflows/docs.yml`](.github/workflows/docs.yml)) publishes **static HTML** to **GitHub Pages** — that is **documentation only**, not the running app.

**Enable Pages:** repository **Settings → Pages → Source: GitHub Actions**. **Preview locally:** `pip install -r requirements-docs.txt` then `mkdocs serve`, or `make docs` to build `./site`.

**Guides:**

| Topic | File |
|-------|------|
| Article (Medium) | [medium.com/…/4d894c404772](https://medium.com/p/4d894c404772?postPublishedType=initial) |
| Vision & story | [docs/vision.md](docs/vision.md) |
| Getting started | [docs/getting-started.md](docs/getting-started.md) |
| Build your own (customize & ship) | [docs/build-your-own.md](docs/build-your-own.md) |
| **Production deployment** (FluentFlow stack) | [docs/deployment.md](docs/deployment.md) |
| Scaling | [docs/scaling.md](docs/scaling.md) |
| Monitoring | [docs/monitoring.md](docs/monitoring.md) |

**Forking:** Badge URLs and `repo_url` in `mkdocs.yml` may still point at the upstream module path; adjust for your fork if needed.

---

## Deployment (production)

FluentFlow is **multi-service** (Postgres, LiveKit, API, web, agent). See **[docs/deployment.md](docs/deployment.md)** for environment variables, TLS/WebRTC, LiveKit, database migrations, and agent replicas.

CI in this repo runs **tests** and **builds the docs site** for GitHub Pages; **hosting** the live application is on **your** cloud or servers (not GitHub Pages).

---

## Scaling

The API is **stateless** behind a load balancer; **Postgres**, **LiveKit**, and **agent workers** scale on different axes. Read the full guide: **[docs/scaling.md](docs/scaling.md)** — connection pooling, PgBouncer, LiveKit clustering / cloud, horizontal agent workers, OpenAI quotas, and cost-aware capacity.

---

## Monitoring

The Go API exposes `GET /healthz` (liveness) and `GET /metrics` (Prometheus). Metrics include `fluentflow_http_requests_total`, `fluentflow_http_request_duration_seconds`, and `fluentflow_session_events_ingested_total`. Internal admin routes under `/internal/v1/` require `ADMIN_TOKEN`.

**Full runbook:** **[docs/monitoring.md](docs/monitoring.md)** — scrape config, PromQL examples, Grafana panels, alerting, logs, and trace extensions.

---

## API surface (v1)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/v1/auth/register`, `/v1/auth/login`, `/v1/auth/guest` | JWT; guest → `is_guest` |
| GET | `/v1/me` | |
| DELETE | `/v1/me/account` | JSON `{"password"}` (omit for guest) |
| GET | `/v1/me/learning-snapshots` | Query `limit` |
| GET/PUT | `/v1/me/profile` | |
| GET | `/v1/scenarios`, `/v1/experiments`, `/v1/feature-flags` | |
| GET | `/v1/sessions`, POST `/v1/sessions` | `scenario_title` on items |
| GET | `/v1/sessions/{id}` | |
| POST | `/v1/sessions/{id}/livekit-token` | |
| POST | `/v1/sessions/{id}/events` | |
| POST | `/v1/sessions/{id}/transcript` | |
| GET | `/v1/sessions/{id}/transcript` | |
| POST | `/v1/sessions/{id}/complete` | |
| GET | `/v1/sessions/{id}/feedback` | `generation_source`, `recommended_scenario_title` |
| POST | `/v1/sessions/{id}/feedback/generate` | |
| POST | `/v1/sessions/{id}/feedback/viewed` | |
| POST | `/v1/sessions/{id}/recommendation-click` | |
| GET | `/v1/sessions/{id}/events` | |
| POST | `/v1/ai/translate`, `/v1/ai/analyze` | |
| GET | `/v1/dashboard/summary` | |
| GET | `/internal/v1/overview`, `/internal/v1/experiments`, `/internal/v1/metrics/summary` | Admin |
| PATCH | `/internal/v1/feature-flags/{key}` | Admin |
| GET | `/healthz`, `/metrics` | Ops |

---

## Key technical decisions

1. **Agent dispatch in the join token** — fewer moving parts; dispatch remains explicit via `agent_name` on the worker.
2. **OpenAI at two speeds** — Realtime for live tutoring; `gpt-4o-mini` for structured post-session feedback (cost/latency tradeoff for the async path).
3. **Postgres as source of truth** — sessions and `session_events` for dashboards, export, and **durability** across deploys and replicas.
4. **Stateless API** — JWT + DB back correctness so you can **scale HTTP** independently of LiveKit and agent pools.
5. **Prometheus first** — `/metrics` and `/healthz` for **operational** scaling (SLIs, alerting, load balancer probes); traces can layer on later.

---

## Repository layout

```
cmd/api/           # HTTP entrypoint
internal/api/      # Handlers, middleware
internal/store/    # Postgres access
internal/livekit/  # Join JWT + roomConfig dispatch
internal/openai/   # Post-session + message tools
internal/migrate/  # Embedded SQL
web/               # Next.js App Router UI
agent/             # LiveKit Agents tutor
docs/              # MkDocs source (plus PRD & matrix)
mkdocs.yml         # Documentation site config
.github/workflows/ # CI: tests + MkDocs → GitHub Pages (docs only)
```

---

## Open to opportunities

I am open to **senior/staff** roles in platform engineering, AI product engineering, and realtime / infrastructure-adjacent teams. Replace the placeholders in **[docs/vision.md](docs/vision.md)** with your LinkedIn and email, or add them to your GitHub profile.

---

## License

Use and extend under your own terms; add an explicit license file if you need a standard OSS license.
