// Package analytics documents the PRD §14.6 event taxonomy (ingested via POST /v1/sessions/{id}/events).
package analytics

// Session lifecycle and learning funnel
const (
	SessionCreated       = "session_created"
	SessionJoined        = "session_joined"
	SessionCompleted     = "session_completed"
	ExperimentExposed    = "experiment_exposed"
	FeedbackViewed       = "feedback_viewed"
	RecommendationClicked = "recommendation_clicked"
)

// Real-time turn instrumentation (client or agent may emit)
const (
	TurnStarted         = "turn_started"
	TurnCompleted       = "turn_completed"
	TranscriptGenerated = "transcript_generated"
	CorrectionGenerated = "correction_generated"
	ErrorEmitted        = "error_emitted"
)
