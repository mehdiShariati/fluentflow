package config

import "testing"

func TestLoad_requiresDatabaseAndJWT(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoad_ok(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x:y@localhost:5432/db")
	t.Setenv("JWT_SECRET", "1234567890123456")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.DatabaseURL == "" {
		t.Fatal("database url")
	}
}
