package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mehdi/fluentflow/internal/analytics"
	"github.com/mehdi/fluentflow/internal/auth"
	"github.com/mehdi/fluentflow/internal/config"
	"github.com/mehdi/fluentflow/internal/experiment"
	openaipkg "github.com/mehdi/fluentflow/internal/openai"
	"github.com/mehdi/fluentflow/internal/store"
)

type Server struct {
	cfg   *config.Config
	store *store.Store
	reg   *prometheus.Registry
	httpRequests  *prometheus.CounterVec
	httpLatency   *prometheus.HistogramVec
	sessionEvents *prometheus.CounterVec
}

func New(cfg *config.Config, st *store.Store) *Server {
	reg := prometheus.NewRegistry()
	req := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fluentflow_http_requests_total",
		Help: "HTTP requests by route pattern and status.",
	}, []string{"method", "path", "code"})
	lat := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fluentflow_http_request_duration_seconds",
		Help:    "HTTP request latency.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
	sev := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fluentflow_session_events_ingested_total",
		Help: "Session-scoped product events POSTed to the API (PRD §14.6).",
	}, []string{"event_type"})
	reg.MustRegister(req)
	reg.MustRegister(lat)
	reg.MustRegister(sev)
	return &Server{cfg: cfg, store: st, reg: reg, httpRequests: req, httpLatency: lat, sessionEvents: sev}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.cfg.CorsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(s.observe)

	r.Get("/healthz", s.healthz)
	r.Handle("/metrics", promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/register", s.register)
		r.Post("/auth/login", s.login)
		r.Post("/auth/guest", s.guestAuth)

		r.Group(func(r chi.Router) {
			r.Use(s.requireUser)
			r.Get("/me", s.getMe)
			r.Delete("/me/account", s.deleteAccount)
			r.Get("/me/learning-snapshots", s.listLearningSnapshots)
			r.Get("/me/profile", s.getProfile)
			r.Put("/me/profile", s.putProfile)
			r.Get("/scenarios", s.getScenarios)
			r.Get("/experiments", s.getExperiments)
			r.Get("/feature-flags", s.getFeatureFlags)
			r.Post("/ai/translate", s.aiTranslate)
			r.Post("/ai/analyze", s.aiAnalyze)
			r.Get("/sessions", s.listSessions)
			r.Post("/sessions", s.createSession)
			r.Get("/sessions/{id}", s.getSession)
			r.Post("/sessions/{id}/livekit-token", s.issueLiveKitToken)
			r.Post("/sessions/{id}/events", s.postSessionEvents)
			r.Post("/sessions/{id}/transcript", s.postTranscript)
			r.Get("/sessions/{id}/transcript", s.getTranscript)
			r.Post("/sessions/{id}/complete", s.completeSession)
			r.Get("/sessions/{id}/feedback", s.getFeedback)
			r.Post("/sessions/{id}/feedback/generate", s.generateFeedback)
			r.Post("/sessions/{id}/feedback/viewed", s.postFeedbackViewed)
			r.Post("/sessions/{id}/recommendation-click", s.postRecommendationClick)
			r.Get("/sessions/{id}/events", s.getSessionEvents)
			r.Get("/dashboard/summary", s.dashboardSummary)
		})
	})

	r.Route("/internal/v1", func(r chi.Router) {
		r.Use(s.requireAdmin)
		r.Get("/overview", s.adminOverview)
		r.Get("/experiments", s.adminListExperiments)
		r.Patch("/feature-flags/{key}", s.adminPatchFeatureFlag)
		r.Get("/metrics/summary", s.adminMetricsSummary)
	})

	return r
}

func (s *Server) observe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, 0)
		next.ServeHTTP(ww, r)
		route := r.URL.Path
		if rc := chi.RouteContext(r.Context()); rc != nil {
			if p := rc.RoutePattern(); p != "" {
				route = p
			}
		}
		code := strconv.Itoa(ww.Status())
		s.httpRequests.WithLabelValues(r.Method, route, code).Inc()
		s.httpLatency.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := s.store.Pool().Ping(ctx); err != nil {
		http.Error(w, `{"error":"db_unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// --- auth ---

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || len(req.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, errResp{"email_and_password_required"})
		return
	}
	exists, err := s.store.UserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Printf("register: %v", err)
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if exists != nil {
		writeJSON(w, http.StatusConflict, errResp{"email_taken"})
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	u, err := s.store.CreateUser(r.Context(), req.Email, hash)
	if err != nil {
		log.Printf("create user: %v", err)
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	token, err := auth.SignToken(u.ID, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user_id": u.ID.String(),
		"token":   token,
	})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	u, err := s.store.UserByEmail(r.Context(), req.Email)
	if err != nil || u == nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		writeJSON(w, http.StatusUnauthorized, errResp{"invalid_credentials"})
		return
	}
	token, err := auth.SignToken(u.ID, s.cfg.JWTSecret, s.cfg.JWTExpiry)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": u.ID.String(),
		"token":   token,
	})
}

type ctxKey int

const userIDKey ctxKey = 1

func (s *Server) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			writeJSON(w, http.StatusUnauthorized, errResp{"missing_token"})
			return
		}
		raw := strings.TrimSpace(h[7:])
		claims, err := auth.ParseToken(raw, s.cfg.JWTSecret)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, errResp{"invalid_token"})
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func userIDFrom(ctx context.Context) uuid.UUID {
	v, _ := ctx.Value(userIDKey).(uuid.UUID)
	return v
}

func (s *Server) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.AdminToken == "" {
			writeJSON(w, http.StatusForbidden, errResp{"admin_disabled"})
			return
		}
		tok := strings.TrimSpace(r.Header.Get("X-Admin-Token"))
		if tok == "" {
			h := r.Header.Get("Authorization")
			if strings.HasPrefix(strings.ToLower(h), "bearer ") {
				tok = strings.TrimSpace(h[7:])
			}
		}
		if tok != s.cfg.AdminToken {
			writeJSON(w, http.StatusUnauthorized, errResp{"unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- profile ---

type profileDTO struct {
	SourceLanguage   string `json:"source_language"`
	TargetLanguage   string `json:"target_language"`
	ProficiencyLevel string `json:"proficiency_level"`
	LearningGoal     string `json:"learning_goal"`
	TutorStyle       string `json:"tutor_style"`
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	p, err := s.store.ProfileByUserID(r.Context(), uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if p == nil {
		writeJSON(w, http.StatusOK, profileDTO{
			SourceLanguage:   "en",
			TargetLanguage:   "",
			ProficiencyLevel: "A2",
			LearningGoal:     "daily_life",
			TutorStyle:       "gentle",
		})
		return
	}
	writeJSON(w, http.StatusOK, profileDTO{
		SourceLanguage:   p.SourceLanguage,
		TargetLanguage:   p.TargetLanguage,
		ProficiencyLevel: p.ProficiencyLevel,
		LearningGoal:     p.LearningGoal,
		TutorStyle:       p.TutorStyle,
	})
}

func (s *Server) putProfile(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	var dto profileDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	if dto.TargetLanguage == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"target_language_required"})
		return
	}
	p := &store.Profile{
		UserID:           uid,
		SourceLanguage:   dto.SourceLanguage,
		TargetLanguage:   dto.TargetLanguage,
		ProficiencyLevel: dto.ProficiencyLevel,
		LearningGoal:     dto.LearningGoal,
		TutorStyle:       dto.TutorStyle,
	}
	if p.SourceLanguage == "" {
		p.SourceLanguage = "en"
	}
	if p.ProficiencyLevel == "" {
		p.ProficiencyLevel = "A2"
	}
	if p.LearningGoal == "" {
		p.LearningGoal = "daily_life"
	}
	if p.TutorStyle == "" {
		p.TutorStyle = "gentle"
	}
	if err := s.store.UpsertProfile(r.Context(), p); err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) getScenarios(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"scenarios": listScenarios()})
}

// --- experiments & flags ---

func (s *Server) ensureAssignments(ctx context.Context, uid uuid.UUID) (map[string]string, error) {
	exps, err := s.store.ListActiveExperiments(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, e := range exps {
		v, err := s.store.GetExperimentAssignment(ctx, uid, e.Key)
		if err != nil {
			return nil, err
		}
		if v == "" {
			v = experiment.VariantForUser(uid.String(), e.Key, e.Variants)
			if err := s.store.SetExperimentAssignment(ctx, uid, e.Key, v); err != nil {
				return nil, err
			}
		}
		out[e.Key] = v
	}
	return out, nil
}

func (s *Server) getExperiments(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	m, err := s.ensureAssignments(r.Context(), uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"assignments": m})
}

func (s *Server) getFeatureFlags(w http.ResponseWriter, r *http.Request) {
	m, err := s.store.ListFeatureFlags(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"flags": m})
}

// --- sessions ---

type createSessionReq struct {
	ScenarioID string `json:"scenario_id"`
}

// sessionListItem adds catalog metadata to each row from ListSessionsForUser.
type sessionListItem struct {
	store.Session
	ScenarioTitle string `json:"scenario_title"`
}

type createSessionResp struct {
	SessionID           string                 `json:"session_id"`
	RoomName            string                 `json:"room_name"`
	LiveKitURL          string                 `json:"livekit_url"`
	LiveKitToken        string                 `json:"livekit_token,omitempty"`
	PromptVersion       string                 `json:"prompt_version"`
	ExperimentSnapshot  map[string]string      `json:"experiment_snapshot"`
	TokenTTLSeconds     int                    `json:"token_ttl_seconds,omitempty"`
	Note                string                 `json:"note,omitempty"`
}

func (s *Server) createSession(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	var req createSessionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	req.ScenarioID = strings.TrimSpace(req.ScenarioID)
	if req.ScenarioID == "" {
		writeJSON(w, http.StatusBadRequest, errResp{"scenario_id_required"})
		return
	}
	okScenario := false
	for _, sc := range catalog {
		if sc.ID == req.ScenarioID {
			okScenario = true
			break
		}
	}
	if !okScenario {
		writeJSON(w, http.StatusBadRequest, errResp{"unknown_scenario"})
		return
	}
	pv, err := s.store.ActivePromptVersion(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	assign, err := s.ensureAssignments(r.Context(), uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	snap, _ := json.Marshal(assign)
	room := "ff-" + uuid.New().String()
	sess, err := s.store.CreateSession(r.Context(), uid, req.ScenarioID, room, pv, snap)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	_ = s.store.InsertSessionEvent(r.Context(), sess.ID, analytics.SessionCreated, json.RawMessage(`{}`))
	_ = s.store.InsertSessionEvent(r.Context(), sess.ID, analytics.ExperimentExposed, snap)
	if s.sessionEvents != nil {
		s.sessionEvents.WithLabelValues(analytics.SessionCreated).Inc()
		s.sessionEvents.WithLabelValues(analytics.ExperimentExposed).Inc()
	}

	ttl := 15 * time.Minute
	resp := createSessionResp{
		SessionID:          sess.ID.String(),
		RoomName:           room,
		LiveKitURL:         s.cfg.LiveKitURL,
		PromptVersion:      pv,
		ExperimentSnapshot: assign,
		TokenTTLSeconds:    int(ttl.Seconds()),
	}
	if s.cfg.LiveKitAPIKey != "" && s.cfg.LiveKitAPISecret != "" {
		tok, err := s.mintLearnerJoinToken(r.Context(), room, uid, sess)
		if err != nil {
			log.Printf("livekit token: %v", err)
			resp.Note = "livekit_token_failed_check_credentials"
		} else {
			resp.LiveKitToken = tok
		}
	} else {
		resp.Note = "configure LIVEKIT_URL LIVEKIT_API_KEY LIVEKIT_API_SECRET for join tokens"
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	limit := 20
	if v := strings.TrimSpace(r.URL.Query().Get("limit")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	list, err := s.store.ListSessionsForUser(r.Context(), uid, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	out := make([]sessionListItem, 0, len(list))
	for _, sess := range list {
		out = append(out, sessionListItem{
			Session:       sess,
			ScenarioTitle: ScenarioTitle(sess.ScenarioID),
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"sessions": out})
}

func (s *Server) getSession(w http.ResponseWriter, r *http.Request) {
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
	var ended interface{}
	if sess.EndedAt != nil {
		ended = sess.EndedAt.UTC().Format(time.RFC3339)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":      sess.ID.String(),
		"scenario_id":     sess.ScenarioID,
		"scenario_title":  ScenarioTitle(sess.ScenarioID),
		"room_name":       sess.RoomName,
		"livekit_url":     s.cfg.LiveKitURL,
		"status":          sess.Status,
		"prompt_version":  sess.PromptVersion,
		"started_at":      sess.StartedAt.UTC().Format(time.RFC3339),
		"ended_at":        ended,
	})
}

func (s *Server) issueLiveKitToken(w http.ResponseWriter, r *http.Request) {
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
	if sess.Status != "active" {
		writeJSON(w, http.StatusConflict, errResp{"session_not_active"})
		return
	}
	if s.cfg.LiveKitAPIKey == "" || s.cfg.LiveKitAPISecret == "" {
		writeJSON(w, http.StatusServiceUnavailable, errResp{"livekit_not_configured"})
		return
	}
	ttl := 15 * time.Minute
	tok, err := s.mintLearnerJoinToken(r.Context(), sess.RoomName, uid, sess)
	if err != nil {
		log.Printf("livekit token refresh: %v", err)
		writeJSON(w, http.StatusInternalServerError, errResp{"token_issue_failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"livekit_token":     tok,
		"livekit_url":       s.cfg.LiveKitURL,
		"token_ttl_seconds": int(ttl.Seconds()),
	})
}

type batchEventsReq struct {
	Events []struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	} `json:"events"`
}

func (s *Server) postSessionEvents(w http.ResponseWriter, r *http.Request) {
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
	var req batchEventsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{"invalid_json"})
		return
	}
	for _, e := range req.Events {
		if e.Type == "" {
			continue
		}
		payload := e.Payload
		if len(payload) == 0 {
			payload = json.RawMessage(`{}`)
		}
		if err := s.store.InsertSessionEvent(r.Context(), sess.ID, e.Type, payload); err != nil {
			writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
			return
		}
		if s.sessionEvents != nil {
			s.sessionEvents.WithLabelValues(e.Type).Inc()
		}
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "accepted"})
}

type completeReq struct {
	SpeakingSeconds int `json:"speaking_seconds"`
	TurnCount       int `json:"turn_count"`
}

func (s *Server) completeSession(w http.ResponseWriter, r *http.Request) {
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
	var body completeReq
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.SpeakingSeconds < 0 {
		body.SpeakingSeconds = 0
	}
	if err := s.store.CompleteSession(r.Context(), id, body.SpeakingSeconds, body.TurnCount); err != nil {
		writeJSON(w, http.StatusConflict, errResp{"session_already_completed_or_missing"})
		return
	}
	_ = s.store.InsertSessionEvent(r.Context(), id, analytics.SessionCompleted, json.RawMessage(`{}`))
	if s.sessionEvents != nil {
		s.sessionEvents.WithLabelValues(analytics.SessionCompleted).Inc()
	}
	if err := s.store.AppendLearningMetricSnapshot(r.Context(), uid); err != nil {
		log.Printf("learning metric snapshot: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (s *Server) getFeedback(w http.ResponseWriter, r *http.Request) {
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
	f, err := s.store.FeedbackBySession(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if f == nil {
		writeJSON(w, http.StatusNotFound, errResp{"feedback_not_ready"})
		return
	}
	writeJSON(w, http.StatusOK, enrichFeedback(f))
}

func (s *Server) generateFeedback(w http.ResponseWriter, r *http.Request) {
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
	if sess.Status != "completed" {
		writeJSON(w, http.StatusConflict, errResp{"complete_session_first"})
		return
	}
	ctx := r.Context()
	prof, _ := s.store.ProfileByUserID(ctx, uid)
	transcript, _ := s.store.TranscriptPlainText(ctx, id)
	var f *store.FeedbackSummary
	var src string
	if s.cfg.OpenAIAPIKey != "" {
		llm, err := openaipkg.SessionFeedback(ctx, s.cfg.OpenAIAPIKey, id, sess.ScenarioID, prof, transcript)
		if err != nil {
			log.Printf("openai feedback: %v", err)
		}
		if llm != nil {
			f = llm
			src = "openai"
		}
	}
	if f == nil {
		f = buildStubFeedback(id, sess.ScenarioID)
		src = "stub"
	}
	if err := s.store.UpsertFeedback(ctx, f); err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	evPayload, _ := json.Marshal(map[string]string{"source": src})
	_ = s.store.InsertSessionEvent(r.Context(), id, analytics.CorrectionGenerated, evPayload)
	if s.sessionEvents != nil {
		s.sessionEvents.WithLabelValues(analytics.CorrectionGenerated).Inc()
	}
	out, err := s.store.FeedbackBySession(ctx, id)
	if err != nil || out == nil {
		writeJSON(w, http.StatusOK, enrichFeedback(f))
		return
	}
	writeJSON(w, http.StatusOK, enrichFeedback(out))
}

func (s *Server) dashboardSummary(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	d, err := s.store.DashboardForUser(r.Context(), uid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (s *Server) adminOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var users, sessions, events int
	_ = s.store.Pool().QueryRow(ctx, `SELECT COUNT(*)::int FROM users`).Scan(&users)
	_ = s.store.Pool().QueryRow(ctx, `SELECT COUNT(*)::int FROM sessions`).Scan(&sessions)
	_ = s.store.Pool().QueryRow(ctx, `SELECT COUNT(*)::int FROM session_events`).Scan(&events)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":          users,
		"sessions":       sessions,
		"session_events": events,
	})
}
