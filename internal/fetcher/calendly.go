package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Justified02/abm/internal/storage"
	"github.com/Justified02/abm/internal/storage/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type CalendlyClient struct {
	apiKey string
	userURI string
	db     *storage.Store
}

// calendly client constructor
func NewCalendlyClient(apiKey string, userURI string, db *storage.Store) *CalendlyClient {
	newClient := &CalendlyClient{
		apiKey: apiKey,
		userURI: userURI,
		db:     db,
	}

	return newClient
}

func (c *CalendlyClient) Fetch(ctx context.Context) ([]byte, error) {
	// create new request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.calendly.com/scheduled_events?user=" + c.userURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// set authentication
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	// chec status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	return body, nil
}

func (c *CalendlyClient) Save(ctx context.Context, data []byte) error {
	_, err := c.db.Queries().SaveSnapshot(ctx, db.SaveSnapshotParams{
		Source:      "calendly",
		Data:        data,
		PeriodStart: pgtype.Timestamptz{},
		PeriodEnd:   pgtype.Timestamptz{},
	})
	if err != nil {
		return fmt.Errorf("cannot save snapshot to the db: %w", err)
	}

	return nil
}
