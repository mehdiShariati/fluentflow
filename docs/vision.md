# Vision & story

Most language products still optimize for **taps and streaks**, not for the uncomfortable part of fluency: **speaking out loud**, often, with feedback. Cheap bandwidth and capable models changed what is possible—but **shipping** voice-native tutoring still means **media**, **latency**, **reliability**, and **measurement**, not prompts alone.

FluentFlow is a **reference implementation** of that idea end to end: scenario-based sessions, a **join token that dispatches a named agent** into the room, durable session state, transcript-aware feedback, and **instrumentation** you can actually run (health, Prometheus metrics, product events).

### For anyone building a language AI product

You do **not** have to treat this as a read-only demo. The intended use includes: **fork**, change **scenarios** and **tutor prompts**, connect your own **domains** and **keys**, and run in **production** on your cloud. The stack is split so teams can own **API**, **realtime**, and **agents** independently as they grow. A practical checklist lives in **[Build your own](build-your-own.md)**.

Helpful background: **[article on Medium](https://medium.com/p/4d894c404772?postPublishedType=initial)**, **[LiveKit on GitHub](https://github.com/livekit)** (WebRTC + agents ecosystem), and **[Building AI Voice Agents for Production](https://learn.deeplearning.ai/courses/building-ai-voice-agents-for-production/information)** (DeepLearning.AI, voice pipelines and production patterns).

### What we optimize for

- **Speaking-first UX** — real-time voice in the room, not a text chat dressed as a tutor.
- **A thin live path** — realtime stays realtime; heavier analysis can follow the session.
- **Explicit systems design** — where dispatch happens, how experiments snap to sessions, what gets stored and why.

### Scalability, durability, and system design

This stack is shaped so you can **reason about growth** without rewriting the domain model:

- **Separation of concerns** — HTTP/API state is **stateless**; **PostgreSQL** is the durable source of truth for users, sessions, events, transcripts, and feedback. Realtime media and agent work live in **LiveKit + workers**, which scale on a different axis than the API.
- **Durability** — session lifecycle and analytics events are **persisted**, not assumed to live only in memory or a single process. Migrations are **versioned** (`internal/migrate`) so schema evolution stays controlled.
- **Operability at scale** — **health** (`/healthz`) and **Prometheus** (`/metrics`) are first-class so you can load-balance replicas, scrape metrics, and alert before users do. Product events map to **counters** you can monitor as usage grows.
- **Clear next steps** — horizontal API replicas, PgBouncer, read replicas, LiveKit clustering or cloud, more agent workers, and async jobs are **documented paths** ([Scaling](scaling.md), [Monitoring](monitoring.md)), not afterthoughts.

### Scope (honest)

This repository is **not** a hosted product or a funding pitch. It is a **working stack** you can run locally or extend: Docker Compose, Go API, Next.js client, LiveKit, Python agent. What is *not* here is spelled out in the [PRD](prd.md)—OAuth breadth, billing, multi-region SFU operations, full OpenTelemetry, and so on—on purpose, so the core story stays legible.

### About the author & opportunities

I care about **AI-native learning** and **production-grade realtime systems**. I am **open to senior/staff roles** in platform engineering, AI product engineering, or realtime/infra-adjacent teams—roles where this kind of problem shows up for real.

**Connect:** replace these with your own links.

- LinkedIn: `https://www.linkedin.com/in/YOUR_PROFILE`
- Email: `you@example.com` (optional)https://github.com/mehdiShariati/fluentflow

---

*If you are reading the repo on GitHub, the technical entry point is the root [`README.md`](https://github.com/mehdi/fluentflow/blob/main/README.md).*
