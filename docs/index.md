# FluentFlow documentation

**FluentFlow** is a real-time, **speaking-first** language learning stack: **Go** API, **Next.js** web, **LiveKit** for voice, and a **Python** agent (OpenAI Realtime). **Anyone** can use this repository as a **starting point** for their own AI language tutor or voice learning product—fork it, change scenarios and prompts, deploy on your infrastructure.

!!! tip "New here?"
    1. **[Getting started](getting-started.md)** — run everything with Docker in a few minutes.  
    2. **[Build your own](build-your-own.md)** — customize tutor behavior, scenarios, and branding.  
    3. **[Vision & story](vision.md)** — design goals and scope.

The architecture keeps **durable data in Postgres**, a **stateless API** you can scale out, and **separate scaling** for LiveKit and agent workers, with **metrics and health** for serious operation.

## Guides for everyone

| Guide | Who it is for |
|-------|----------------|
| [Getting started](getting-started.md) | First run: Docker, ports, env vars, Windows tips |
| [Build your own](build-your-own.md) | Forking, customizing tutor + scenarios, shipping your product |
| [Vision & story](vision.md) | Why this exists, durability and scale framing |
| [Deployment](deployment.md) | Production: services, TLS/WebRTC, LiveKit, database, agents |
| [Scaling](scaling.md) | Capacity: API, Postgres, LiveKit, workers |
| [Monitoring](monitoring.md) | `/healthz`, `/metrics`, operations |

## Deep reference (engineering)

- [Product requirements (PRD)](prd.md) — full product / systems specification.
- [Implementation matrix](IMPLEMENTATION_MATRIX.md) — PRD coverage mapped to the codebase.

## Learn more

- **[FluentFlow article (Medium)](https://medium.com/p/4d894c404772?postPublishedType=initial)** — project story and technical overview on Medium.
- **[LiveKit on GitHub](https://github.com/livekit)** — realtime WebRTC infrastructure and the [`agents`](https://github.com/livekit/agents) framework FluentFlow uses.
- **[Building AI Voice Agents for Production](https://learn.deeplearning.ai/courses/building-ai-voice-agents-for-production/information)** — DeepLearning.AI short course with LiveKit on architecture, latency, and production voice agents.

## Repository root

API tables, diagrams, and layout: **`README.md`** in the GitHub repository.
