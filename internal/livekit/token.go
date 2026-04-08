package livekit

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// VideoGrant matches LiveKit JWT video grant claims (subset for room join).
type VideoGrant struct {
	RoomJoin     bool   `json:"roomJoin"`
	Room         string `json:"room"`
	CanPublish   bool   `json:"canPublish"`
	CanSubscribe bool   `json:"canSubscribe"`
}

// roomAgentDispatch is embedded in join tokens so the server dispatches the tutor when the learner connects.
// See https://docs.livekit.io/agents/server/agent-dispatch/
type roomAgentDispatch struct {
	AgentName string `json:"agentName"`
	Metadata  string `json:"metadata,omitempty"`
}

type roomConfiguration struct {
	Agents []roomAgentDispatch `json:"agents"`
}

type livekitClaims struct {
	Video      *VideoGrant        `json:"video,omitempty"`
	RoomConfig *roomConfiguration `json:"roomConfig,omitempty"`
	Name       string             `json:"name,omitempty"`
	Metadata   string             `json:"metadata,omitempty"`
	jwt.RegisteredClaims
}

// JoinTokenParams configures a learner room join token (optional tutor dispatch via roomConfig).
type JoinTokenParams struct {
	APIKey    string
	APISecret string
	Room      string
	Identity  string
	TTL       time.Duration
	// DisplayName is shown as participant name in the room.
	DisplayName string
	// ParticipantMetadata is optional participant metadata (JSON string).
	ParticipantMetadata string
	// AgentName triggers RoomAgentDispatch when non-empty (must match worker agent_name).
	AgentName string
	// AgentMetadata is JSON passed to the agent job (scenario, level, etc.).
	AgentMetadata string
}

// JoinToken builds a short-lived JWT for joining a LiveKit room (HS256).
// If AgentName is set, the agent is dispatched when this participant connects.
func JoinToken(p JoinTokenParams) (string, error) {
	if p.APIKey == "" || p.APISecret == "" {
		return "", errors.New("livekit: missing API key or secret")
	}
	now := time.Now()
	claims := livekitClaims{
		Video: &VideoGrant{
			RoomJoin:     true,
			Room:         p.Room,
			CanPublish:   true,
			CanSubscribe: true,
		},
		Name:     p.DisplayName,
		Metadata: p.ParticipantMetadata,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.APIKey,
			Subject:   p.Identity,
			ExpiresAt: jwt.NewNumericDate(now.Add(p.TTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	}
	if p.AgentName != "" {
		claims.RoomConfig = &roomConfiguration{
			Agents: []roomAgentDispatch{
				{AgentName: p.AgentName, Metadata: p.AgentMetadata},
			},
		}
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return t.SignedString([]byte(p.APISecret))
}

// DecodeGrantJSON is a helper for debugging (not used in hot path).
func DecodeGrantJSON(token string) (string, error) {
	parts := jwt.Parser{}
	_, _, err := parts.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	segs := splitJWT(token)
	if len(segs) < 2 {
		return "", errors.New("invalid jwt")
	}
	raw, err := base64.RawURLEncoding.DecodeString(segs[1])
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return "", err
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b), nil
}

func splitJWT(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
