package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	StripeKey     string
	LLMKey        string
	CronSchedule  string
	N8NWebhookURL string
	Port          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		StripeKey:     os.Getenv("STRIPE_SECRET_KEY"),
		LLMKey:        os.Getenv("LLM_API_KEY"),
		CronSchedule:  os.Getenv("CRON_SCHEDULE"),
		N8NWebhookURL: os.Getenv("N8N_WEBHOOK_URL"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.StripeKey == "" {
		return nil, fmt.Errorf("STRIPE_SECRET_KEY is required")
	}
	if cfg.LLMKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY is required")
	}

	// Optional fields get sensible defaults
	if cfg.CronSchedule == "" {
		cfg.CronSchedule = "45 6 * * *"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
