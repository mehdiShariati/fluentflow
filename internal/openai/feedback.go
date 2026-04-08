package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mehdi/fluentflow/internal/store"
)

// chatCompletionsURL is overridable in tests.
var chatCompletionsURL = "https://api.openai.com/v1/chat/completions"

type chatReq struct {
	Model          string          `json:"model"`
	Messages       []chatMsg       `json:"messages"`
	Temperature    float64         `json:"temperature"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type feedbackJSON struct {
	Strengths            []string `json:"strengths"`
	TopMistakes          []string `json:"top_mistakes"`
	Suggestions          []string `json:"suggestions"`
	RecommendedScenario  string   `json:"recommended_scenario"`
	Score                float64  `json:"score"`
	Notes                string   `json:"notes"`
	TranscriptSummary    string   `json:"transcript_summary"`
}

// SessionFeedback calls OpenAI for structured coaching. TranscriptPlainText may be empty (scenario-only mode).
func SessionFeedback(ctx context.Context, apiKey string, sessionID uuid.UUID, scenarioID string, prof *store.Profile, transcriptPlainText string) (*store.FeedbackSummary, error) {
	if apiKey == "" {
		return nil, nil
	}
	var b strings.Builder
	b.WriteString("Scenario: ")
	b.WriteString(scenarioID)
	b.WriteString(". ")
	if prof != nil {
		b.WriteString(fmt.Sprintf("Learner target language: %s, level %s, goal %s. ",
			prof.TargetLanguage, prof.ProficiencyLevel, prof.LearningGoal))
	}
	if strings.TrimSpace(transcriptPlainText) != "" {
		b.WriteString("\nConversation transcript (speaker: line):\n")
		b.WriteString(transcriptPlainText)
		b.WriteString("\n")
	}
	b.WriteString(`Produce JSON only with keys: strengths (array of 2 short strings), top_mistakes (2 items), suggestions (2 items), recommended_scenario (id like coffee_shop), score (0-10 number), notes (one short string), transcript_summary (2-3 sentence recap of the chat). Be encouraging and practical.`)

	body, _ := json.Marshal(chatReq{
		Model:       "gpt-4o-mini",
		Messages:    []chatMsg{{Role: "user", Content: b.String()}},
		Temperature: 0.6,
		ResponseFormat: &responseFormat{
			Type: "json_object",
		},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{Timeout: 45 * time.Second}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("openai: %s: %s", res.Status, string(raw))
	}
	var outer chatResp
	if err := json.Unmarshal(raw, &outer); err != nil {
		return nil, err
	}
	if len(outer.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty choices")
	}
	var fj feedbackJSON
	if err := json.Unmarshal([]byte(outer.Choices[0].Message.Content), &fj); err != nil {
		return nil, err
	}
	score := fj.Score
	rec := fj.RecommendedScenario
	notes := fj.Notes
	ts := fj.TranscriptSummary
	src := "openai"
	return &store.FeedbackSummary{
		SessionID:           sessionID,
		Strengths:           fj.Strengths,
		TopMistakes:         fj.TopMistakes,
		Suggestions:         fj.Suggestions,
		RecommendedScenario: strOrNil(rec),
		Score:               &score,
		RawNotes:            strOrNil(notes),
		TranscriptSummary:   strOrNil(ts),
		GenerationSource:    &src,
	}, nil
}

func strOrNil(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}
