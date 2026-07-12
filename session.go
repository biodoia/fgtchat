package fgtchat

import (
	"sync"
	"time"
)

// Session tracks a conversation.
type Session struct {
	ID      string
	msgs    []Message
	Created time.Time
	mu      sync.RWMutex
}

// Message is a single message in a conversation.
type Message struct {
	Role      string // "user", "assistant", "system"
	Content   string
	Timestamp time.Time
}

// NewSession creates a new conversation session.
func NewSession(id string) *Session {
	return &Session{
		ID:      id,
		Created: time.Now(),
	}
}

// Append adds a message to the session.
func (s *Session) Append(role, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msgs = append(s.msgs, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
}

// Messages returns a copy of all messages.
func (s *Session) Messages() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Message, len(s.msgs))
	copy(out, s.msgs)
	return out
}

// LastN returns the last N messages.
func (s *Session) LastN(n int) []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if n > len(s.msgs) {
		n = len(s.msgs)
	}
	out := make([]Message, n)
	copy(out, s.msgs[len(s.msgs)-n:])
	return out
}
