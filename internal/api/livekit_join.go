package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mehdi/fluentflow/internal/livekit"
	"github.com/mehdi/fluentflow/internal/store"
)

func joinDisplayName(ctx context.Context, st *store.Store, uid uuid.UUID) string {
	u, err := st.UserByID(ctx, uid)
	if err != nil || u == nil {
		return "Learner"
	}
	parts := strings.Split(u.Email, "@")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "Learner"
}

func joinAgentMetadataJSON(sess *store.Session, prof *store.Profile) string {
	m := map[string]interface{}{
		"scenario_id":    sess.ScenarioID,
		"session_id":     sess.ID.String(),
		"prompt_version": sess.PromptVersion,
	}
	if prof != nil {
		m["target_language"] = prof.TargetLanguage
		m["proficiency_level"] = prof.ProficiencyLevel
		m["learning_goal"] = prof.LearningGoal
		m["tutor_style"] = prof.TutorStyle
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func (s *Server) mintLearnerJoinToken(ctx context.Context, room string, uid uuid.UUID, sess *store.Session) (string, error) {
	prof, _ := s.store.ProfileByUserID(ctx, uid)
	return livekit.JoinToken(livekit.JoinTokenParams{
		APIKey:        s.cfg.LiveKitAPIKey,
		APISecret:     s.cfg.LiveKitAPISecret,
		Room:          room,
		Identity:      uid.String(),
		TTL:           15 * time.Minute,
		DisplayName:   joinDisplayName(ctx, s.store, uid),
		AgentName:     s.cfg.LiveKitAgentName,
		AgentMetadata: joinAgentMetadataJSON(sess, prof),
	})
}
