package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestSessionFeedback_httptest(t *testing.T) {
	fj := feedbackJSON{
		Strengths:           []string{"Clear sentences"},
		TopMistakes:         []string{"Articles"},
		Suggestions:         []string{"Drill phrases"},
		RecommendedScenario: "business_small_talk",
		Score:               7.5,
		Notes:               "Good session",
		TranscriptSummary:   "Discussed coffee order.",
	}
	content, _ := json.Marshal(fj)
	apiBody, _ := json.Marshal(map[string]interface{}{
		"choices": []map[string]interface{}{{
			"message": map[string]string{"content": string(content)},
		}},
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(apiBody)
	}))
	defer srv.Close()

	old := chatCompletionsURL
	chatCompletionsURL = srv.URL
	defer func() { chatCompletionsURL = old }()

	sid := uuid.New()
	out, err := SessionFeedback(context.Background(), "sk-test", sid, "coffee_shop", nil, "user: hello")
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("nil result")
	}
	if len(out.Strengths) != 1 || out.Score == nil || *out.Score != 7.5 {
		t.Fatalf("unexpected: %+v", out)
	}
	if out.TranscriptSummary == nil || *out.TranscriptSummary != fj.TranscriptSummary {
		t.Fatalf("transcript summary: %+v", out.TranscriptSummary)
	}
}
