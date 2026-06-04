package llm

import (

)

type LLMClient struct {
	apiKey string
	model string
}

func NewLLMClient(apiKey string, model string) *LLMClient {
	newLLMClient := &LLMClient{
		apiKey: apiKey,
		model: model,
	}

	return newLLMClient
}