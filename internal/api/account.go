package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mehdi/fluentflow/internal/auth"
	"github.com/mehdi/fluentflow/internal/store"
)

type deleteAccountReq struct {
	Password string `json:"password"`
}

// deleteAccount removes the current user and cascades related rows (PRD data deletion).
// Registered users must send the correct password; guests may omit it.
func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	ctx := r.Context()
	u, err := s.store.UserByID(ctx, uid)
	if err != nil || u == nil {
		writeJSON(w, http.StatusUnauthorized, errResp{"unauthorized"})
		return
	}
	var body deleteAccountReq
	_ = json.NewDecoder(r.Body).Decode(&body)
	if !u.IsGuest {
		if !auth.CheckPassword(u.PasswordHash, body.Password) {
			writeJSON(w, http.StatusUnauthorized, errResp{"invalid_password"})
			return
		}
	}
	if err := s.store.DeleteUser(ctx, uid); err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) listLearningSnapshots(w http.ResponseWriter, r *http.Request) {
	uid := userIDFrom(r.Context())
	limit := 30
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	list, err := s.store.ListLearningMetricSnapshots(r.Context(), uid, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errResp{"server_error"})
		return
	}
	if list == nil {
		list = []store.LearningMetricSnapshot{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"snapshots": list})
}
