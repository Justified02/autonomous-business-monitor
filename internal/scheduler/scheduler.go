package scheduler

import (
	"context"
	"fmt"

	"time"

	"github.com/Justified02/abm/internal/anomaly"
	"github.com/Justified02/abm/internal/fetcher"
	"github.com/Justified02/abm/internal/llm"
	"github.com/Justified02/abm/internal/notify"
	"github.com/Justified02/abm/internal/storage"
	"github.com/Justified02/abm/internal/storage/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
)

type fetchResult struct {
	source string
	data   []byte
	err    error
}

type Scheduler struct {
	stripe   *fetcher.StripeClient
	gmail    *fetcher.GmailClient
	calendly *fetcher.CalendlyClient
	engine   *anomaly.Engine
	llm      *llm.LLMClient
	db       *storage.Store
	webhook  *notify.WebhookClient
}

func NewScheduler(s *fetcher.StripeClient, e *anomaly.Engine, l *llm.LLMClient, db *storage.Store, g *fetcher.GmailClient, c *fetcher.CalendlyClient, w *notify.WebhookClient) *Scheduler {
	newScheduler := &Scheduler{
		stripe:   s,
		engine:   e,
		llm:      l,
		db:       db,
		gmail:    g,
		calendly: c,
		webhook: w,
	}

	return newScheduler
}

// every morning, the scheduler needs to collect data from all sources (Stripe, Gmail, Calendly)
// FetchAllSources means "run the data collection process"
func (s *Scheduler) FetchAllSources(ctx context.Context) {
	fmt.Println("starting fetch all sources...")
	fetchedResult := make(chan fetchResult, 3)

	go func() {
		data, err := s.stripe.Fetch(ctx)
		fetchedResult <- fetchResult{source: "stripe", data: data, err: err}
	}()

	go func() {
		data, err := s.gmail.Fetch(ctx)
		fetchedResult <- fetchResult{source: "gmail", data: data, err: err}
	}()

	go func() {
		data, err := s.calendly.Fetch(ctx)
		fetchedResult <- fetchResult{source: "calendly", data: data, err: err}
	}()

	var stripeRevenue float64
	var stripeFailedPayments int
	var stripeAnomaly anomaly.AnomalyResult
	var calendlyEvents []byte
	var gmailMessages []byte

	for i := 0; i < 3; i++ {
		result := <-fetchedResult
		if result.err != nil {
			fmt.Println(result.source, "fetch error:", result.err)
			continue // skip this source, don't return
		}

		if result.source == "stripe" {
			fmt.Println("processing stripe data...")
			// save the stripe snapshot
			err := s.stripe.Save(ctx, result.data)
			if err != nil {
				fmt.Println("error saving snapshot:", err)
				continue
			}
			fmt.Println("stripe snapshot saved successfully")

			// parse the stripe raw data
			revenue, failedCounts, err := s.stripe.Parse(result.data)
			if err != nil {
				fmt.Println("error parsing raw data:", err)
				continue
			}

			// save to daily_metrics
			var pgRevenue pgtype.Numeric
			pgRevenue.Scan(fmt.Sprintf("%.2f", revenue))

			_, err = s.db.Queries().SaveDailyMetrics(ctx, db.SaveDailyMetricsParams{
				Source:         "stripe",
				MetricDate:     pgtype.Date{Time: time.Now(), Valid: true},
				Revenue:        pgRevenue,
				FailedPayments: int32(failedCounts),
			})
			if err != nil {
				fmt.Println("error saving daily metrics:", err)
				continue
			}

			// run the anomaly engine
			anomResult, err := s.engine.Analyze(ctx, "stripe", revenue)
			if err != nil {
				fmt.Println("error running anomaly engine:", err)
				continue
			}

			stripeRevenue = revenue
			stripeFailedPayments = failedCounts
			stripeAnomaly = anomResult

			fmt.Println("stripe fetch complete, err:", result.err)

		} else if result.source == "gmail" {
			fmt.Println("gmail data received, length:", len(result.data))

			// save gmail snapshot
			err := s.gmail.Save(ctx, result.data)
			if err != nil {
				fmt.Println("error saving gmail snapshot:", err)
				continue
			}
			fmt.Println("gmail snapshot saved successfully")

			gmailMessages = result.data

		} else if result.source == "calendly" {
			fmt.Println("calendly data received, length:", len(result.data))

			// save calendly snapshot
			err := s.calendly.Save(ctx, result.data)
			if err != nil {
				fmt.Println("error saving calendly snashot:", err)
				continue
			}
			fmt.Println("calendly snapshot saved sucessfully")

			calendlyEvents = result.data
		}
	}

	digestData := llm.DigestData{
		StripeRevenue:        stripeRevenue,
		StripeFailedPayments: stripeFailedPayments,
		StripeAnomaly:        stripeAnomaly.IsAnomaly,
		StripeDelta:          stripeAnomaly.DeltaPct,
		CalendlyEvents:       string(calendlyEvents),
		GmailMessages:        string(gmailMessages),
	}

	// build prompt from template
	prompt, err := llm.BuildPrompt(digestData, "prompts/morning_digest.txt")
	if err != nil {
		fmt.Println("error building prompt:", err)
		return
	}

	// call the llm
	llmResp, err := s.llm.Generate(ctx, prompt)
	if err != nil {
		fmt.Println("error generating brief:", err)
		return
	}

	// save the digest
	_, err = s.db.Queries().SaveDigest(ctx, db.SaveDigestParams{
		Content:           llmResp,
		HasCriticalAlerts: stripeAnomaly.IsAnomaly,
	})
	if err != nil {
		fmt.Println("error saving digest:", err)
		return
	}
	fmt.Println("digest saved successfully")

	// send the digest
	s.webhook.Send(ctx, llmResp, stripeAnomaly.IsAnomaly)
	fmt.Println("webhook payload sent successfully")
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
