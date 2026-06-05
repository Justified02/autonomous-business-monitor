package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	StripeKey     string
	LLMKey        string
	CronSchedule  string
	N8NWebhookURL string
	Port          string
	LLMModel	  string
	RunNow		  bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	runNow, _ := strconv.ParseBool(os.Getenv("RUN_NOW"))

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		StripeKey:     os.Getenv("STRIPE_SECRET_KEY"),
		LLMKey:        os.Getenv("LLM_API_KEY"),
		CronSchedule:  os.Getenv("CRON_SCHEDULE"),
		N8NWebhookURL: os.Getenv("N8N_WEBHOOK_URL"),
		LLMModel: 	   os.Getenv("LLM_MODEL"),
		RunNow:		   runNow,
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
	if cfg.LLMModel == "" {
		return nil, fmt.Errorf("LLM_MODEL is required")
	}

	// Optional fields get sensible defaults
	if cfg.CronSchedule == "" {
		cfg.CronSchedule = "45 6 * * *"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	// if cfg.RunNow == false {
	// 	cfg.RunNow = false
	// }

	return cfg, nil
}
