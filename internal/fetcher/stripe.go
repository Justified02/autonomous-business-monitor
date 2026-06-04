package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"time"

	"github.com/Justified02/abm/internal/storage"
	"github.com/Justified02/abm/internal/storage/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// StripeClient holds the HTTP client and API Key
type StripeClient struct {
	apiKey string
	db     *storage.Store
}

type stripeResponse struct {
	Data []stripeCharge `json:"data"`
}

type stripeCharge struct {
	ID 		string `json:"id"`
	Amount  int64 `json:"amount"`
	Currency string `json:"currency"`
	Status  string `json:"status"`
}

func (c *StripeClient) Parse(data []byte) (float64, int, error) {
	var response stripeResponse
	err := json.Unmarshal(data, &response)
	if err != nil{
		return 0, 0, fmt.Errorf("parsing stripe response: %w", err)
	}

	var revenue float64
	var failedCount int

	for _, charge := range response.Data {
		if charge.Status == "succeeded" {
			revenue += float64(charge.Amount)/100
		}

		if charge.Status == "failed" {
			failedCount++
		}
	}

	return revenue, failedCount, nil
}

func (c *StripeClient) Fetch(ctx context.Context) ([]byte, error) {

	// create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.stripe.com/v1/charges?limit=10", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Stripe's authentication requires only api key and no password. 
	// They designed it as basic auth where the password is always blank
	req.SetBasicAuth(c.apiKey, "")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func (c *StripeClient) Save(ctx context.Context, data []byte) error {
	_, err := c.db.Queries().SaveSnapshot(ctx, db.SaveSnapshotParams{
		Source: "stripe",
		Data: data,
		PeriodStart: pgtype.Timestamptz{},
		PeriodEnd: pgtype.Timestamptz{},
	})
	if err != nil {
		return fmt.Errorf("cannot save snapshot to the db: %w", err)
	}

	return nil
}


func NewStripeClient(apiKey string, db *storage.Store) *StripeClient {
	newClient := &StripeClient{
		apiKey: apiKey,
		db:     db,
	}

	return newClient
}
