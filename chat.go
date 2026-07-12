// Package fgtchat provides a conversational agent interface for fgt-sdk apps.
// Import this library and register your app's tools to get a chat UI.
package fgtchat

import (
	"log/slog"
	"sync"
)

// Chat is the conversational agent engine.
type Chat struct {
	llm      LLMBackend
	tools    map[string]Tool
	sessions map[string]*Session
	mu       sync.RWMutex
	logger   *slog.Logger
}

// Option configures the Chat engine.
type Option func(*Chat)

// WithLLM sets the LLM backend.
func WithLLM(llm LLMBackend) Option {
	return func(c *Chat) { c.llm = llm }
}

// WithLogger sets the logger.
func WithLogger(l *slog.Logger) Option {
	return func(c *Chat) { c.logger = l }
}

// New creates a new Chat engine.
func New(opts ...Option) *Chat {
	c := &Chat{
		tools:    make(map[string]Tool),
		sessions: make(map[string]*Session),
		logger:   slog.Default(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.llm == nil {
		c.llm = NewLLMFromEnv()
	}
	return c
}

// RegisterTool adds a tool the agent can call.
func (c *Chat) RegisterTool(t Tool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tools[t.Name] = t
	c.logger.Info("registered chat tool", "name", t.Name)
}

// Tools returns all registered tools.
func (c *Chat) Tools() []Tool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Tool, 0, len(c.tools))
	for _, t := range c.tools {
		out = append(out, t)
	}
	return out
}
