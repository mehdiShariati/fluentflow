package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mehdi/fluentflow/internal/analytics"
	"github.com/mehdi/fluentflow/internal/auth"
	"github.com/mehdi/fluentflow/internal/store"
)

func (s *Server) guestAuth(w http.ResponseWriter, r *http.Request) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	email := "guest_" + hex.EncodeToString(raw) + "@guest.fluentflow.internal"
	pass := hex.EncodeToString(raw) + hex.EncodeToString(raw)
	hash, err := auth.HashPassword(pass)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	u, err := s.store.CreateGuestUser(r.Context(), email, hash)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	token, err := auth.SignToken(u.ID, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user_id":  u.ID.String(),
		"token":    token,
		"is_guest": true,
	})
}

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	u, err := s.store.UserByID(r.Context(), uid)
	if err != nil || u == nil {
		writeJSON(w, http.StatusNotFound, errResp{"user_not_found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":  u.ID.String(),
		"email":    u.Email,
		"is_guest": u.IsGuest,
	})
}

type transcriptBatchReq struct {
	Segments []struct {
		Speaker  string `json:"speaker"`
		Text     string `json:"text"`
		OffsetMs *int   `json:"offset_ms"`
	} `json:"segments"`
}

func (s *Server) postTranscript(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_session_id"})
		return
	}
	sess, err := s.store.SessionByID(r.Context(), id)
	if err != nil || sess == nil || sess.UserID != uid {
		writeJSON(w, http.StatusNotFound, errResp{"session_not_found"})
		return
	}
	var req transcriptBatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	var rows []struct {
		Speaker  string
		Text     string
		OffsetMs *int
	}
	for _, seg := range req.Segments {
		rows = append(rows, struct {
			Speaker  string
			Text     string
			OffsetMs *int
		}{Speaker: seg.Speaker, Text: strings.TrimSpace(seg.Text), OffsetMs: seg.OffsetMs})
	}
	inserted, err := s.store.InsertTranscriptSegments(r.Context(), id, rows)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if inserted > 0 {
		pay, _ := json.Marshal(map[string]int{"count": inserted})
		_ = s.store.InsertSessionEvent(r.Context(), id, analytics.TranscriptGenerated, pay)
		if s.sessionEvents != nil {
			s.sessionEvents.WithLabelValues(analytics.TranscriptGenerated).Add(float64(inserted))
		}
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"ingested": inserted})
}

func (s *Server) getTranscript(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_session_id"})
		return
	}
	sess, err := s.store.SessionByID(r.Context(), id)
	if err != nil || sess == nil || sess.UserID != uid {
		writeJSON(w, http.StatusNotFound, errResp{"session_not_found"})
		return
	}
	list, err := s.store.ListTranscriptSegments(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if list == nil {
		list = []store.TranscriptSegment{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"segments": list})
}

func (s *Server) getSessionEvents(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_session_id"})
		return
	}
	sess, err := s.store.SessionByID(r.Context(), id)
	if err != nil || sess == nil || sess.UserID != uid {
		writeJSON(w, http.StatusNotFound, errResp{"session_not_found"})
		return
	}
	limit := 100
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	list, err := s.store.ListSessionEvents(r.Context(), id, uid, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if list == nil {
		list = []store.SessionEventRow{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"events": list})
}

type recommendationClickReq struct {
	ScenarioID string `json:"scenario_id"`
	Source     string `json:"source,omitempty"`
}

func (s *Server) postRecommendationClick(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_session_id"})
		return
	}
	sess, err := s.store.SessionByID(r.Context(), id)
	if err != nil || sess == nil || sess.UserID != uid {
		writeJSON(w, http.StatusNotFound, errResp{"session_not_found"})
		return
	}
	var body recommendationClickReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	body.ScenarioID = strings.TrimSpace(body.ScenarioID)
	if body.ScenarioID == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"scenario_id_required"})
		return
	}
	if !IsValidScenarioID(body.ScenarioID) {
		writeJSON(w, http.StatusBadRequest, errResp{"unknown_scenario"})
		return
	}
	pay, _ := json.Marshal(map[string]string{
		"scenario_id": body.ScenarioID,
		"source":      body.Source,
	})
	_ = s.store.InsertSessionEvent(r.Context(), id, analytics.RecommendationClicked, pay)
	if s.sessionEvents != nil {
		s.sessionEvents.WithLabelValues(analytics.RecommendationClicked).Inc()
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) postFeedbackViewed(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_session_id"})
		return
	}
	sess, err := s.store.SessionByID(r.Context(), id)
	if err != nil || sess == nil || sess.UserID != uid {
		writeJSON(w, http.StatusNotFound, errResp{"session_not_found"})
		return
	}
	_ = s.store.InsertSessionEvent(r.Context(), id, analytics.FeedbackViewed, json.RawMessage(`{}`))
	if s.sessionEvents != nil {
		s.sessionEvents.WithLabelValues(analytics.FeedbackViewed).Inc()
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) adminListExperiments(w http.ResponseWriter, r *http.Request) {
	list, err := s.store.ListAllExperiments(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if list == nil {
		list = []store.Experiment{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"experiments": list})
}

type flagPatchReq struct {
	Enabled bool `json:"enabled"`
}

func (s *Server) adminPatchFeatureFlag(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if key == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"missing_key"})
		return
	}
	var body flagPatchReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	if err := s.store.SetFeatureFlag(r.Context(), key, body.Enabled); err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"key": key, "enabled": body.Enabled})
}

func (s *Server) adminMetricsSummary(w http.ResponseWriter, r *http.Request) {
	m, err := s.store.SessionEventCounts(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"session_events_by_type": m})
}
