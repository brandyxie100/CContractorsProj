package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration from environment variables.
type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	AccessTTL          time.Duration
	RefreshTTL         time.Duration
	CORSOrigin         string
	SeedAdminEmail     string
	SeedAdminPassword  string
}

// Load reads .env (if present) and environment variables.
func Load() Config {
	_ = godotenv.Load()
	accessMin := envInt("JWT_ACCESS_TTL_MINUTES", 15)
	refreshHours := envInt("JWT_REFRESH_TTL_HOURS", 168)
	return Config{
		Port:              env("PORT", "8080"),
		DatabaseURL:       env("DATABASE_URL", "postgres://clements:clements@localhost:5432/equipment?sslmode=disable"),
		JWTSecret:         env("JWT_SECRET", "dev-change-me-in-production-use-long-secret"),
		AccessTTL:         time.Duration(accessMin) * time.Minute,
		RefreshTTL:        time.Duration(refreshHours) * time.Hour,
		CORSOrigin:        env("CORS_ORIGIN", "http://localhost:5173"),
		SeedAdminEmail:    env("SEED_ADMIN_EMAIL", "admin@clements.local"),
		SeedAdminPassword: env("SEED_ADMIN_PASSWORD", "admin123"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
