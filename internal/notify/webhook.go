package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookClient struct {
	url 	string
}

type webhookPayload struct {
	Content 		string `json:"content"`
	HasCriticalAlerts bool `json:"has_critical_alerts"`
}

func NewWebhookClient(url string) *WebhookClient {
	newClient := &WebhookClient{
		url: url,
	}

	return newClient
}

func (c *WebhookClient) Send(ctx context.Context, content string, hasCriticalAlerts bool) error {
	// build the request body
	payload := webhookPayload{
		Content: 		content,
		HasCriticalAlerts: hasCriticalAlerts,
	}

	// Marshal the struct to JSON
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling request: %w", err)
	}

	// create the request
	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}