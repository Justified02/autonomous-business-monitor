package scheduler

import (
	"context"
	"fmt"

	"time"

	"github.com/Justified02/abm/internal/anomaly"
	"github.com/Justified02/abm/internal/fetcher"
	"github.com/Justified02/abm/internal/llm"
	"github.com/Justified02/abm/internal/storage"
	"github.com/Justified02/abm/internal/storage/db"
	//"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
)

type fetchResult struct {
	data []byte
	err  error
}

type Scheduler struct {
	stripe *fetcher.StripeClient
	engine *anomaly.Engine
	llm    *llm.LLMClient
	db     *storage.Store
}

func NewScheduler(s *fetcher.StripeClient, e *anomaly.Engine, l *llm.LLMClient, db  *storage.Store) *Scheduler {
	newScheduler := &Scheduler{
		stripe: s,
		engine: e,
		llm:    l,
		db:		db,
	}

	return newScheduler
}

// every morning, the scheduler needs to collect data from all sources (Stripe, Gmail, Calendly)
// FetchAllSources means "run the data collection process"
func (s *Scheduler) FetchAllSources(ctx context.Context) {
	fmt.Println("starting fetch all sources...")
	fetchedResult := make(chan fetchResult, 1)

	go func() {
		data, err := s.stripe.Fetch(ctx)
		fetchedResult <- fetchResult{data: data, err: err}
	}()

	result := <-fetchedResult

	if result.err != nil {
		fmt.Println("error fetching data:", result.err)
		return
	}

	// Save the stripe snapshot
	err := s.stripe.Save(ctx, result.data)
	if err != nil {
		fmt.Println("error saving snapshot:", err)
		return
	}

	// Parse the raw data
	revenue, failedCounts, err := s.stripe.Parse(result.data)
	if err != nil {
		fmt.Println("error parsing raw data:", err)
		return
	}

	fmt.Println("revenue:", revenue)
	fmt.Println("failedCounts:", failedCounts)

	// save to daily_metrics
	// var pgRevenue pgtype.Numeric
	// pgRevenue.Scan(fmt.Sprintf("%.2f", revenue))

	// _, err = s.db.Queries().SaveDailyMetrics(ctx, db.SaveDailyMetricsParams{
	// 	Source: "stripe",
	// 	MetricDate: pgtype.Date{Time: time.Now(), Valid: true},
	// 	Revenue: pgRevenue,
	// 	FailedPayments: int32(failedCounts),
	// })
	// if err != nil {
	// 	fmt.Println("error saving daily metrics:", err)
	// 	return
	// }

	// run the anomaly engine
	anomResult, err := s.engine.Analyze(ctx, "stripe", revenue)
	if err != nil {
		fmt.Println("error running anomaly engine:", err)
		return
	}

	// Build the prompt
	prompt := fmt.Sprintf(
		"Generate a morning briefing. Stripe revenue today: $%.2f. Failed payments: %d. Anomaly detected: %v. Delta from 7-day average: %.2f%%.",
		revenue, failedCounts, anomResult.IsAnomaly, anomResult.DeltaPct,
	)

	// call the llm
	llmResp, err := s.llm.Generate(ctx, prompt)
	if err != nil {
		fmt.Println("error generating briefing:", err)
		return
	}

	// save the digest to the db
	_, err = s.db.Queries().SaveDigest(ctx, db.SaveDigestParams{
		Content: llmResp,
		HasCriticalAlerts: anomResult.IsAnomaly,
	})
	if err != nil {
		fmt.Println("error saving digest:", err)
		return
	}

	fmt.Println("stripe fetch complete, err:", result.err)
}

// Start the cron job - call the fetchAllSources function in the process
func (s *Scheduler) Start(cronSchedule string) {
	// create a cron scheduler
	c := cron.New()

	// add a cron job
	_, err := c.AddFunc(cronSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		s.FetchAllSources(ctx)
	})
	if err != nil {
		fmt.Println("adding cron job:", err)
		return
	}

	c.Start()
	select {}
}
