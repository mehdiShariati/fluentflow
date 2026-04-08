# Deploying FluentFlow (production)

FluentFlow is a **multi-service** application. In production you run these parts together (or split across hosts), with **TLS** for the browser and **WSS** for LiveKit:

| Component | Role |
|-----------|------|
| **PostgreSQL** | Durable storage: users, sessions, events, transcripts, feedback, experiments |
| **LiveKit server** | WebRTC SFU: learner and agent audio in the same room |
| **Go API** | Auth, session lifecycle, join tokens (with agent dispatch), metrics, feedback |
| **Next.js web** | Learner UI; must call the **public** API URL |
| **Python agent** | `livekit-agents` worker: voice tutor (OpenAI Realtime + VAD when configured) |

Local development uses **`docker compose`** with [`docker-compose.yml`](https://github.com/mehdiShariati/fluentflow/blob/main/docker-compose.yml) as a reference for **service names, ports, and environment wiring**. Production uses the same **images** (`Dockerfile.api`, `web/Dockerfile`, `agent/Dockerfile`) on your own infrastructure.

---

## Environment variables (must stay consistent)

| Area | Purpose |
|------|---------|
| **API** | `DATABASE_URL`, `JWT_SECRET`, `CORS_ORIGINS` (your real web origins), `LIVEKIT_URL` (browser-reachable **`wss://`** in production), `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`, `LIVEKIT_AGENT_NAME` (must match the agent worker), `OPENAI_API_KEY` for feedback/tools, optional `ADMIN_TOKEN` |
| **Web** | `NEXT_PUBLIC_API_URL` — full public base URL of the API (e.g. `https://api.example.com`) |
| **Agent** | Same `LIVEKIT_*` credentials as the API; `LIVEKIT_URL` is often an **internal** URL to LiveKit inside your network; `LIVEKIT_AGENT_NAME` identical to dispatch in tokens; `OPENAI_*` for realtime speech |

Never commit secrets; inject them from your host’s secret store.

---

## HTTPS, WebRTC, and CORS

- Serve the **web app over HTTPS**.
- Point **`LIVEKIT_URL`** at a **`wss://`** endpoint learners can reach.
- Configure **CORS** on the API for your real web origin(s).
- Production networks often need **TURN** for WebRTC behind strict NATs; plan ICE/TURN with your LiveKit deployment ([LiveKit docs](https://docs.livekit.io/)).

---

## LiveKit: self-hosted vs managed

- **Source & community:** [github.com/livekit](https://github.com/livekit) — server, [`agents`](https://github.com/livekit/agents), SDKs, and related tools.
- **Self-hosted:** run `livekit-server` with a proper config, clustering if needed, and correct **UDP/TCP** exposure behind your firewall or cloud SGs.
- **Managed:** [LiveKit Cloud](https://livekit.io/cloud) — point `LIVEKIT_URL` and API keys at the cloud project; reduces SFU operations work.

For voice-agent architecture and latency, **[Building AI Voice Agents for Production](https://learn.deeplearning.ai/courses/building-ai-voice-agents-for-production/information)** (DeepLearning.AI × LiveKit) complements this stack.

---

## Database

- Use **managed PostgreSQL** (TLS, backups) in production.
- Run **migrations** once per release (`internal/migrate` in the API binary) before or during rollout.

---

## Agent workers

- Build from **`agent/Dockerfile`**.
- Run **multiple replicas** with the same **`LIVEKIT_AGENT_NAME`** so LiveKit can distribute room jobs.
- Ensure **`OPENAI_API_KEY`** (and model env vars) are set if you want real voice.

---

## Hosting patterns (typical)

- **Containers:** Kubernetes, ECS, Nomad, Docker Swarm — one deployment per service above.
- **CI/CD:** Build and push images to a registry; deploy with your platform; run migrations as a job.

GitHub hosts **this documentation site** (static HTML) only; it does **not** run LiveKit, Postgres, or your API. See the repository **README** for how the docs are built.

---

## See also

- [Getting started](getting-started.md) — local Docker Compose
- [Build your own](build-your-own.md) — fork, customize, then deploy
- [Scaling](scaling.md) — capacity and components
- [Monitoring](monitoring.md) — `/healthz`, `/metrics`, operations
