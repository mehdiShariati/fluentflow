package api

import (
	"encoding/json"
	"net/http"
	"strings"

	openaipkg "github.com/mehdi/fluentflow/internal/openai"
)

type translateReq struct {
	Text           string `json:"text"`
	TargetLanguage string `json:"target_language,omitempty"`
}

func (s *Server) aiTranslate(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	var body translateReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	text := strings.TrimSpace(body.Text)
	if text == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"text_required"})
		return
	}
	target := strings.TrimSpace(body.TargetLanguage)
	if target == "" {
		if p, _ := s.store.ProfileByUserID(r.Context(), uid); p != nil && strings.TrimSpace(p.TargetLanguage) != "" {
			target = p.TargetLanguage
		}
	}
	if target == "" {
		target = "English"
	}
	if s.cfg.OpenAIAPIKey == "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"translated_text": text,
			"source":          "stub",
			"target_language": target,
		})
		return
	}
	out, err := openaipkg.TranslateText(r.Context(), s.cfg.OpenAIAPIKey, text, target)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"translated_text": text,
			"source":          "stub_fallback",
			"target_language": target,
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"translated_text": out,
		"source":          "openai",
		"target_language": target,
	})
}

type analyzeReq struct {
	Text           string `json:"text"`
	TargetLanguage string `json:"target_language,omitempty"`
}

func (s *Server) aiAnalyze(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	var body analyzeReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	text := strings.TrimSpace(body.Text)
	if text == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"text_required"})
		return
	}
	target := strings.TrimSpace(body.TargetLanguage)
	if target == "" {
		if p, _ := s.store.ProfileByUserID(r.Context(), uid); p != nil && strings.TrimSpace(p.TargetLanguage) != "" {
			target = p.TargetLanguage
		}
	}
	if target == "" {
		target = "English"
	}
	if s.cfg.OpenAIAPIKey == "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"analysis": "- Grammar: N/A (OpenAI disabled)\n- Vocabulary: N/A\n- Pronunciation/Fluency: N/A\n- Better version: " + text,
			"source":   "stub",
		})
		return
	}
	out, err := openaipkg.AnalyzeUtterance(r.Context(), s.cfg.OpenAIAPIKey, text, target)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"analysis": "- Grammar: temporary analysis unavailable\n- Vocabulary: temporary analysis unavailable\n- Pronunciation/Fluency: temporary analysis unavailable\n- Better version: " + text,
			"source":   "stub_fallback",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"analysis": out,
		"source":   "openai",
	})
}

