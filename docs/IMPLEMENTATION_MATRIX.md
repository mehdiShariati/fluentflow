# PRD → implementation matrix (FluentFlow MVP)

Cross-reference for reviewers: [`prd.md`](prd.md) section → what exists in this repository.

| PRD area | § | Implementation |
|----------|---|------------------|
| Onboarding / profile | 12.1, 14.2 | `PUT/GET /v1/me/profile`, Next.js `/onboarding` |
| Scenario catalog | 12.2 | `GET /v1/scenarios`, `internal/api/scenarios.go` |
| LiveKit session | 12.3, 14.3 | `POST /v1/sessions`, join token + `roomConfig` in `internal/livekit/token.go`, Compose `livekit` service |
| Session list UX | 12.2 | `GET /v1/sessions` returns `scenario_title` per row; `limit` clamped (default 20, max 100) |
| Voice agent | 12.3, 14.4 | `agent/tutor_agent.py`, OpenAI Realtime + Silero VAD |
| Post-session feedback | 12.4, 14.5 | `POST …/feedback/generate`, OpenAI JSON or stub, `transcript_summary`, `generation_source` (`openai` / `stub`), `correction_generated` event + Prometheus |
| Progress dashboard | 12.6 | `GET /v1/dashboard/summary`, Next.js `/dashboard` (also shows `learning_metric_snapshots`) |
| Learning history snapshots | 12.6 | `learning_metric_snapshots`, `AppendLearningMetricSnapshot` on `POST …/complete`, `GET /v1/me/learning-snapshots` |
| Account deletion | 14.x / GDPR-ish | `DELETE /v1/me/account`, Next.js `/settings` |
| Scenario deep-link | 12.2 | `/scenarios?highlight=` scroll + ring; recommendation-click validates catalog |
| Experiments | 12.7, 14.7 | `experiments`, `experiment_assignments`, `GET /v1/experiments`, snapshot on session |
| Feature flags | 12.7 | `feature_flags`, `GET /v1/feature-flags`, `PATCH /internal/v1/feature-flags/{key}` |
| Event taxonomy | 14.6 | `internal/analytics/events.go`, `POST /v1/sessions/{id}/events`, `GET …/events` (read-back), `POST …/recommendation-click`, Prometheus `fluentflow_session_events_ingested_total` |
| Transcript chunks | 14.4 | `transcript_segments`, `POST/GET …/transcript`, emits `transcript_generated` |
| Feedback viewed | 14.6 | `POST …/feedback/viewed` → `feedback_viewed` |
| Monitoring | 14.8, §17 | `/healthz`, `/metrics`, internal `…/metrics/summary` (event counts) |
| Guest mode | 14.1 | `POST /v1/auth/guest`, `users.is_guest` |
| Auth | 14.1 | Register/login, JWT, bcrypt; OAuth not implemented |
| Admin / ops | 12.8 | `GET /internal/v1/overview`, `…/experiments`, `…/metrics/summary` (+ flag patch) |
| OAuth | 14.1 | **Not in MVP** (documented gap) |
| Redis / queue workers | 16.5, 16.6 | **Not in MVP** (Postgres-only hot path) |
| Full OpenTelemetry | §17 | **Not in MVP** (metrics + structured event stream) |

## Non-goals (PRD §10) still respected

No multiplayer classrooms, tutor marketplace, social feed, or native app in this codebase; web + API + agent only.
