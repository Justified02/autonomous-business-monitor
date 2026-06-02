package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"honnef.co/go/tools/config"
)

func main() {
	// 1. Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		level: slog.LevelInfo,
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


	// 4. Build the dependency chain - each component receives exactly what it needs - nothing more
	stripeClient := fetcher.NewStripe(cfg.StripeKey)

	llmClient, err := llm.New(cfg.LLMKey, "prompts/morning_digest.txt")
	if err != nil {
		slog.Error("failed to initialise llm client", "error", err)
		os.Exit(1)
	}

	sched := scheduler.New(store, stripeClient, llmClient)


	// 5. Run mode - RUN_NOW=true skips the cron schedule and fires immediately
	if os.Getenv("RUN_NOW") == "true" {
		slog.Info("RUN_NOW=true - running collection immediately")
		err := sched.RunNow(ctx)
		if err != nil {
			slog.Error("collection failed", "error", err)
			os.Exit(1)
		}
		return
	}

	// 6. Start the scheduler
	err = sched.Start(cfg.CronSchedule)
	if err != nil {
		slog.Error("failed to start scheduler", "error", err)
		os.Exit(1)
	}
	defer sched.Stop()

	slog.Info("scheduler running", "schedule", cfg.CronSchedule)


	// 7. Graceful shutdown - blocks the main goroutine - keeps the program alive - until
	// it receives SIGNIT (Ctrl+C from you) or SIGTERM (from Railway when it deploys a new version)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())
	slog.Info("shutdown complete")
}
