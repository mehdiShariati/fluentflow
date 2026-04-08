package api

import (
	"testing"

	"github.com/google/uuid"

	"github.com/mehdi/fluentflow/internal/store"
)

func TestEnrichFeedbackRecommendedScenarioTitle(t *testing.T) {
	rec := "coffee_shop"
	f := &store.FeedbackSummary{
		SessionID:           uuid.New(),
		RecommendedScenario: &rec,
	}
	out := enrichFeedback(f)
	if out.RecommendedScenarioTitle != "Ordering coffee" {
		t.Fatalf("title: got %q want Ordering coffee", out.RecommendedScenarioTitle)
	}
}

func TestEnrichFeedbackNil(t *testing.T) {
	out := enrichFeedback(nil)
	if out.SessionID != uuid.Nil {
		t.Fatalf("enrichFeedback(nil): SessionID = %v, want zero UUID", out.SessionID)
	}
	if out.RecommendedScenarioTitle != "" {
		t.Fatalf("RecommendedScenarioTitle should be empty")
	}
}
