package scheduler

import (
	"context"
	"fmt"

	"github.com/Justified02/abm/internal/fetcher"
	"github.com/robfig/cron/v3"
	"time"
)

type fetchResult struct {
	data []byte
	err  error
}

type Scheduler struct {
	stripe *fetcher.StripeClient
}

func NewScheduler(s *fetcher.StripeClient) *Scheduler {
	newScheduler := &Scheduler{
		stripe: s,
	}

	return newScheduler
}

// every morning, the scheduler needs to collect data from all sources (Stripe, Gmail, Calendly)
// fetchAllSources means "run the data collection process"
func (s *Scheduler) fetchAllSources(ctx context.Context) {
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

	fmt.Println(string(result.data))
}

// Start the cron job - call the fetchAllSources function in the process
func (s *Scheduler) Start(cronSchedule string) {
	// create a cron scheduler
	c := cron.New()

	// add a cron job
	_, err := c.AddFunc(cronSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
		defer cancel()

		s.fetchAllSources(ctx)
	})
	if err != nil {
		fmt.Println("adding cron job:", err)
    	return 
	}

	c.Start()
	select {}
}