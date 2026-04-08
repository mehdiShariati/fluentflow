# Getting started

This guide covers local setup with **Docker Compose** (recommended) and optional **component-by-component** development.

## Prerequisites

- **Docker Desktop** (or a compatible container runtime) for the one-command stack.
- **OpenAI API key** if you want real voice (OpenAI Realtime) and AI-powered post-session feedback — optional for exploring the UI and stub paths.
- **Node.js 20+** and **Go 1.18+** only if you run services outside Docker (see below).

## One-command run (Docker)

1. Copy the environment template and set secrets:

   ```bash
   cp .env.example .env
   # Set at least OPENAI_API_KEY=sk-... for live speech and LLM feedback
   ```

2. From the repository root:

   ```bash
   docker compose up --build
   ```

   For background services (recommended for daily development):

   ```bash
   docker compose up -d --build
   ```

3. Open the app at **http://localhost:3000** — register (or use guest), complete your profile, pick a scenario, then **Connect to room** and speak when the tutor is present.

### Default ports

| Port | Service |
|------|---------|
| 3000 | Next.js web |
| 8080 | Go API |
| 7880 | LiveKit (WebSocket) |
| 5432 | PostgreSQL |

### LiveKit development credentials

The compose file aligns with common LiveKit local defaults: API key `devkey`, secret `secret`. Match these in your `.env` if you override them.

### Disable agent dispatch

To debug the browser without spawning the tutor worker, set `LIVEKIT_AGENT_NAME=` (empty) in `.env` so join tokens omit `roomConfig.agents`.

### Windows and UDP port conflicts

If Docker reports UDP bind errors on ports in the `50000` range, Windows Hyper-V may reserve that range. This repository ships [`livekit-docker.yaml`](https://github.com/mehdi/fluentflow/blob/main/livekit-docker.yaml) with UDP mux on **7882** only to avoid exposing a wide ephemeral port range.

### Agent environment (Realtime)

The Python worker reads model-related variables (see `.env.example`), for example:

- `OPENAI_REALTIME_MODEL`
- `OPENAI_TRANSCRIPTION_MODEL`
- `OPENAI_TTS_VOICE`

## Quality checks

From the repository root:

```bash
make test
```

Full verification (Go vet, tests, Next.js lint, production build):

```bash
make verify
```

Avoid `go test ./...` from the repo root if `web/node_modules` contains nested Go packages; use `make test` or scope packages under `./cmd/...` and `./internal/...`.

## Local development without full Docker rebuild

- **API:** Run Postgres, then `go run -buildvcs=false ./cmd/api` with `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS`, and `LIVEKIT_*` as documented in `.env.example`.
- **Web:** `cd web`, copy `.env.local.example` to `.env.local`, then `npm run dev`.
- **Agent:** Python venv, `pip install -r requirements.txt`, export `LIVEKIT_*` pointing at your LiveKit instance, run `python tutor_agent.py dev`.

For a guided tour of the codebase, see [`teach.md`](https://github.com/mehdi/fluentflow/blob/main/teach.md) in the repository root.
