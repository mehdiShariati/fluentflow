package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// LearningMetricSnapshot is a point-in-time rollup (PRD §12.6 / learning history).
type LearningMetricSnapshot struct {
	ID                  int64      `json:"id"`
	UserID              uuid.UUID  `json:"user_id"`
	CapturedAt          time.Time  `json:"captured_at"`
	TotalSessions       int        `json:"total_sessions"`
	TotalSpeakingMins   float64    `json:"total_speaking_minutes"`
	AvgSessionScore     *float64   `json:"avg_session_score"`
}

// AppendLearningMetricSnapshot inserts one row with aggregates after a milestone (e.g. session completed).
func (s *Store) AppendLearningMetricSnapshot(ctx context.Context, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO learning_metric_snapshots (user_id, total_sessions, total_speaking_minutes, avg_session_score)
		SELECT $1,
			COUNT(*)::int,
			(COALESCE(SUM(speaking_seconds), 0)::float8 / 60.0),
			(SELECT AVG(fs.score::float8) FROM feedback_summaries fs
			 INNER JOIN sessions s2 ON s2.id = fs.session_id
			 WHERE s2.user_id = $1 AND fs.score IS NOT NULL)
		FROM sessions s WHERE s.user_id = $1
	`, userID)
	return err
}

func (s *Store) ListLearningMetricSnapshots(ctx context.Context, userID uuid.UUID, limit int) ([]LearningMetricSnapshot, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, captured_at, total_sessions, total_speaking_minutes, avg_session_score
		FROM learning_metric_snapshots
		WHERE user_id = $1
		ORDER BY captured_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LearningMetricSnapshot
	for rows.Next() {
		var m LearningMetricSnapshot
		var avg sql.NullFloat64
		if err := rows.Scan(&m.ID, &m.UserID, &m.CapturedAt, &m.TotalSessions, &m.TotalSpeakingMins, &avg); err != nil {
			return nil, err
		}
		if avg.Valid {
			v := avg.Float64
			m.AvgSessionScore = &v
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	ct, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
