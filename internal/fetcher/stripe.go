package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// StripeClient holds the HTTP client and API Key
type StripeClient struct {
	apiKey string
	db     string
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

func NewStripeClient(apiKey string, db string) *StripeClient {
	newClient := &StripeClient{
		apiKey: apiKey,
		db:     db,
	}

	return newClient
}
