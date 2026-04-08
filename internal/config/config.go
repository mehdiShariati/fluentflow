package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr         string
	DatabaseURL      string
	JWTSecret        string
	JWTExpiry        time.Duration
	LiveKitURL       string
	LiveKitAPIKey    string
	LiveKitAPISecret string
	// LiveKitAgentName must match the Python worker @rtc_session(agent_name=...). Empty disables dispatch in tokens.
	LiveKitAgentName string
	// OpenAIAPIKey enables LLM-generated post-session feedback (optional; stub used if empty).
	OpenAIAPIKey string
	AdminToken   string
	CorsOrigins  []string
}

func Load() (*Config, error) {
	jwtExp := 24 * time.Hour
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		h, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("JWT_EXPIRY_HOURS: %w", err)
		}
		jwtExp = time.Duration(h) * time.Hour
	}
	origins := []string{"http://localhost:3000"}
	if v := strings.TrimSpace(os.Getenv("CORS_ORIGINS")); v != "" {
		var parsed []string
		for _, p := range strings.Split(v, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				parsed = append(parsed, p)
			}
		}
		if len(parsed) > 0 {
			origins = parsed
		}
	}
	c := &Config{
		HTTPAddr:         getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpiry:        jwtExp,
		LiveKitURL:       os.Getenv("LIVEKIT_URL"),
		LiveKitAPIKey:    os.Getenv("LIVEKIT_API_KEY"),
		LiveKitAPISecret: os.Getenv("LIVEKIT_API_SECRET"),
		LiveKitAgentName: strings.TrimSpace(getenv("LIVEKIT_AGENT_NAME", "fluentflow-tutor")),
		OpenAIAPIKey:     strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		AdminToken:       os.Getenv("ADMIN_TOKEN"),
		CorsOrigins:      origins,
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" || len(c.JWTSecret) < 16 {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least 16 characters")
	}
	return c, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
