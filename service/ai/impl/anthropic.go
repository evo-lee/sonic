package aiimpl

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/service/ai"
)

type anthropicProvider struct {
	client anthropic.Client
	model  string
}

// NewAnthropicProvider returns an Anthropic-backed Provider when api_key is
// configured, or a no-op Provider that returns ErrNotConfigured otherwise.
// This keeps FX happy when AI is not set up.
func NewAnthropicProvider(cfg *config.Config) ai.Provider {
	aiCfg := cfg.AI
	if aiCfg.APIKey == "" {
		return noopProvider{}
	}
	model := aiCfg.Model
	if model == "" {
		model = anthropic.ModelClaudeHaiku4_5_20251001
	}
	client := anthropic.NewClient(option.WithAPIKey(aiCfg.APIKey))
	return &anthropicProvider{client: client, model: model}
}

func (p *anthropicProvider) Complete(ctx context.Context, req ai.CompletionRequest) (ai.CompletionResponse, error) {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 1024
	}

	params := anthropic.MessageNewParams{
		Model:     p.model,
		MaxTokens: maxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		},
	}
	if req.System != "" {
		params.System = []anthropic.TextBlockParam{{Text: req.System}}
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return ai.CompletionResponse{}, fmt.Errorf("anthropic complete: %w", err)
	}

	var text string
	for _, block := range msg.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return ai.CompletionResponse{
		Content:      text,
		InputTokens:  int(msg.Usage.InputTokens),
		OutputTokens: int(msg.Usage.OutputTokens),
	}, nil
}

type noopProvider struct{}

func (noopProvider) Complete(_ context.Context, _ ai.CompletionRequest) (ai.CompletionResponse, error) {
	return ai.CompletionResponse{}, fmt.Errorf("AI not configured: set ai.api_key in config")
}
