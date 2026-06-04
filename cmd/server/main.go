package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Justified02/abm/config"
	"github.com/Justified02/abm/internal/anomaly"
	"github.com/Justified02/abm/internal/fetcher"
	"github.com/Justified02/abm/internal/llm"
	"github.com/Justified02/abm/internal/scheduler"
	"github.com/Justified02/abm/internal/storage"
)

func main() {
	// 1. Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// 2. Configuration - load all environment variables into one typed struct
	// if any required variable is missing, this exits immediately with a clear error message
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded", "port", cfg.Port, "schedule", cfg.CronSchedule)

	// 3. Database - context.Background() is the root context
	ctx := context.Background()

	store, err := storage.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	slog.Info("database connected")

	// 4. Create all clients
	stripeClient := fetcher.NewStripeClient(cfg.StripeKey, store)
	engineClient := anomaly.NewEngine(store.Queries())
	llmClient := llm.NewLLMClient(cfg.LLMKey, cfg.LLMModel)

	// 5. Pass the stripeClient to the scheduler to create a new scheduler
	newScheduler := scheduler.NewScheduler(stripeClient, engineClient, llmClient, store)

	// 6. Start the scheduler
	newScheduler.Start(cfg.CronSchedule)
}
