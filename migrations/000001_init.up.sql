-- FluentFlow MVP schema (PRD §21)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    source_language TEXT NOT NULL DEFAULT 'en',
    target_language TEXT NOT NULL,
    proficiency_level TEXT NOT NULL DEFAULT 'A2',
    learning_goal TEXT NOT NULL DEFAULT 'daily_life',
    tutor_style TEXT NOT NULL DEFAULT 'gentle',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE prompt_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version TEXT NOT NULL UNIQUE,
    label TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE experiments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    variants TEXT[] NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE experiment_assignments (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    experiment_key TEXT NOT NULL,
    variant TEXT NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, experiment_key)
);

CREATE TABLE feature_flags (
    key TEXT PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    scenario_id TEXT NOT NULL,
    room_name TEXT NOT NULL,
    prompt_version TEXT NOT NULL,
    experiment_snapshot JSONB NOT NULL DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at TIMESTAMPTZ,
    speaking_seconds INT NOT NULL DEFAULT 0,
    turn_count INT NOT NULL DEFAULT 0
);

CREATE INDEX sessions_user_started ON sessions (user_id, started_at DESC);

CREATE TABLE session_events (
    id BIGSERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX session_events_session ON session_events (session_id, created_at);

CREATE TABLE feedback_summaries (
    session_id UUID PRIMARY KEY REFERENCES sessions (id) ON DELETE CASCADE,
    strengths TEXT[] NOT NULL DEFAULT '{}',
    top_mistakes TEXT[] NOT NULL DEFAULT '{}',
    suggestions TEXT[] NOT NULL DEFAULT '{}',
    recommended_scenario TEXT,
    score NUMERIC(4,2),
    raw_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE learning_metric_snapshots (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    total_sessions INT NOT NULL,
    total_speaking_minutes NUMERIC(12,2) NOT NULL,
    avg_session_score NUMERIC(4,2)
);

CREATE INDEX learning_metric_user ON learning_metric_snapshots (user_id, captured_at DESC);
