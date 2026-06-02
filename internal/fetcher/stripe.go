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

func (c *StripeClient) Fetch(ctx context.Context) {

	// create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.stripe.com/v1/charges?limit=10", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(c.apiKey, "")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Println("Wrong status code")
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}

func NewStripeClient(apiKey string, db string) *StripeClient {
	newClient := &StripeClient{
		apiKey: apiKey,
		db:     db,
	}

	return newClient
}
