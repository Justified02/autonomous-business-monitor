package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LLMClient struct {
	apiKey string
	model string
}

// Request Structs
type generateRequest struct {
	Contents []content `json:"contents"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

// Response Structs
type generateResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content content `json:"content"`
}


func NewLLMClient(apiKey string, model string) *LLMClient {
	newLLMClient := &LLMClient{
		apiKey: apiKey,
		model: model,
	}

	return newLLMClient
}

func (c *LLMClient) Generate(ctx context.Context, prompt string) (string, error) {
	// Build the request body
	reqBody := generateRequest{
		Contents: []content{
			{Parts: []part{{Text: prompt}}},
		},
	}

	// Marshal the struct to JSON
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}

	// build the URL with the API key as a query parameter
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)

	// Create the POST request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// execute request and handle response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// decode the response and extract the text
	var geminiResp generateResponse
	err = json.NewDecoder(resp.Body).Decode(&geminiResp)
	if err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}