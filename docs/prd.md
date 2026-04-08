
# Product Requirements Document

# Project: **FluentFlow**

## AI-Native Real-Time Language Learning Platform

**Document owner:** Mehdi Shariati / Product Designer / Systems Architect 
**Status:** Portfolio PRD
**Version:** 1.0

---

# 1. Executive Summary

FluentFlow is an AI-native language learning application focused on **real-time speaking practice**, **adaptive learning**, and **production-grade product infrastructure**. The product enables learners to practice a target language through live voice conversations with AI tutors powered by a **LiveKit AI agent**, while the platform continuously measures learning outcomes, latency, reliability, retention, and experiment performance.

This project is designed not only as a usable language-learning product, but also as a proof of capability in:

* end-to-end product thinking
* real-time AI systems
* system design for scale
* observability and monitoring
* experimentation and A/B testing
* AI orchestration and prompt management
* reliability engineering
* shipping AI features with measurable business impact

The portfolio value of this project is that it demonstrates the ability to think beyond UI and features and into the full operating model of a modern AI product.

---

# 2. Product Vision

Build the most effective AI speaking coach for language learners by combining:

* **real-time conversation**
* **personalized feedback**
* **adaptive lesson progression**
* **low-latency AI voice interaction**
* **experiment-driven product development**
* **deep observability across every system layer**

The app should feel like a real language tutor, not a chatbot with a microphone.

---

# 3. Product Thesis

Most language learning apps are strong at memorization and weak at actual language production. They optimize for streaks, taps, and vocabulary drills, but they do not provide enough live, personalized speaking practice.

Users need:

* speaking confidence
* immediate correction
* adaptive difficulty
* realistic conversation flow
* measurable improvement over time

AI, especially voice-native agents, now makes it possible to deliver this at scale.

But the real differentiator is not just “using AI.” It is building an AI product with:

* robust orchestration
* latency-aware architecture
* monitoring by default
* experiment infrastructure from day one
* reliable production thinking

---

# 4. Portfolio Positioning Goal

This project must signal to recruiters, hiring managers, founders, and engineers that the builder can do all of the following:

## Product thinking

* Identify a real user problem
* Define user segments
* Prioritize features
* Connect features to outcomes

## Systems thinking

* Design low-latency real-time flows
* Choose architecture based on tradeoffs
* Plan for scale and failure modes

## AI product thinking

* Use LLMs intentionally, not cosmetically
* Design prompt routing, model selection, and evaluation loops
* Account for hallucination, latency, and cost

## Operational excellence

* Instrument the product from the start
* Define success metrics
* Use A/B tests and feature flags
* Monitor uptime, cost, performance, and user outcomes

---

# 5. Problem Statement

Language learners struggle to become fluent because they lack enough opportunities for safe, accessible, real-time conversation practice.

Current alternatives have major gaps:

## Traditional apps

* Too passive
* Mostly reading, tapping, flashcards
* Weak speaking practice

## Human tutors

* Expensive
* Hard to access frequently
* Scheduling friction

## Generic AI chatbots

* Not built for language pedagogy
* No structured progression
* No learner-specific feedback loop
* Weak observability and learning analytics

The opportunity is to create a real-time AI language coach that behaves like a tutor, adapts like software, and operates like a modern production system.

---

# 6. Primary Users

## Segment A: Busy professionals

**Profile:** 24–40, learning for work, travel, or relocation
**Needs:** flexible speaking practice, low friction, measurable progress
**Pain points:** no time for tutors, inconsistent practice, fear of speaking

## Segment B: Intermediate learners

**Profile:** A2–B2 language learners who know basics but cannot speak smoothly
**Needs:** confidence building, correction, practical conversation
**Pain points:** vocabulary exists but active speaking is weak

## Segment C: Exam-focused learners

**Profile:** IELTS, TOEFL, DELE, DELF learners
**Needs:** speaking simulations, structured feedback, targeted drills
**Pain points:** lack of realistic practice and scoring feedback

## Segment D: Expats / immigrants

**Profile:** learning for daily life
**Needs:** real-world conversation scenarios
**Pain points:** anxiety in live interactions, need for practical speaking

---

# 7. Core User Jobs to Be Done

Users hire FluentFlow to:

* practice speaking without judgment
* receive corrections in real time
* improve conversation skills faster
* build confidence before talking to real people
* get personalized lessons based on mistakes
* measure progress clearly over time

---

# 8. Product Principles

1. **Conversation first**
   The product must optimize for active speaking, not passive content consumption.

2. **Real-time feels magical**
   AI response latency must feel immediate enough to preserve conversational flow.

3. **Feedback must be useful, not overwhelming**
   Corrections should improve confidence, not discourage users.

4. **Everything measurable**
   Every important flow should have metrics, logs, traces, and business outcomes attached.

5. **Experimentation by design**
   Product decisions should be tested, not guessed.

6. **Reliability is part of UX**
   A voice tutor that drops audio or stalls is a broken product.

7. **AI is infrastructure, not decoration**
   AI features must exist to improve learning outcomes, not to sound impressive.

---

# 9. Product Goals

## User goals

* Help users practice speaking daily
* Increase confidence in spoken language
* Improve fluency and grammar over time
* Deliver personalized learning loops

## Business/product goals

* Strong session completion rate
* Strong day-7 and day-30 retention
* High repeat speaking session rate
* Efficient infrastructure and AI cost per session
* Fast experiment cycle time

## Portfolio goals

* Showcase scalable architecture
* Showcase monitoring stack
* Showcase A/B test infrastructure
* Showcase LiveKit AI agent integration
* Showcase AI evaluation and product analytics maturity

---

# 10. Non-Goals for V1

To keep the project focused, the following are explicitly not first-priority:

* multiplayer classrooms
* marketplace for human tutors
* full offline support
* social feed/community features
* advanced gamification economy
* support for dozens of languages on day one
* fully custom speech models from scratch

These can be future roadmap items.

---

# 11. MVP Scope

The MVP must be polished and complete enough to show depth.

## Included in MVP

* user onboarding
* language selection
* proficiency self-assessment
* scenario-based AI voice conversation
* real-time LiveKit AI tutor
* post-session feedback and correction summary
* user progress dashboard
* event tracking
* feature flags
* experimentation support
* observability stack
* admin experiment console or simple internal dashboard
* prompt versioning
* usage analytics
* session replay metadata
* failure and latency instrumentation

## Excluded from MVP

* native mobile app if web is enough for portfolio
* social/community
* human tutor escalation
* subscription billing unless needed
* advanced enterprise admin

---

# 12. Key Product Features

## 12.1 Onboarding and Profile Setup

### User can:

* choose source language and target language
* set learning goal (travel, business, exam, daily life)
* select current level
* choose preferred lesson style
* take a quick placement interaction

### Why it matters:

This creates personalization inputs that affect lesson generation, tutor behavior, and experiment segmentation.

### Metrics:

* onboarding completion rate
* placement completion rate
* first session start conversion

---

## 12.2 Real-Time AI Conversation Practice

### Description

User joins a live speaking session and talks with an AI tutor in the target language. The tutor uses voice, contextual memory, correction strategy, and adaptive prompts.

### Key capabilities

* voice-based turn-taking
* streaming speech-to-text
* LLM-driven response generation
* text-to-speech playback
* interruption handling
* context retention within session
* scenario-based guidance

### Example scenarios

* ordering coffee
* job interview
* airport immigration
* business small talk
* apartment rental conversation
* exam speaking simulation

### Why it matters

This is the hero feature and the strongest proof of AI system design ability.

### Metrics

* session start rate
* session completion rate
* average session duration
* latency per turn
* speech recognition success rate
* user speaking time ratio
* interruption recovery rate

---

## 12.3 LiveKit AI Agent Integration

### Description

The platform uses a LiveKit AI agent to manage real-time voice interactions.

### What to highlight in PRD and README

* low-latency media transport
* real-time audio rooms
* server-side or agent-side orchestration
* barge-in support
* voice activity detection
* conversational turn management
* resilience under packet loss or reconnects

### Why it matters

This is a concrete differentiator. Many portfolio apps say “AI voice.” Few explain a real-time voice agent stack.

### Specific architectural value

* proves familiarity with real-time systems
* proves latency-sensitive design
* proves practical AI shipping beyond text chat

---

## 12.4 Adaptive Feedback Engine

### Description

After each session, the app analyzes the conversation and gives structured feedback.

### Feedback categories

* pronunciation issues
* grammar corrections
* vocabulary suggestions
* fluency observations
* filler-word overuse
* missed expressions
* confidence estimate
* next recommended exercise

### UX output

* “What you did well”
* “Top 3 mistakes”
* “How to say it better”
* “Practice these phrases next”

### Metrics

* feedback open rate
* practice recommendation click-through
* repeat session rate after feedback

---

## 12.5 Personalized Learning Path

### Description

The system updates the user’s learning path based on actual performance, not only stated level.

### Inputs

* repeated grammar mistakes
* hesitation patterns
* low-confidence areas
* scenario performance
* vocabulary gaps
* session completion behavior

### Outputs

* recommend easier or harder scenarios
* focus on specific grammar categories
* suggest targeted speaking exercises

### Why it matters

Shows product depth and AI personalization beyond conversation.

---

## 12.6 Progress Dashboard

### Description

User sees measurable progression over time.

### Example dashboard components

* total speaking minutes
* sessions completed
* fluency trend
* vocabulary growth
* correction frequency trend
* pronunciation confidence trend
* streaks and consistency

### Why it matters

Users need proof of progress; recruiters need proof you understand retention mechanics.

---

## 12.7 Experimentation and Feature Flag Infrastructure

### Description

Every important learning and UX decision can be rolled out behind flags and measured via experiments.

### Candidate experiments

* strict correction vs gentle correction
* tutor personality A vs B
* text feedback vs audio feedback
* scenario recommendation style
* onboarding length
* speaking session countdown UX
* first-session handholding

### Requirements

* user assignment to variants
* event tagging with experiment metadata
* experiment analysis dashboard
* safe rollout and kill switch
* holdout control group support

### Why it matters

This is one of the strongest signals in the project. It shows you think like someone who ships and measures.

---

## 12.8 Admin / Internal Observability Dashboard

### Description

A lightweight internal dashboard for the founder/operator to monitor:

* session health
* latency
* errors
* experiment performance
* prompt version performance
* cost per session
* drop-off points

### Why it matters

Shows operator mindset, not just builder mindset.

---

# 13. User Stories

## Learner stories

* As a learner, I want to choose my target language so the app feels personalized.
* As a learner, I want to practice realistic scenarios so my speaking improves for real situations.
* As a learner, I want instant feedback so I know what to fix.
* As a learner, I want the AI tutor to adapt to my level so the session is not too easy or too hard.
* As a learner, I want to review my mistakes after a session so I can improve quickly.
* As a learner, I want to track progress over time so I stay motivated.

## Operator stories

* As a founder, I want to monitor session latency so I can keep real-time conversations smooth.
* As a founder, I want to compare experiment variants so I can improve retention scientifically.
* As a founder, I want to detect STT/TTS/LLM failures so I can maintain reliability.
* As a founder, I want to see AI cost per session so I can control unit economics.
* As a founder, I want prompt versions tied to outcomes so I can improve tutor quality safely.

---

# 14. Functional Requirements

## 14.1 Authentication

* Email/password or OAuth
* Secure session management
* Support anonymous guest mode optionally
* Persist user progress

## 14.2 User Profile

* target language
* source language
* level
* goal
* preferred tutor style
* learning history

## 14.3 Session Creation

* user selects scenario
* create room/session ID
* connect to LiveKit room
* initialize AI tutor config
* bind experiment flags and prompt version

## 14.4 Real-Time Session Engine

* capture audio from user
* stream audio to real-time media infrastructure
* generate transcript chunks
* send to orchestration layer
* stream response text and TTS back to user
* store structured events and metadata

## 14.5 Post-Session Summary

* transcript summary
* mistakes and corrections
* performance highlights
* recommended next scenario
* session score

## 14.6 Analytics and Event Tracking

Must track:

* session_created
* session_joined
* turn_started
* turn_completed
* transcript_generated
* correction_generated
* session_completed
* feedback_viewed
* recommendation_clicked
* experiment_exposed
* error_emitted

## 14.7 Experiment Framework

* define experiments
* allocate variants
* attach metadata to user and session
* log exposures and outcomes
* allow staged rollout
* allow rollback

## 14.8 Monitoring

* service uptime
* API latency
* AI pipeline latency
* audio room quality
* token usage
* model error rate
* queue size
* DB health
* client error rate

---

# 15. Non-Functional Requirements

## Performance

* user joins session in under 3 seconds
* AI first response target under 1.5 seconds perceived delay
* transcription chunk delivery near real time
* UI remains responsive during streaming

## Scalability

* support growth from 100 daily users to 100k+ monthly active users
* horizontally scalable stateless APIs
* asynchronous background jobs for heavy post-session processing
* queue-based decoupling for non-critical workflows

## Reliability

* 99.9% session service uptime target for production-grade design
* graceful fallback if one AI provider fails
* retry and timeout strategies
* partial degradation instead of total failure

## Security

* secure auth
* encrypt sensitive data in transit
* protect room/session tokens
* redact PII from logs where needed
* follow basic GDPR-minded design if storing transcripts

## Maintainability

* modular services
* clear service boundaries
* shared schema contracts
* versioned prompts and APIs

## Observability

* every critical path traced
* structured logs
* metrics dashboard
* alerting for service degradation

---

# 16. Detailed System Design

This section is crucial for proving your skill.

## 16.1 Recommended Architecture Style

For a portfolio project, the best approach is:

**Modular monolith at first, with service-oriented boundaries**
or
**small service architecture if you want to show stronger infra maturity**

Recommended compromise:

### Frontend

* Next.js web app

### Core backend services

* API / BFF service
* Real-time session orchestrator
* AI pipeline service
* analytics/event ingestion service
* experiment service
* background worker
* monitoring stack

This keeps the project realistic while still showing architectural maturity.

---

## 16.2 High-Level Architecture Components

### Client

* Web app for learners
* Handles auth, onboarding, dashboard, conversation UI
* Streams audio input
* Receives transcripts and AI responses

### API Gateway / Backend-for-Frontend

* auth
* profile
* scenario config
* dashboard data
* experiment fetch
* session creation

### Live Session Service

* provisions session state
* creates/join rooms
* binds session metadata
* coordinates real-time pipeline

### LiveKit Layer

* media transport
* room management
* low-latency audio
* participant state

### AI Orchestration Service

* STT handling
* prompt assembly
* model routing
* TTS output coordination
* context memory
* retry/fallback logic

### User Progress Service

* stores learning history
* calculates progress metrics
* recommends next lessons

### Analytics Pipeline

* ingest events
* enrich with experiment metadata
* aggregate product metrics

### Experimentation Service

* assign variants
* expose feature flags
* log exposures

### Monitoring Stack

* metrics
* logs
* traces
* alerts

### Data Stores

* PostgreSQL for product data
* Redis for transient session/cache state
* object store for session artifacts if needed
* warehouse or analytics DB optionally for event analysis

---

## 16.3 Real-Time Session Flow

### Flow

1. User starts session
2. Client requests session token/config
3. Backend creates session record and assigns experiment flags
4. User joins LiveKit room
5. Audio stream begins
6. Speech is transcribed in streaming mode
7. AI orchestration receives transcript chunks
8. Orchestrator determines response
9. LLM generates tutor reply
10. TTS converts reply to speech
11. Audio streamed back through LiveKit
12. Turn metadata logged
13. After session, transcript and feedback pipeline runs
14. Summary and insights stored
15. Analytics and experiment metrics updated

### Important portfolio note

In your PRD and README, call out where latency matters:

* microphone capture
* VAD/turn detection
* STT chunk delay
* LLM inference time
* TTS synthesis time
* playback startup time

That shows deep real-time thinking.

---

## 16.4 AI Pipeline Design

### Pipeline stages

* audio ingestion
* voice activity detection
* streaming speech-to-text
* language understanding and tutor reasoning
* response generation
* corrective annotation
* speech synthesis
* event emission

### Design considerations

* prompt versioning
* model switching by use case
* cache reusable lesson metadata
* split real-time response generation from slower post-session analysis
* allow “fast response mode” vs “deep feedback mode”

### Example model strategy

* fast model for live turn response
* stronger model for post-session summary and pedagogy analysis

This proves cost and latency awareness.

---

## 16.5 State Management Strategy

### Persistent state

Store in PostgreSQL:

* users
* sessions
* transcripts
* experiment assignments
* lesson outcomes
* progress metrics
* prompt versions

### Ephemeral state

Store in Redis:

* current room/session context
* short-term conversational memory
* active turn status
* idempotency keys
* rate limiting counters

This shows you understand hot-path vs durable storage.

---

## 16.6 Scalability Strategy

You specifically want to prove scale knowledge, so be explicit:

### Phase 1: Early product

* single region
* modular backend
* managed Postgres
* managed Redis
* managed LiveKit deployment
* simple queue

### Phase 2: Growth

* autoscaling stateless services
* separate background workers
* partition analytics ingestion
* isolate heavy post-processing from live path
* introduce CDN and stronger caching

### Phase 3: Large-scale

* multi-region considerations
* region-local media routing
* read replicas
* event streaming backbone
* data warehouse for experiment analysis
* advanced cost controls

### What to say explicitly

The real-time session path must remain thin. Heavy analysis, scoring, and enrichment should be decoupled from the live user interaction path.

That line alone signals strong architectural judgment.

---

# 17. Monitoring and Observability Stack

This is one of your signature strengths. Treat it as a first-class product requirement.

## 17.1 Why observability matters here

A language app with voice AI can fail in many ways:

* user cannot connect
* STT misses words
* LLM is slow
* TTS is unnatural
* room quality degrades
* summaries fail
* experiment tagging breaks
* costs spike invisibly

Without observability, the product is impossible to improve reliably.

---

## 17.2 Monitoring Stack Proposal

### Metrics

* Prometheus

### Dashboards

* Grafana

### Tracing

* OpenTelemetry + Tempo / Jaeger

### Logs

* structured JSON logs to Loki / ELK

### Error tracking

* Sentry for frontend and backend

### Product analytics

* PostHog / Amplitude-style setup

### Feature flags / experiments

* PostHog or custom lightweight service

You can present this as a stack recommendation, not a rigid rule.

---

## 17.3 Critical Metrics to Instrument

### Product metrics

* DAU / WAU / MAU
* day-1 / day-7 / day-30 retention
* first-session completion
* average speaking minutes
* repeat sessions per user
* feature adoption

### Learning metrics

* correction frequency trend
* fluency score trend
* vocabulary reuse
* session confidence score
* scenario completion success

### System metrics

* API p50 / p95 / p99 latency
* STT latency
* LLM response latency
* TTS latency
* session join failure rate
* room reconnect rate
* queue lag
* DB query latency
* cache hit rate

### AI metrics

* tokens per session
* cost per session
* prompt version success rate
* model fallback rate
* hallucination/invalid response markers
* moderation event count

### Experiment metrics

* exposure counts
* variant conversion
* retention by variant
* session length by variant
* feedback usefulness by variant

---

## 17.4 Logging Strategy

Every critical service should emit structured logs with:

* request ID
* session ID
* user ID or anonymized equivalent
* experiment assignment
* prompt version
* model used
* latency
* error type
* fallback path

### Important

Do not log raw sensitive user data carelessly. Mention redaction and sampling rules.

That makes the document feel real and mature.

---

## 17.5 Distributed Tracing

Trace the entire turn lifecycle:

* client event
* room connect
* audio start
* STT chunk ready
* LLM start
* LLM finish
* TTS start
* TTS finish
* playback begin

This is an excellent recruiter signal because it proves you understand where user experience is won or lost.

---

## 17.6 Alerting

Create alerts for:

* session join failures exceed threshold
* p95 AI response latency exceeds threshold
* STT provider error spike
* TTS failure spike
* experiment assignment service unavailable
* token cost anomaly
* DB saturation
* elevated client crash rate

---

# 18. Experimentation and A/B Testing Infrastructure

This section should feel serious, not decorative.

## 18.1 Why experimentation is core

The effectiveness of an AI tutor depends on many product decisions that should be tested:

* correction frequency
* persona warmth
* prompt structure
* onboarding friction
* timing of feedback
* recommendation logic

You want to prove you build products scientifically.

---

## 18.2 Experimentation Architecture

### Components

* feature flag store
* assignment engine
* experiment exposure logger
* event pipeline
* analysis dashboard
* kill switch support

### Assignment rules

* deterministic by user ID
* sticky across sessions
* allow segmentation by language, level, and user cohort

### Exposure logging

Log when a user is actually exposed, not only assigned.

This is a subtle but strong detail. It shows experiment rigor.

---

## 18.3 Example Experiments

### Experiment 1: Tutor correction style

* Variant A: correct every major mistake
* Variant B: correct only after user completes sentence

**Success metrics:**

* session completion
* feedback usefulness
* repeat session rate

### Experiment 2: Tutor tone

* Variant A: professional tutor
* Variant B: friendly conversation partner

**Success metrics:**

* average session duration
* day-7 retention

### Experiment 3: Post-session summary format

* Variant A: text bullet summary
* Variant B: scorecard + examples
* Variant C: short audio recap

**Success metrics:**

* feedback open rate
* next-session conversion

### Experiment 4: Onboarding length

* Variant A: 5-step setup
* Variant B: lightweight setup then infer later

**Success metrics:**

* onboarding completion
* first-session start rate

---

## 18.4 Feature Flags

Use flags for:

* rolling out new tutor persona
* enabling new STT provider
* switching TTS engine
* changing prompt version
* activating new scenario engine
* controlling risky features safely

---

# 19. AI Design and Prompt Strategy

You said you want the project to show AI shipping maturity. This section matters a lot.

## 19.1 AI roles in the system

AI is used for:

* live tutoring
* correction generation
* feedback summarization
* progression recommendation
* personalization analysis
* scenario adaptation

## 19.2 Prompt design principles

* keep system prompts stable and versioned
* separate pedagogical rules from style
* include learner context but avoid excessive prompt bloat
* make correction policy explicit
* define allowed behavior for uncertain cases
* structure outputs for downstream parsing when needed

## 19.3 Prompt versioning

Each session should persist:

* prompt version
* model version
* experiment variant
* language pair
* scenario type

This makes performance and regression analysis possible.

## 19.4 Model routing strategy

Use different models for:

* real-time fast conversational responses
* slower high-quality feedback summaries
* analytics enrichment

This shows understanding of tradeoffs among latency, quality, and cost.

## 19.5 AI safety and quality controls

* block unsafe content
* constrain tutor persona
* prevent harmful or abusive outputs
* detect malformed structured outputs
* retry or fallback on invalid response

---

# 20. Reliability and Failure Mode Design

This section is extremely valuable because most portfolio projects ignore it.

## 20.1 Failure scenarios to plan for

### STT failure

* partial transcript missing
* high latency
* wrong language detection

**Mitigation:**

* retry provider
* fallback to text input mode
* surface graceful UI message

### LLM timeout

* tutor pauses too long
* user loses flow

**Mitigation:**

* timeout threshold
* lightweight fallback response
* queue post-turn analysis separately

### TTS failure

* no audio playback
* robotic voice quality

**Mitigation:**

* fallback voice provider
* show text response if audio unavailable

### LiveKit connection loss

* user disconnects mid-session

**Mitigation:**

* reconnect logic
* preserve session state briefly
* restore room if possible

### Database slowness

* session writes lag
* dashboard stale

**Mitigation:**

* async non-critical writes
* circuit breaker
* degrade non-essential features first

### Analytics outage

* events not delivered

**Mitigation:**

* local buffering / retry
* avoid blocking user flow on analytics success

---

## 20.2 Graceful Degradation Strategy

The real-time conversation must remain prioritized. If non-critical services fail:

* session may continue
* analytics may be delayed
* summary generation may complete later
* dashboard may temporarily show stale data

This demonstrates mature architecture judgment.

---

# 21. Data Model Overview

You do not need full schema in the PRD, but enough to show you’ve thought through it.

## Core entities

* User
* UserProfile
* Session
* RoomConnection
* TranscriptSegment
* TutorResponse
* FeedbackSummary
* LearningMetricSnapshot
* Experiment
* ExperimentAssignment
* FeatureFlag
* PromptVersion
* EventLog
* ErrorEvent

## Key relationships

* user has many sessions
* session has many transcript segments
* session has one feedback summary
* user has many experiment assignments
* session references prompt version and model metadata

---

# 22. Analytics Plan

## 22.1 North Star Metric

**Weekly completed speaking sessions per active learner**

Why:

* directly tied to behavior that drives learning
* more meaningful than raw signups
* captures value delivery and habit formation

## 22.2 Secondary metrics

* first-session completion
* day-7 retention
* average speaking minutes per week
* repeat session rate
* feedback open rate
* scenario progression rate

## 22.3 Diagnostic metrics

* onboarding drop-off by step
* session failure reasons
* latency breakdown
* model errors by provider/version
* cost by user cohort
* experiment exposure integrity

---

# 23. UX Requirements

## Core UX goals

* minimal friction to start a session
* obvious microphone and connection status
* clear speaking turn state
* visible tutor feedback without overload
* reassuring failure messages
* progress dashboard that feels motivating, not punitive

## Important voice UX details

* show “listening” state
* show “thinking” state
* show transcript in real time
* allow interrupting tutor politely
* allow replay of corrections
* allow switching to text if audio quality fails

That level of detail helps the PRD feel credible.

---

# 24. Accessibility Requirements

* keyboard-navigable web experience
* captions/transcripts for voice sessions
* clear focus states
* readable contrast
* audio controls
* alternative input methods where possible

This makes the project stronger and more thoughtful.

---

# 25. Security and Privacy Requirements

* secure authentication and token handling
* room access tokens should be short-lived
* transcripts stored securely
* avoid logging raw sensitive content unnecessarily
* explicit consent if storing conversation history
* allow user to delete history
* internal dashboards must be access-controlled

---

# 26. Cost and Unit Economics Thinking

This is optional in many PRDs, but for your portfolio it is a killer addition.

## Cost drivers

* STT cost
* LLM tokens
* TTS cost
* Live media infra
* storage and analytics

## Cost controls

* use smaller fast model during live turns where acceptable
* run heavy analysis only post-session
* summarize transcript before passing to expensive model
* cache repeated scenario assets
* rate-limit abuse
* track cost per completed session

This tells recruiters you think like a founder, not just a builder.

---

# 27. Success Criteria for MVP

The MVP is successful if it can demonstrate:

## Product success

* users can complete end-to-end speaking sessions
* feedback is generated and useful
* dashboard shows progress
* at least one meaningful experiment is runnable

## System success

* real-time interaction feels smooth enough
* latency is measured and visible
* errors are traceable
* architecture can scale with clear next steps

## Portfolio success

A recruiter or hiring manager should be able to look at the project and say:

* this person understands product
* this person understands systems
* this person understands operating AI in production
* this person measures what they build

---

# 28. Roadmap

## V1

* onboarding
* core voice tutor
* post-session feedback
* dashboard
* feature flags
* observability
* 1–2 experiments

## V2

* deeper personalization
* stronger lesson planning
* richer analytics
* prompt experimentation dashboard
* roleplay scenarios library expansion

## V3

* exam prep modes
* social/cohort features
* native mobile client
* multi-agent tutor orchestration
* multilingual support expansion

---

# 29. Risks

## Product risks

* users may enjoy the novelty but not retain
* correction style may reduce confidence
* feedback may feel repetitive

## Technical risks

* live latency may degrade experience
* AI costs can rise quickly
* voice quality inconsistency can hurt trust
* observability setup may become too complex for early stage

## Mitigations

* strong instrumentation from day one
* focused MVP
* experiments on user experience
* cheap/fast vs rich/slow model split
* degrade gracefully

---

# 30. Why This Project Is a Strong Portfolio Piece

This project proves:

* product strategy
* deep system design
* real-time architecture
* AI infrastructure knowledge
* observability maturity
* experiment-driven product development
* practical founder thinking
* ability to build software with business and reliability awareness

It is much stronger than a generic AI chatbot or simple SaaS clone.

---

# 31. How to Present This Project So It Looks Elite

Your project should not be presented as just:

> “An AI language learning app.”

It should be presented as:

> “A production-oriented real-time AI language learning platform with LiveKit voice agents, experiment infrastructure, observability, and adaptive learning systems.”

That framing immediately upgrades how people read your work.

---

# 32. Recruiter-Killer README Structure

This matters a lot. A weak README can make a strong project look average.

Your README should not be a dump of setup commands. It should sell the product, the architecture, and your thinking in under 2 minutes.

Below is the structure I recommend.

---

## README Title

# FluentFlow — Real-Time AI Language Tutor

## One-line value proposition

AI-native language learning platform for real-time speaking practice, built with LiveKit voice agents, observability-first architecture, and experiment infrastructure.

---

## Section 1: Why this project exists

Explain the problem in 4–6 lines.

Example:

> Most language learning apps optimize for passive repetition instead of real spoken fluency. FluentFlow explores how real-time AI voice tutors, adaptive feedback, and measurable experimentation can create a better speaking-first learning experience.

This tells recruiters you are solving a real problem.

---

## Section 2: What this project demonstrates

This section is extremely important.

Example:

* Real-time AI voice interaction with LiveKit
* AI orchestration for STT → LLM → TTS loop
* Product analytics and event taxonomy
* Feature flags and A/B testing support
* Observability stack with metrics, logs, and traces
* Scalable system design decisions
* Adaptive learning and post-session feedback

This section lets recruiters instantly understand why the project is impressive.

---

## Section 3: Product walkthrough

Add screenshots or GIFs:

* onboarding
* live session screen
* post-session feedback
* progress dashboard
* internal experiment/monitoring panel

Use visuals. They make the project memorable.

---

## Section 4: Architecture overview

Include a diagram.

At minimum show:

* frontend
* backend/API
* LiveKit
* STT
* LLM
* TTS
* Postgres
* Redis
* analytics pipeline
* observability stack
* experiment service

The architecture section is one of the biggest recruiter hooks.

---

## Section 5: Real-time AI flow

Write a short but sharp flow:

1. User joins a LiveKit room
2. Audio is streamed and transcribed
3. The AI orchestration layer builds tutor context
4. LLM generates response and correction strategy
5. TTS returns speech in real time
6. Events, traces, and experiment metadata are logged
7. Post-session analysis updates user progress

This section should feel clear and technical.

---

## Section 6: Key technical decisions

This is gold.

Write:

* Why LiveKit over building raw WebRTC primitives
* Why separate real-time path from async analysis
* Why Redis for ephemeral session state
* Why feature flags were included from the start
* Why observability is part of core architecture

Recruiters love deliberate tradeoffs.

---

## Section 7: Observability and experiments

Most READMEs skip this. Yours should not.

Write:

* what metrics are collected
* how traces follow turn lifecycle
* what experiments can be run
* example dashboard views
* why this matters for product quality

This section will make your project look senior.

---

## Section 8: Scalability notes

Explain:

* how the architecture scales horizontally
* how asynchronous jobs reduce load on live paths
* what would change from 1k users to 1M users
* what bottlenecks are anticipated

This helps prove system design maturity.

---

## Section 9: Local development

Only here do you add setup.

Keep it clean:

* prerequisites
* environment variables
* install
* run frontend/backend/workers
* monitoring stack commands
* seed data

Do not put this before the value proposition.

---

## Section 10: Future improvements

This shows thinking range.

Examples:

* multi-region session routing
* stronger pronunciation scoring
* deeper model evaluation
* personalized lesson planning
* mobile clients
* exam-specific tutoring modes

---

## Section 11: What I’d build next with more time

A very strong founder-style section.

Examples:

* experiment analysis UI
* cost-aware model router
* regional failover
* better learner knowledge graph
* richer automated quality evaluation

This makes the project feel alive and strategic.

---

# 33. README Writing Rules That Make It Strong

## Do this

* lead with problem and value
* explain what makes the system special
* use diagrams and screenshots
* include tradeoffs
* include architecture decisions
* include metrics/observability
* write like a builder-operator, not a student

## Avoid this

* starting with “This is a full-stack app built with X, Y, Z”
* giant walls of setup text at the top
* generic feature lists with no depth
* no product motivation
* no architectural explanation
* no screenshots
* no outcomes or tradeoffs

---

# 34. Best Portfolio Packaging

To make this project truly strong, package it as four assets:

## 1. PRD

What you asked for here.

## 2. README

The recruiter-facing version.

## 3. Architecture diagram

Clean, polished, easy to read.

## 4. Case study

A concise write-up:

* problem
* approach
* system design
* metrics
* tradeoffs
* lessons learned

This combination is far more powerful than code alone.

---

# 35. Final Recommendation

This is absolutely the right kind of project for you.

But the winning version is not:

> “AI language app with many features.”

The winning version is:

> “A real-time AI language tutoring platform that proves I can design and ship production-grade systems with observability, experimentation, scalability, and AI orchestration.”

That is the story.

And this PRD should support that exact story.

If you want, next I can turn this into a **polished portfolio-ready document with cleaner formatting and executive style**, or I can write the **recruiter-killer README in full**.
