package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Justified02/abm/internal/storage"
)

// GmailCLient struct
type GmailClient struct {
	GmailClientID     string
	GmailClientSecret string
	GmailRefreshToken string
	db 				  *storage.Store
}

// Token response struct
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// new GmailClient constructor
func NewGmailClient(clientID string, clientSecret string, refreshToken string, db *storage.Store) *GmailClient {
	newClient := &GmailClient{
		GmailClientID: clientID,
		GmailClientSecret: clientSecret,
		GmailRefreshToken: refreshToken,
		db:			  db,
	}

	return newClient
}

// Method to get access token
func (c *GmailClient) getAccessToken(ctx context.Context) (string, error) {
	// Gmail's token endpoint expects form-encoded data
	data := url.Values{}
	data.Set("client_id", c.GmailClientID)
	data.Set("client_secret", c.GmailClientSecret)
	data.Set("refresh_token", c.GmailRefreshToken)
	data.Set("grant_type", "refresh_token")

	body := strings.NewReader(data.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// set request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response tokenResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return response.AccessToken, nil
}

func (c *GmailClient) Fetch(ctx context.Context) ([]byte, error) {
	// get access token
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	// create request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://gmail.googleapis.com/gmail/v1/users/me/messages?maxResults=10", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// set auth header
	req.Header.Set("Authorization", "Bearer " + accessToken)

	// define http client
	client := &http.Client{}

	// make request
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