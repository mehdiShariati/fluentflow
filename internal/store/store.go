package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsGuest      bool
	CreatedAt    time.Time
}

type Profile struct {
	UserID            uuid.UUID
	SourceLanguage    string
	TargetLanguage    string
	ProficiencyLevel  string
	LearningGoal      string
	TutorStyle        string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Session struct {
	ID                 uuid.UUID       `json:"id"`
	UserID             uuid.UUID       `json:"user_id"`
	ScenarioID         string          `json:"scenario_id"`
	RoomName           string          `json:"room_name"`
	PromptVersion      string          `json:"prompt_version"`
	ExperimentSnapshot json.RawMessage `json:"experiment_snapshot"`
	Status             string          `json:"status"`
	StartedAt          time.Time       `json:"started_at"`
	EndedAt            *time.Time      `json:"ended_at,omitempty"`
	SpeakingSeconds    int             `json:"speaking_seconds"`
	TurnCount          int             `json:"turn_count"`
}

type FeedbackSummary struct {
	SessionID            uuid.UUID `json:"session_id"`
	Strengths            []string  `json:"strengths"`
	TopMistakes          []string  `json:"top_mistakes"`
	Suggestions          []string  `json:"suggestions"`
	RecommendedScenario      *string `json:"recommended_scenario"`
	RecommendedScenarioTitle string  `json:"recommended_scenario_title,omitempty"`
	Score                *float64  `json:"score"`
	RawNotes             *string   `json:"raw_notes"`
	TranscriptSummary    *string   `json:"transcript_summary,omitempty"`
	GenerationSource     *string   `json:"generation_source,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// TranscriptSegment is a stored STT / caption chunk (PRD §14.4).
type TranscriptSegment struct {
	ID        int64     `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	Speaker   string    `json:"speaker"`
	Text      string    `json:"text"`
	OffsetMs  *int      `json:"offset_ms,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Experiment struct {
	Key      string
	Name     string
	Variants []string
	Active   bool
}

type DashboardSummary struct {
	TotalSessions     int        `json:"total_sessions"`
	TotalSpeakingMins float64    `json:"total_speaking_mins"`
	CompletedSessions int        `json:"completed_sessions"`
	LastSessionAt     *time.Time `json:"last_session_at"`
	AvgScore          *float64   `json:"avg_score"`
}

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *Store) CreateGuestUser(ctx context.Context, email, passwordHash string) (*User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, is_guest) VALUES ($1, $2, true)
		RETURNING id, email, password_hash, is_guest, created_at
	`, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsGuest, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, is_guest) VALUES ($1, $2, false)
		RETURNING id, email, password_hash, is_guest, created_at
	`, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsGuest, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, COALESCE(is_guest,false), created_at FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsGuest, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, COALESCE(is_guest,false), created_at FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsGuest, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UpsertProfile(ctx context.Context, p *Profile) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO user_profiles (user_id, source_language, target_language, proficiency_level, learning_goal, tutor_style)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			source_language = EXCLUDED.source_language,
			target_language = EXCLUDED.target_language,
			proficiency_level = EXCLUDED.proficiency_level,
			learning_goal = EXCLUDED.learning_goal,
			tutor_style = EXCLUDED.tutor_style,
			updated_at = now()
	`, p.UserID, p.SourceLanguage, p.TargetLanguage, p.ProficiencyLevel, p.LearningGoal, p.TutorStyle)
	return err
}

func (s *Store) ProfileByUserID(ctx context.Context, userID uuid.UUID) (*Profile, error) {
	var p Profile
	err := s.pool.QueryRow(ctx, `
		SELECT user_id, source_language, target_language, proficiency_level, learning_goal, tutor_style, created_at, updated_at
		FROM user_profiles WHERE user_id = $1
	`, userID).Scan(
		&p.UserID, &p.SourceLanguage, &p.TargetLanguage, &p.ProficiencyLevel,
		&p.LearningGoal, &p.TutorStyle, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) ActivePromptVersion(ctx context.Context) (string, error) {
	var v string
	err := s.pool.QueryRow(ctx, `SELECT version FROM prompt_versions WHERE active = true ORDER BY created_at DESC LIMIT 1`).Scan(&v)
	if err == pgx.ErrNoRows {
		return "v1-default", nil
	}
	if err != nil {
		return "", err
	}
	return v, nil
}

func (s *Store) ListActiveExperiments(ctx context.Context) ([]Experiment, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, name, variants, active FROM experiments WHERE active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Experiment
	for rows.Next() {
		var e Experiment
		if err := rows.Scan(&e.Key, &e.Name, &e.Variants, &e.Active); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) GetExperimentAssignment(ctx context.Context, userID uuid.UUID, experimentKey string) (string, error) {
	var v string
	err := s.pool.QueryRow(ctx, `
		SELECT variant FROM experiment_assignments WHERE user_id = $1 AND experiment_key = $2
	`, userID, experimentKey).Scan(&v)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return v, nil
}

func (s *Store) SetExperimentAssignment(ctx context.Context, userID uuid.UUID, experimentKey, variant string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO experiment_assignments (user_id, experiment_key, variant) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, experiment_key) DO UPDATE SET variant = EXCLUDED.variant, assigned_at = now()
	`, userID, experimentKey, variant)
	return err
}

func (s *Store) ListFeatureFlags(ctx context.Context) (map[string]bool, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, enabled FROM feature_flags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]bool)
	for rows.Next() {
		var k string
		var en bool
		if err := rows.Scan(&k, &en); err != nil {
			return nil, err
		}
		m[k] = en
	}
	return m, rows.Err()
}

func (s *Store) CreateSession(ctx context.Context, userID uuid.UUID, scenarioID, roomName, promptVersion string, expSnap json.RawMessage) (*Session, error) {
	var sess Session
	err := s.pool.QueryRow(ctx, `
		INSERT INTO sessions (user_id, scenario_id, room_name, prompt_version, experiment_snapshot)
		VALUES ($1, $2, $3, $4, COALESCE($5::jsonb, '{}'::jsonb))
		RETURNING id, user_id, scenario_id, room_name, prompt_version, experiment_snapshot, status, started_at, ended_at, speaking_seconds, turn_count
	`, userID, scenarioID, roomName, promptVersion, expSnap).Scan(
		&sess.ID, &sess.UserID, &sess.ScenarioID, &sess.RoomName, &sess.PromptVersion,
		&sess.ExperimentSnapshot, &sess.Status, &sess.StartedAt, &sess.EndedAt,
		&sess.SpeakingSeconds, &sess.TurnCount,
	)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *Store) SessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	var sess Session
	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, scenario_id, room_name, prompt_version, experiment_snapshot, status, started_at, ended_at, speaking_seconds, turn_count
		FROM sessions WHERE id = $1
	`, id).Scan(
		&sess.ID, &sess.UserID, &sess.ScenarioID, &sess.RoomName, &sess.PromptVersion,
		&sess.ExperimentSnapshot, &sess.Status, &sess.StartedAt, &sess.EndedAt,
		&sess.SpeakingSeconds, &sess.TurnCount,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

func (s *Store) ListSessionsForUser(ctx context.Context, userID uuid.UUID, limit int) ([]Session, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, scenario_id, room_name, prompt_version, experiment_snapshot, status, started_at, ended_at, speaking_seconds, turn_count
		FROM sessions WHERE user_id = $1 ORDER BY started_at DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		var sess Session
		if err := rows.Scan(
			&sess.ID, &sess.UserID, &sess.ScenarioID, &sess.RoomName, &sess.PromptVersion,
			&sess.ExperimentSnapshot, &sess.Status, &sess.StartedAt, &sess.EndedAt,
			&sess.SpeakingSeconds, &sess.TurnCount,
		); err != nil {
			return nil, err
		}
		out = append(out, sess)
	}
	return out, rows.Err()
}

func (s *Store) CompleteSession(ctx context.Context, id uuid.UUID, speakingSec, turnCount int) error {
	ct, err := s.pool.Exec(ctx, `
		UPDATE sessions SET status = 'completed', ended_at = now(), speaking_seconds = $2, turn_count = $3
		WHERE id = $1 AND status = 'active'
	`, id, speakingSec, turnCount)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (s *Store) InsertSessionEvent(ctx context.Context, sessionID uuid.UUID, eventType string, payload json.RawMessage) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO session_events (session_id, event_type, payload) VALUES ($1, $2, COALESCE($3::jsonb, '{}'::jsonb))
	`, sessionID, eventType, payload)
	return err
}

// SessionEventRow is a persisted analytics row (PRD §14.6).
type SessionEventRow struct {
	ID        int64           `json:"id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func (s *Store) ListSessionEvents(ctx context.Context, sessionID, userID uuid.UUID, limit int) ([]SessionEventRow, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		SELECT e.id, e.event_type, e.payload, e.created_at
		FROM session_events e
		INNER JOIN sessions s ON s.id = e.session_id
		WHERE e.session_id = $1 AND s.user_id = $2
		ORDER BY e.id ASC
		LIMIT $3
	`, sessionID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SessionEventRow
	for rows.Next() {
		var r SessionEventRow
		if err := rows.Scan(&r.ID, &r.EventType, &r.Payload, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) UpsertFeedback(ctx context.Context, f *FeedbackSummary) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO feedback_summaries (session_id, strengths, top_mistakes, suggestions, recommended_scenario, score, raw_notes, transcript_summary, generation_source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (session_id) DO UPDATE SET
			strengths = EXCLUDED.strengths,
			top_mistakes = EXCLUDED.top_mistakes,
			suggestions = EXCLUDED.suggestions,
			recommended_scenario = EXCLUDED.recommended_scenario,
			score = EXCLUDED.score,
			raw_notes = EXCLUDED.raw_notes,
			transcript_summary = EXCLUDED.transcript_summary,
			generation_source = EXCLUDED.generation_source,
			created_at = now()
	`, f.SessionID, f.Strengths, f.TopMistakes, f.Suggestions, f.RecommendedScenario, f.Score, f.RawNotes, f.TranscriptSummary, f.GenerationSource)
	return err
}

func (s *Store) FeedbackBySession(ctx context.Context, sessionID uuid.UUID) (*FeedbackSummary, error) {
	var f FeedbackSummary
	var ts sql.NullString
	var gen sql.NullString
	err := s.pool.QueryRow(ctx, `
		SELECT session_id, strengths, top_mistakes, suggestions, recommended_scenario, score, raw_notes,
		       transcript_summary, generation_source, created_at
		FROM feedback_summaries WHERE session_id = $1
	`, sessionID).Scan(
		&f.SessionID, &f.Strengths, &f.TopMistakes, &f.Suggestions,
		&f.RecommendedScenario, &f.Score, &f.RawNotes, &ts, &gen, &f.CreatedAt,
	)
	if err == nil && gen.Valid {
		v := gen.String
		f.GenerationSource = &v
	}
	if err == nil && ts.Valid {
		f.TranscriptSummary = &ts.String
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *Store) InsertTranscriptSegments(ctx context.Context, sessionID uuid.UUID, segments []struct {
	Speaker  string
	Text     string
	OffsetMs *int
}) (inserted int, err error) {
	if len(segments) == 0 {
		return 0, nil
	}
	for _, seg := range segments {
		if seg.Text == "" {
			continue
		}
		sp := seg.Speaker
		if sp == "" {
			sp = "user"
		}
		_, execErr := s.pool.Exec(ctx, `
			INSERT INTO transcript_segments (session_id, speaker, text, offset_ms) VALUES ($1, $2, $3, $4)
		`, sessionID, sp, seg.Text, seg.OffsetMs)
		if execErr != nil {
			return inserted, execErr
		}
		inserted++
	}
	return inserted, nil
}

func (s *Store) ListTranscriptSegments(ctx context.Context, sessionID uuid.UUID) ([]TranscriptSegment, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, session_id, speaker, text, offset_ms, created_at
		FROM transcript_segments WHERE session_id = $1 ORDER BY id ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TranscriptSegment
	for rows.Next() {
		var t TranscriptSegment
		var off sql.NullInt32
		if err := rows.Scan(&t.ID, &t.SessionID, &t.Speaker, &t.Text, &off, &t.CreatedAt); err != nil {
			return nil, err
		}
		if off.Valid {
			v := int(off.Int32)
			t.OffsetMs = &v
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) TranscriptPlainText(ctx context.Context, sessionID uuid.UUID) (string, error) {
	segs, err := s.ListTranscriptSegments(ctx, sessionID)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for _, seg := range segs {
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(seg.Speaker)
		b.WriteString(": ")
		b.WriteString(seg.Text)
	}
	return b.String(), nil
}

func (s *Store) ListAllExperiments(ctx context.Context) ([]Experiment, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, name, variants, active FROM experiments ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Experiment
	for rows.Next() {
		var e Experiment
		if err := rows.Scan(&e.Key, &e.Name, &e.Variants, &e.Active); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) SetFeatureFlag(ctx context.Context, key string, enabled bool) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO feature_flags (key, enabled) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET enabled = EXCLUDED.enabled, updated_at = now()
	`, key, enabled)
	return err
}

func (s *Store) DashboardForUser(ctx context.Context, userID uuid.UUID) (*DashboardSummary, error) {
	var d DashboardSummary
	var avg sql.NullFloat64
	err := s.pool.QueryRow(ctx, `
		SELECT
			COUNT(*)::int,
			COALESCE(SUM(speaking_seconds), 0)::float8 / 60.0,
			COUNT(*) FILTER (WHERE status = 'completed')::int,
			MAX(started_at) FILTER (WHERE status = 'completed'),
			(SELECT AVG(score) FROM feedback_summaries fs
			 INNER JOIN sessions s ON s.id = fs.session_id
			 WHERE s.user_id = $1 AND fs.score IS NOT NULL)
		FROM sessions WHERE user_id = $1
	`, userID).Scan(&d.TotalSessions, &d.TotalSpeakingMins, &d.CompletedSessions, &d.LastSessionAt, &avg)
	if err != nil {
		return nil, err
	}
	if avg.Valid {
		v := avg.Float64
		d.AvgScore = &v
	}
	return &d, nil
}
