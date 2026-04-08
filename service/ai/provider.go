package ai

import "context"

// CompletionRequest is a request to generate text from a prompt.
type CompletionRequest struct {
	System string
	Prompt string
	// MaxTokens defaults to 1024 if zero.
	MaxTokens int
}

// CompletionResponse holds the generated text and token usage.
type CompletionResponse struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

// Provider is the abstraction over any LLM backend.
type Provider interface {
	Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}
