package api

import (
	"github.com/google/uuid"

	"github.com/mehdi/fluentflow/internal/store"
)

// buildStubFeedback produces deterministic placeholder feedback when no LLM worker is configured (PRD §12.4).
func buildStubFeedback(sessionID uuid.UUID, scenarioID string) *store.FeedbackSummary {
	rec := "coffee_shop"
	if scenarioID == "coffee_shop" {
		rec = "business_small_talk"
	}
	score := 7.5
	src := "stub"
	return &store.FeedbackSummary{
		SessionID: sessionID,
		Strengths: []string{
			"You kept the conversation moving and used complete sentences.",
			"Good job staying on topic for the scenario.",
		},
		TopMistakes: []string{
			"Article usage before countable nouns.",
			"Past tense consistency when telling a short story.",
		},
		Suggestions: []string{
			"Practice short monologues (30s) and record yourself.",
			"Drill 5 useful phrases for this scenario before the next session.",
		},
		RecommendedScenario: &rec,
		Score:               &score,
		RawNotes:            strPtr("Stub feedback: connect an LLM worker to replace this with real analysis."),
		GenerationSource:    &src,
	}
}

func strPtr(s string) *string { return &s }

// enrichFeedback copies f and sets recommended_scenario_title from the static catalog (response-only; not stored).
func enrichFeedback(f *store.FeedbackSummary) store.FeedbackSummary {
	if f == nil {
		return store.FeedbackSummary{}
	}
	out := *f
	out.RecommendedScenarioTitle = ""
	if f.RecommendedScenario != nil && *f.RecommendedScenario != "" {
		out.RecommendedScenarioTitle = ScenarioTitle(*f.RecommendedScenario)
	}
	return out
}
