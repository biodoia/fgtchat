package fgtchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LLMBackend is the interface for LLM providers.
type LLMBackend interface {
	// Chat sends a conversation and returns the response.
	Chat(ctx context.Context, session *Session, tools []Tool) (response string, toolsUsed []string, err error)
	// Stream sends a conversation and returns a channel of chunks.
	Stream(ctx context.Context, session *Session, tools []Tool) (<-chan string, error)
}

// NewLLMFromEnv creates an LLM backend from environment variables.
func NewLLMFromEnv() LLMBackend {
	provider := os.Getenv("FGTCHAT_LLM_PROVIDER")
	if provider == "" {
		provider = "9router"
	}

	switch provider {
	case "9router":
		return NewNineRouterLLM(os.Getenv("NINEROUTER_URL"))
	case "ollama":
		return NewOllamaLLM(os.Getenv("FGTCHAT_LLM_URL"))
	case "openai":
		return NewOpenAILLM(os.Getenv("OPENAI_API_KEY"))
	default:
		return NewNineRouterLLM("")
	}
}

// NineRouterLLM uses 9Router as the LLM backend.
type NineRouterLLM struct {
	baseURL string
	client  *http.Client
}

// NewNineRouterLLM creates a 9Router LLM backend.
func NewNineRouterLLM(baseURL string) *NineRouterLLM {
	if baseURL == "" {
		baseURL = "http://localhost:9090"
	}
	return &NineRouterLLM{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (l *NineRouterLLM) Chat(ctx context.Context, session *Session, tools []Tool) (string, []string, error) {
	msgs := session.Messages()
	reqBody := map[string]interface{}{
		"model":    "gpt-4o-mini",
		"messages": formatMessages(msgs),
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("9router error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", nil, fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil, nil
}

func (l *NineRouterLLM) Stream(ctx context.Context, session *Session, tools []Tool) (<-chan string, error) {
	msgs := session.Messages()
	reqBody := map[string]interface{}{
		"model":    "gpt-4o-mini",
		"messages": formatMessages(msgs),
		"stream":   true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("9router error %d: %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan string, 64)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := decoder.Decode(&chunk); err != nil {
				break
			}
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				ch <- chunk.Choices[0].Delta.Content
			}
		}
	}()

	return ch, nil
}

// OllamaLLM uses Ollama as the LLM backend.
type OllamaLLM struct {
	baseURL string
	client  *http.Client
}

func NewOllamaLLM(baseURL string) *OllamaLLM {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaLLM{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (l *OllamaLLM) Chat(ctx context.Context, session *Session, tools []Tool) (string, []string, error) {
	msgs := session.Messages()
	reqBody := map[string]interface{}{
		"model":    os.Getenv("FGTCHAT_LLM_MODEL"),
		"messages": formatMessages(msgs),
		"stream":   false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Message.Content, nil, nil
}

func (l *OllamaLLM) Stream(ctx context.Context, session *Session, tools []Tool) (<-chan string, error) {
	msgs := session.Messages()
	reqBody := map[string]interface{}{
		"model":    os.Getenv("FGTCHAT_LLM_MODEL"),
		"messages": formatMessages(msgs),
		"stream":   true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	ch := make(chan string, 64)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}
			if err := decoder.Decode(&chunk); err != nil {
				break
			}
			if chunk.Message.Content != "" {
				ch <- chunk.Message.Content
			}
		}
	}()

	return ch, nil
}

// OpenAILLM uses OpenAI API directly.
type OpenAILLM struct {
	apiKey string
	client *http.Client
}

func NewOpenAILLM(apiKey string) *OpenAILLM {
	return &OpenAILLM{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (l *OpenAILLM) Chat(ctx context.Context, session *Session, tools []Tool) (string, []string, error) {
	return NewNineRouterLLM("https://api.openai.com").Chat(ctx, session, tools)
}

func (l *OpenAILLM) Stream(ctx context.Context, session *Session, tools []Tool) (<-chan string, error) {
	return NewNineRouterLLM("https://api.openai.com").Stream(ctx, session, tools)
}

// formatMessages converts session messages to OpenAI format.
func formatMessages(msgs []Message) []map[string]string {
	out := make([]map[string]string, len(msgs))
	for i, m := range msgs {
		out[i] = map[string]string{
			"role":    m.Role,
			"content": m.Content,
		}
	}
	return out
}
