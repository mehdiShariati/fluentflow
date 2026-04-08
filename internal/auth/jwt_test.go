package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSignParseRoundTrip(t *testing.T) {
	uid := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	secret := "test-secret-at-least-16"
	tok, err := SignToken(uid, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseToken(tok, secret)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != uid {
		t.Fatalf("uid: got %v want %v", claims.UserID, uid)
	}
}

func TestParseToken_wrongSecret(t *testing.T) {
	uid := uuid.New()
	tok, _ := SignToken(uid, "secret-one-at-least-16", time.Hour)
	_, err := ParseToken(tok, "secret-two-at-least-16")
	if err == nil {
		t.Fatal("expected error")
	}
}
