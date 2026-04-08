package livekit

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestJoinToken_roomConfigWhenAgentSet(t *testing.T) {
	tok, err := JoinToken(JoinTokenParams{
		APIKey:        "devkey",
		APISecret:     "secret",
		Room:          "ff-test-room",
		Identity:      "learner-1",
		TTL:           10 * time.Minute,
		DisplayName:   "Alex",
		AgentName:     "fluentflow-tutor",
		AgentMetadata: `{"scenario_id":"coffee_shop"}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(tok, ".")
	if len(parts) < 2 {
		t.Fatal("invalid jwt")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	rc, ok := payload["roomConfig"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing roomConfig: %#v", payload)
	}
	agents, ok := rc["agents"].([]interface{})
	if !ok || len(agents) == 0 {
		t.Fatalf("missing agents: %#v", rc)
	}
	ag0 := agents[0].(map[string]interface{})
	if ag0["agentName"] != "fluentflow-tutor" {
		t.Fatalf("agentName: %#v", ag0)
	}
}

func TestJoinToken_noAgentOmitsRoomConfig(t *testing.T) {
	tok, err := JoinToken(JoinTokenParams{
		APIKey:    "k",
		APISecret: "secretsecretsecret",
		Room:      "r1",
		Identity:  "u1",
		TTL:       time.Minute,
		AgentName: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(tok, ".")
	raw, _ := base64.RawURLEncoding.DecodeString(parts[1])
	var payload map[string]interface{}
	_ = json.Unmarshal(raw, &payload)
	if _, ok := payload["roomConfig"]; ok {
		t.Fatal("expected no roomConfig")
	}
}
