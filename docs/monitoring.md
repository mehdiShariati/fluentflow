# Monitoring and observability

FluentFlow exposes **Prometheus** metrics and **health** endpoints on the Go API. This page ties them together with **logging**, **dashboards**, and **alerting** so you can operate the whole project coherently.

## Endpoints (Go API)

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | **Liveness** — returns HTTP 200 when the process is running and ready to answer. Use for load balancer health checks and Kubernetes liveness probes. |
| `GET /metrics` | **Prometheus** text exposition format (pull-based scraping). |
| `GET /internal/v1/metrics/summary` | **Admin** aggregate view of `session_events` counts by `event_type` (requires admin auth). |
| `GET /internal/v1/overview` | **Admin** high-level overview (requires admin auth). |

Configure **`ADMIN_TOKEN`** (or Bearer) for internal routes as documented in the root `README.md`.

## Prometheus metrics (names)

The API registers a custom registry with metrics including:

| Metric | Type | Labels | Meaning |
|--------|------|--------|---------|
| `fluentflow_http_requests_total` | Counter | `method`, `path`, `code` | Total HTTP requests by route pattern and status. |
| `fluentflow_http_request_duration_seconds` | Histogram | `method`, `path` | Request latency distribution. |
| `fluentflow_session_events_ingested_total` | Counter | `event_type` | Product/session events accepted by the API (analytics taxonomy). |

Use these names in **PromQL** queries and Grafana panels.

### Example PromQL

```promql
# Request rate by route
sum(rate(fluentflow_http_requests_total[5m])) by (path)

# Error ratio (5xx)
sum(rate(fluentflow_http_requests_total{code=~"5.."}[5m]))
/
sum(rate(fluentflow_http_requests_total[5m]))

# p95 latency
histogram_quantile(0.95,
  sum(rate(fluentflow_http_request_duration_seconds_bucket[5m])) by (le, path)
)
```

## Scraping configuration

### Prometheus `scrape_configs` snippet

```yaml
scrape_configs:
  - job_name: fluentflow-api
    metrics_path: /metrics
    static_configs:
      - targets: ["api:8080"]   # replace with your service DNS or IP
    relabel_configs:
      - source_labels: [__address__]
        target_label: instance
```

In Kubernetes, prefer **PodMonitor** or **ServiceMonitor** (if using Prometheus Operator) with a selector on the API Service.

## Grafana

1. Add **Prometheus** as a data source pointing at your Prometheus server.
2. Create a dashboard with:
   - **Stat** or **Graph** panels for `fluentflow_http_requests_total` rate.
   - **Heatmap** or **Graph** for latency histograms.
   - **Bar gauge** for `fluentflow_session_events_ingested_total` by `event_type`.
3. Optional: import a **Node** or **Postgres** dashboard for infrastructure correlation (separate exporters).

## Alerting (recommended rules)

Define alerts in Prometheus or Grafana Alerting, for example:

- **High error rate:** `sum(rate(fluentflow_http_requests_total{code=~"5.."}[5m])) / sum(rate(fluentflow_http_requests_total[5m])) > 0.05` for 5 minutes.
- **Latency SLO breach:** p95 `fluentflow_http_request_duration_seconds` above a threshold for sustained periods.
- **Target down:** scrape failures on the `fluentflow-api` job.

Route notifications to **Slack**, **PagerDuty**, or email.

## Logs

- **API:** Structured logs from the Go process (stdout in containers) — aggregate with **Loki**, **CloudWatch**, **Datadog**, or **ELK**.
- **Agent:** Python worker logs for LiveKit connection and OpenAI errors.
- **LiveKit:** Server logs for media and signaling issues.

Correlate logs with **request IDs** where middleware attaches them.

## Traces (extension)

The MVP does not ship OpenTelemetry. A natural upgrade is **OTLP** traces from the API and optionally the agent for distributed tracing (Jaeger, Tempo, Honeycomb).

## Product analytics

Session events stored in Postgres and counted in Prometheus complement each other:

- **Database:** Drill-down per user and session history.
- **Metrics:** Real-time rates and SLOs.

Use **`/internal/v1/metrics/summary`** for operator snapshots and **`/metrics`** for time-series monitoring.

## Security note

- **`/metrics`** and **`/healthz`** are often exposed inside the cluster only; if exposed publicly, protect **`/metrics`** with network policy or authentication to avoid leaking internal labels.

For deployment context, see [Deployment](deployment.md). For capacity planning, see [Scaling](scaling.md).
