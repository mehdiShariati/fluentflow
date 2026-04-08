# FluentFlow documentation

Welcome to the **FluentFlow** documentation — a real-time, speaking-first language learning platform built with **Go**, **Next.js**, **LiveKit**, and a **Python** voice agent (OpenAI Realtime). The architecture emphasizes **durable state in Postgres**, a **stateless API** you can replicate, **separate scaling** for realtime media and agents, and **observability** (metrics, health) for running the system seriously.

## Start here

| Guide | What you will learn |
|-------|---------------------|
| [Vision & story](vision.md) | Why the project exists, scalability/durability framing, and scope boundaries. |
| [Getting started](getting-started.md) | Run the stack locally with Docker or bare metal, environment variables, and verification commands. |
| [Deployment](deployment.md) | Publish the repo on GitHub, enable **GitHub Pages** for this documentation site, and deploy the application to production. |
| [Scaling](scaling.md) | Scale the API, database, LiveKit, and agent workers; connection pooling; queues; and cost-aware capacity planning. |
| [Monitoring](monitoring.md) | Prometheus metrics, health checks, admin endpoints, Grafana dashboards, and alerting patterns. |

## Product and engineering reference

- [Product requirements (PRD)](prd.md) — full product and systems specification.
- [Implementation matrix](IMPLEMENTATION_MATRIX.md) — PRD coverage mapped to the codebase.

## Repository

The canonical project README with API tables, architecture diagrams, and repository layout lives in the GitHub repository root: **`README.md`**.
