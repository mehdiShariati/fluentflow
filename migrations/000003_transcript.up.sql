-- Mirror of internal/migrate/sql/000003_transcript.up.sql (external migration runners)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_guest BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE IF NOT EXISTS transcript_segments (
    id BIGSERIAL PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions (id) ON DELETE CASCADE,
    speaker TEXT NOT NULL,
    text TEXT NOT NULL,
    offset_ms INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS transcript_segments_session ON transcript_segments (session_id, id);

ALTER TABLE feedback_summaries ADD COLUMN IF NOT EXISTS transcript_summary TEXT;
