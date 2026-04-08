package migrate

import (
	"context"
	"embed"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed sql/*.sql
var files embed.FS

// EnsureSchema runs embedded SQL if core tables are missing.
func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	var ok bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'
		)
	`).Scan(&ok)
	if err != nil {
		return err
	}
	if ok {
		if err := runSeedIfNeeded(ctx, pool); err != nil {
			return err
		}
		if err := ensureTranscriptSchema(ctx, pool); err != nil {
			return err
		}
		return ensureFeedbackGenerationSource(ctx, pool)
	}
	b, err := files.ReadFile("sql/000001_init.up.sql")
	if err != nil {
		return err
	}
	if err := execStatements(ctx, pool, string(b)); err != nil {
		return err
	}
	if err := runSeedIfNeeded(ctx, pool); err != nil {
		return err
	}
	if err := ensureTranscriptSchema(ctx, pool); err != nil {
		return err
	}
	return ensureFeedbackGenerationSource(ctx, pool)
}

// ensureTranscriptSchema applies PRD §14.4/14.5 transcript storage for existing DBs.
func ensureTranscriptSchema(ctx context.Context, pool *pgxpool.Pool) error {
	var ok bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'transcript_segments'
		)
	`).Scan(&ok)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	b, err := files.ReadFile("sql/000003_transcript.up.sql")
	if err != nil {
		return err
	}
	return execStatements(ctx, pool, string(b))
}

func ensureFeedbackGenerationSource(ctx context.Context, pool *pgxpool.Pool) error {
	var ok bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns
			WHERE table_schema = 'public' AND table_name = 'feedback_summaries' AND column_name = 'generation_source'
		)
	`).Scan(&ok)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	b, err := files.ReadFile("sql/000004_feedback_generation_source.up.sql")
	if err != nil {
		return err
	}
	return execStatements(ctx, pool, string(b))
}

func runSeedIfNeeded(ctx context.Context, pool *pgxpool.Pool) error {
	var n int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM experiments`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	b, err := files.ReadFile("sql/000002_seed.up.sql")
	if err != nil {
		return err
	}
	return execStatements(ctx, pool, string(b))
}

func execStatements(ctx context.Context, pool *pgxpool.Pool, sql string) error {
	parts := strings.Split(sql, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, err := pool.Exec(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
