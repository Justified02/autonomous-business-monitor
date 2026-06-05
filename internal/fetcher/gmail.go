package fetcher

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/Justified02/abm/internal/storage"
)

// GmailCLient struct
type GmailClient struct {
	clientID 		string
	clientSecret 	string
	refreshToken	string
	db 				*storage.Store
}

// Token response struct
type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// new GmailClient constructor
func NewGmailClient(clientID string, clientSecret string, refreshToken string, db *storage.Store) *GmailClient {
	newClient := &GmailClient{
		clientID: clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		db:			  db,
	}

	return newClient
}

// Method to get access token
func (c *GmailClient) getAccessToken(ctx context.Context) (string, error) {
	// Gmail's token endpoint expects form-encoded data
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("refresh_token", c.refreshToken)
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
}	