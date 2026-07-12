package fgtchat

import "context"

// Tool is a function the agent can call during conversation.
type Tool struct {
	Name        string
	Description string
	Handler     func(ctx context.Context, args map[string]interface{}) (string, error)
}
