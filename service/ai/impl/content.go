package aiimpl

import (
	"context"
	"strings"

	"github.com/go-sonic/sonic/service/ai"
)

type contentServiceImpl struct {
	provider ai.Provider
}

func NewContentService(provider ai.Provider) ai.ContentService {
	return &contentServiceImpl{provider: provider}
}

func (s *contentServiceImpl) Summarize(ctx context.Context, content string) (string, error) {
	resp, err := s.provider.Complete(ctx, ai.CompletionRequest{
		System: "You are a concise technical writer. Reply with a single paragraph summary, no more than 3 sentences, in the same language as the input.",
		Prompt: "Summarize the following article:\n\n" + content,
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Content), nil
}

func (s *contentServiceImpl) SuggestTags(ctx context.Context, title, content string) ([]string, error) {
	resp, err := s.provider.Complete(ctx, ai.CompletionRequest{
		System: "You are a content tagger. Reply with a comma-separated list of 3–6 lowercase tags, no explanation.",
		Prompt: "Title: " + title + "\n\nContent:\n" + content,
	})
	if err != nil {
		return nil, err
	}
	raw := strings.Split(resp.Content, ",")
	tags := make([]string, 0, len(raw))
	for _, t := range raw {
		if tag := strings.TrimSpace(t); tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func (s *contentServiceImpl) Polish(ctx context.Context, content string) (string, error) {
	resp, err := s.provider.Complete(ctx, ai.CompletionRequest{
		System: "You are an expert editor. Improve the clarity, flow, and readability of the text. Preserve the original language, meaning, and markdown formatting.",
		Prompt: content,
		MaxTokens: 4096,
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Content), nil
}
