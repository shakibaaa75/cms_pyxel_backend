package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	FRONTEND_URL        = os.Getenv("FRONTEND_URL")
	SECOND_FRONTEND_URL = os.Getenv("SECOND_FRONTEND_URL")
	ADMIN_URL           = os.Getenv("ADMIN_URL")

	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USERNAME string
	SMTP_PASSWORD string
	SMTP_FROM     string

	ADMIN_USERNAME     string
	ADMIN_PASSWORD_RAW string
	AdminPasswordHash  []byte

	JWT_SECRET []byte

	COMPANY_LOGO_URL string
	COMPANY_NAME     string
	COMPANY_PHONE    string
	COMPANY_EMAIL    string
	COMPANY_WEBSITE  string
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	FRONTEND_URL = getEnv("FRONTEND_URL", "http://localhost:5174")
	SECOND_FRONTEND_URL = getEnv("SECOND_FRONTEND_URL", "http://localhost:5175")
	ADMIN_URL = getEnv("ADMIN_URL", "http://localhost:5173")

	SMTP_HOST = getEnv("SMTP_HOST", "")
	SMTP_PORT = getEnv("SMTP_PORT", "587")
	SMTP_USERNAME = getEnv("SMTP_USERNAME", "")
	SMTP_PASSWORD = getEnv("SMTP_PASSWORD", "")
	SMTP_FROM = getEnv("SMTP_FROM", "Pyxel Construction")

	ADMIN_USERNAME = getEnv("ADMIN_USERNAME", "")
	ADMIN_PASSWORD_RAW = getEnv("ADMIN_PASSWORD", "")

	jwtSecret := getEnv("JWT_SECRET", "")

	COMPANY_LOGO_URL = getEnv("COMPANY_LOGO_URL", "")
	COMPANY_NAME = getEnv("COMPANY_NAME", "Pyxel Construction")
	COMPANY_PHONE = getEnv("COMPANY_PHONE", "(916) 888-8281")
	COMPANY_EMAIL = getEnv("COMPANY_EMAIL", "contact@pyxelconstruction.com")
	COMPANY_WEBSITE = getEnv("COMPANY_WEBSITE", "https://pyxelconstruction.com")

	required := []struct{ name, value string }{
		{"SMTP_HOST", SMTP_HOST},
		{"SMTP_USERNAME", SMTP_USERNAME},
		{"SMTP_PASSWORD", SMTP_PASSWORD},
		{"ADMIN_USERNAME", ADMIN_USERNAME},
		{"ADMIN_PASSWORD", ADMIN_PASSWORD_RAW},
		{"JWT_SECRET", jwtSecret},
	}
	for _, r := range required {
		if r.value == "" {
			slog.Error("required environment variable is missing", "name", r.name)
			os.Exit(1)
		}
	}

	if len(jwtSecret) < 32 {
		slog.Error("JWT_SECRET must be at least 32 characters")
		os.Exit(1)
	}
	JWT_SECRET = []byte(jwtSecret)

	hash, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PASSWORD_RAW), 12)
	if err != nil {
		slog.Error("failed to hash admin password", "err", err)
		os.Exit(1)
	}
	AdminPasswordHash = hash
	ADMIN_PASSWORD_RAW = ""

	slog.Info("company loaded", "name", COMPANY_NAME)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
func GetEnv(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}
