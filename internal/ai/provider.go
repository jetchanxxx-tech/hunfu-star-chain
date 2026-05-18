package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/huifu/star-chain/internal/config"
)

type ChatRequest struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`    // system / user / assistant
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Provider interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

type openAICompatProvider struct {
	endpoint   string
	apiKey     string
	model      string
	timeout    time.Duration
	maxRetries int
	client     *http.Client
}

func NewProvider(cfg config.AIConfig) Provider {
	pcfg, ok := cfg.Providers[cfg.DefaultProvider]
	if !ok || pcfg.APIKey == "" {
		return &noopProvider{}
	}
	dur, _ := time.ParseDuration(pcfg.Timeout)
	if dur == 0 {
		dur = 30 * time.Second
	}
	return &openAICompatProvider{
		endpoint:   pcfg.Endpoint,
		apiKey:     pcfg.APIKey,
		model:      pcfg.Model,
		timeout:    dur,
		maxRetries: pcfg.MaxRetries,
		client:     &http.Client{Timeout: dur},
	}
}

func (p *openAICompatProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body := map[string]any{
		"model":    p.model,
		"messages": req.Messages,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		resp, err := p.send(ctx, b)
		if err == nil {
			return resp, nil
		}
		if attempt == p.maxRetries {
			return nil, err
		}
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}
	return nil, fmt.Errorf("max retries exceeded")
}

func (p *openAICompatProvider) send(ctx context.Context, body []byte) (*ChatResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(b))
	}

	var cr ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &cr, nil
}

type noopProvider struct{}

func (n *noopProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	return &ChatResponse{
		Choices: []Choice{{Message: Message{Role: "assistant", Content: "AI service not configured"}}},
	}, nil
}
