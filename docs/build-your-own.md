# Build your own AI language tutor (from this repo)

This project is meant to be **reused**. You can fork it, rename it, and ship your own **speaking-first** language product—startup, classroom tool, internal training, or portfolio—without starting from zero on WebRTC, sessions, or agent dispatch.

You do **not** need permission to use it as a starting point. Add your own **license** file in the repository root if you redistribute code; the upstream README notes educational/portfolio use until you choose something explicit.

---

## What you get out of the box

| Piece | What it does |
|-------|----------------|
| **Web app** | Register, profile, pick a **scenario**, join a **voice room**, live transcript, translate/analyze helpers |
| **API** | Auth, sessions, LiveKit **join tokens** with **agent dispatch**, feedback, metrics, events |
| **Voice agent** | Python worker: listens and speaks via **OpenAI Realtime** when you set `OPENAI_API_KEY` |
| **Database** | Postgres schema for users, sessions, transcripts, feedback, experiments |

If you only want to **try** it locally, follow **[Getting started](getting-started.md)** first (`docker compose up`), then come back here to customize.

**External resources:** explore the **[LiveKit organization on GitHub](https://github.com/livekit)** for servers, SDKs, and examples, and the DeepLearning.AI course **[Building AI Voice Agents for Production](https://learn.deeplearning.ai/courses/building-ai-voice-agents-for-production/information)** for voice pipelines, WebRTC vs HTTP, and scaling agents.

---

## Customize in the right order

### 1. Environment and secrets

- Copy **`.env.example`** → **`.env`** (never commit real keys).
- Set **`OPENAI_API_KEY`** for real voice and post-session feedback.
- Keep **`LIVEKIT_API_KEY`**, **`LIVEKIT_API_SECRET`**, and **`LIVEKIT_AGENT_NAME`** aligned between the **API** and the **agent** service (see [Deployment](deployment.md)).

### 2. Tutor personality and instructions

Edit **`agent/tutor_agent.py`**. The function **`scenario_prompt`** builds the system instructions sent to the model. Change tone (strict vs friendly), languages, length of replies, or add your brand name. The worker reads **job metadata** (from the join token) so scenario and level still flow from the API.

If you rename the agent in code, set **`LIVEKIT_AGENT_NAME`** everywhere to match and ensure the **Go** token code uses the same name when minting JWTs (`internal/livekit`).

### 3. Scenarios (lesson topics)

Scenarios are a **static catalog** in Go: **`internal/api/scenarios.go`**. Edit the **`catalog`** slice: add IDs, titles, descriptions, and levels. Any new **`ID`** must stay in sync with validation (`IsValidScenarioID`) used when creating sessions—so define scenarios **only** in this catalog (or refactor to your own store later).

### 4. Branding and UI

- **Next.js** app lives under **`web/`**—titles, copy, colors (e.g. Tailwind), and routes are yours to change.
- **`NEXT_PUBLIC_API_URL`** must point at your running API when you deploy.

### 5. Production

Read **[Deployment](deployment.md)** for HTTPS, **`wss://`** LiveKit, CORS, and running multiple **agent** replicas. For scale and ops, see **[Scaling](scaling.md)** and **[Monitoring](monitoring.md)**.

---

## If you get stuck

1. Confirm **Docker** runs and **`docker compose up --build`** completes.
2. Confirm **`OPENAI_API_KEY`** is set if you expect **speech** (without it, the agent process may stay up but not speak).
3. Check **browser → API** (CORS) and **browser → LiveKit** (correct WebSocket URL).
4. Run **`make verify`** locally (Go tests + web lint/build) before large changes.

The repository **`README.md`** has the full API table and architecture diagrams. **`teach.md`** (in the repo root) walks the codebase for contributors.

---

## Summary

| Goal | Where to look |
|------|----------------|
| Run locally | [Getting started](getting-started.md) |
| Change tutor behavior | `agent/tutor_agent.py` |
| Change lesson list | `internal/api/scenarios.go` |
| Ship to production | [Deployment](deployment.md) |
| Operate at scale | [Scaling](scaling.md), [Monitoring](monitoring.md) |
| Article (overview) | [Medium](https://medium.com/p/4d894c404772?postPublishedType=initial) |
| LiveKit + voice-agent learning | [github.com/livekit](https://github.com/livekit), [Building AI Voice Agents for Production](https://learn.deeplearning.ai/courses/building-ai-voice-agents-for-production/information) |

You can treat FluentFlow as a **template**: swap scenarios, adjust prompts, deploy on your cloud, and own the product narrative—while keeping a **clear** split between **API**, **LiveKit**, and **agents** that scales with you.
