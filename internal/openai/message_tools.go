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
)

func oneShotText(ctx context.Context, apiKey, prompt string) (string, error) {
	body, _ := json.Marshal(chatReq{
		Model:       "gpt-4o-mini",
		Messages:    []chatMsg{{Role: "user", Content: prompt}},
		Temperature: 0.2,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{Timeout: 35 * time.Second}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("openai: %s: %s", res.Status, string(raw))
	}
	var outer chatResp
	if err := json.Unmarshal(raw, &outer); err != nil {
		return "", err
	}
	if len(outer.Choices) == 0 {
		return "", fmt.Errorf("openai: empty choices")
	}
	return strings.TrimSpace(outer.Choices[0].Message.Content), nil
}

func TranslateText(ctx context.Context, apiKey, text, targetLanguage string) (string, error) {
	prompt := fmt.Sprintf(
		"Translate the following text into %s. Keep meaning and tone. Return only translated text.\n\nText:\n%s",
		targetLanguage,
		text,
	)
	return oneShotText(ctx, apiKey, prompt)
}

func AnalyzeUtterance(ctx context.Context, apiKey, text, targetLanguage string) (string, error) {
	prompt := fmt.Sprintf(
		"Analyze this learner utterance for %s speaking practice. Return concise markdown with 4 bullets: Grammar, Vocabulary, Pronunciation/Fluency, Better version.\n\nUtterance:\n%s",
		targetLanguage,
		text,
	)
	return oneShotText(ctx, apiKey, prompt)
}

