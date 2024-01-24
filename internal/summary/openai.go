package summary

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"sync"
	"time"
)

type OpenAISummariser struct {
	client  *openai.Client
	prompt  string
	model   string
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummariser(apiKey, model, prompt string) *OpenAISummariser {
	s := &OpenAISummariser{
		client: openai.NewClient(apiKey),
		prompt: prompt,
		model:  model,
	}

	log.Printf("openai summariser is enabled: %v", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAISummariser) Summarise(text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return text, nil
	}

	request := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: s.prompt},
			{Role: openai.ChatMessageRoleUser, Content: text},
		},
		MaxTokens:   1024,
		Temperature: 1,
		TopP:        1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in openai response")
	}

	summary := resp.Choices[0].Message.Content
	summary = strings.TrimSpace(summary)
	if !strings.HasSuffix(summary, ".") {
		summary += "."
	}

	return summary, nil
}
