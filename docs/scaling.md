# Scaling

This guide describes how to scale **FluentFlow** components safely: the **Go API**, **PostgreSQL**, **LiveKit**, **Next.js**, and **Python agent workers**, plus cross-cutting concerns (sessions, cost, and backpressure).

## Architecture recap

- **Stateless API:** The Go service holds no in-memory session state required for correctness; JWTs and Postgres back persistence.
- **Stateful realtime:** LiveKit is the SFU; scaling it is different from scaling HTTP replicas.
- **Agents:** Workers connect to LiveKit and OpenAI; scale horizontally with more processes or pods.

## API tier (Go)

### Horizontal scaling

- Run **multiple API replicas** behind a load balancer (round-robin or least connections).
- Ensure **the same `JWT_SECRET`** across replicas so tokens validate everywhere.
- Keep **`CORS_ORIGINS`** aligned with every frontend origin (including preview URLs if applicable).

### Database connections

- Each replica opens its own pool to Postgres. Use **pgx pool limits** appropriate for `(replicas × max conns per replica) < Postgres max_connections`.
- Under load, prefer **raising the API pool** modestly and **adding replicas** rather than huge per-process connection counts.

### Read scaling (future-friendly)

- The current codebase uses a **single primary** for reads and writes. For heavy read dashboards, introduce **read replicas** and route read-only queries (for example dashboard aggregates) to replicas after verifying replication lag is acceptable.

## PostgreSQL

- **Vertical scaling:** More CPU/RAM and fast SSD for OLTP workloads.
- **Connection pooling:** Use **PgBouncer** (transaction or session mode) in front of Postgres when many API instances connect.
- **Migrations:** Run migrations once per release (`internal/migrate`) before or during a controlled rollout.

## LiveKit (realtime)

### Self-hosted

- LiveKit supports **multi-node** clustering; see [LiveKit self-hosting](https://docs.livekit.io/home/self-hosting/) for Redis-backed coordination and UDP/TCP port planning.
- **Network:** Ensure ICE/TURN is configured for clients behind strict NATs; production almost always needs **TURN** for reliability.

### LiveKit Cloud

- Offloads SFU scaling and global edge routing; update `LIVEKIT_URL` and API keys to match the cloud project.

## Agent workers (Python)

- **Horizontal scale:** Run **multiple agent worker processes** with the same `LIVEKIT_AGENT_NAME` and credentials; LiveKit dispatches room jobs to available workers.
- **Resource limits:** Each session uses CPU for VAD and network for media; size containers with headroom for OpenAI Realtime latency spikes.
- **OpenAI quotas:** Realtime sessions consume **concurrent connection** and **token** limits; monitor OpenAI usage dashboards and set alerts.

## Next.js (web)

- **Static optimization:** `next build` produces an optimized server; scale with more Node processes or containers behind a CDN for static assets.
- **`NEXT_PUBLIC_API_URL`:** Must be the **browser-reachable** API base URL in production.

## Optional: Redis and async work (extensions)

The MVP keeps most logic synchronous. Common extensions:

- **Redis** for rate limiting, short-lived session hints, or pub/sub between services.
- **Queue workers** (SQS, RabbitMQ, NATS) for **feedback generation** or analytics export so API latency stays bounded under spikes.

## Capacity checklist

1. **Measure** API latency (`fluentflow_http_request_duration_seconds`) and error rates before adding replicas blindly.
2. **Cap** DB connections and add PgBouncer before Postgres hits `max_connections`.
3. **Match** agent worker count to expected **concurrent voice sessions**, not just HTTP RPS.
4. **Plan** LiveKit and TURN capacity for **concurrent publishers/subscribers** and geographic distribution.

## Cost-aware scaling

- **OpenAI Realtime** is priced per minute and model; optimize prompts and session length in product design.
- **Post-session `gpt-4o-mini`** feedback is cheaper than Realtime; keep feedback generation idempotent and batched where possible.

For observability details, see [Monitoring](monitoring.md).
