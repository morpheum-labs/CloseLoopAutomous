package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ChatClient calls an OpenAI-compatible POST {base}/chat/completions API.
type ChatClient struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model       string         `json:"model"`
	Messages    []chatMessage  `json:"messages"`
	Temperature float64        `json:"temperature,omitempty"`
	MaxTokens   int            `json:"max_tokens,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// ChatCompletion runs a single chat completion. ctx should carry a deadline.
func (c *ChatClient) ChatCompletion(ctx context.Context, model, systemPrompt, userPrompt string, temperature float64, maxTokens int) (string, error) {
	base := strings.TrimSuffix(strings.TrimSpace(c.BaseURL), "/")
	if base == "" {
		return "", fmt.Errorf("llm base url is empty")
	}
	if strings.TrimSpace(model) == "" {
		return "", fmt.Errorf("llm model is empty")
	}
	url := base + "/chat/completions"
	body := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(c.APIKey) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.APIKey))
	}
	cli := c.HTTP
	if cli == nil {
		cli = http.DefaultClient
	}
	res, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()
	respBytes, err := io.ReadAll(io.LimitReader(res.Body, 8<<20))
	if err != nil {
		return "", err
	}
	var parsed chatCompletionResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("llm response json: %w (http %d)", err, res.StatusCode)
	}
	if parsed.Error != nil && strings.TrimSpace(parsed.Error.Message) != "" {
		return "", fmt.Errorf("llm api error: %s", parsed.Error.Message)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("llm http %d: %s", res.StatusCode, strings.TrimSpace(string(respBytes)))
	}
	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llm response has no choices")
	}
	out := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if out == "" {
		return "", fmt.Errorf("llm returned empty content")
	}
	return out, nil
}

// DefaultHTTPClient returns an HTTP client suitable for LLM calls (no global timeout; use context).
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			MaxIdleConns:          8,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   15 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
