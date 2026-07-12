package fgtchat

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ChatRequest is the incoming message from the user.
type ChatRequest struct {
	Message string `json:"message"`
	Session string `json:"session,omitempty"`
}

// ChatResponse is the outgoing message to the user.
type ChatResponse struct {
	Message string   `json:"message"`
	Session string   `json:"session"`
	Tools   []string `json:"tools_used,omitempty"`
}

// HandleEcho registers chat routes on an Echo instance.
func (c *Chat) HandleEcho(e *echo.Echo) {
	e.POST("/api/chat", c.handleMessage)
	e.GET("/api/chat/stream", c.handleStream)
	e.GET("/api/chat/tools", c.handleListTools)
	e.GET("/chat", c.handleChatPage)
}

func (c *Chat) handleMessage(ctx echo.Context) error {
	var req ChatRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	sessionID := req.Session
	if sessionID == "" {
		sessionID = generateID()
	}

	session := c.getOrCreateSession(sessionID)
	session.Append("user", req.Message)

	response, toolsUsed, err := c.llm.Chat(ctx.Request().Context(), session, c.Tools())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	session.Append("assistant", response)

	return ctx.JSON(http.StatusOK, ChatResponse{
		Message: response,
		Session: sessionID,
		Tools:   toolsUsed,
	})
}

func (c *Chat) handleStream(ctx echo.Context) error {
	message := ctx.QueryParam("message")
	sessionID := ctx.QueryParam("session")

	if message == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "message required"})
	}

	if sessionID == "" {
		sessionID = generateID()
	}

	session := c.getOrCreateSession(sessionID)
	session.Append("user", message)

	ctx.Response().Header().Set("Content-Type", "text/event-stream")
	ctx.Response().Header().Set("Cache-Control", "no-cache")
	ctx.Response().Header().Set("Connection", "keep-alive")

	flusher, ok := ctx.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
	}

	stream, err := c.llm.Stream(ctx.Request().Context(), session, c.Tools())
	if err != nil {
		return err
	}

	var fullResponse strings.Builder
	for chunk := range stream {
		fullResponse.WriteString(chunk)
		fmt.Fprintf(ctx.Response().Writer, "data: %s\n\n", chunk)
		flusher.Flush()
	}

	session.Append("assistant", fullResponse.String())
	fmt.Fprintf(ctx.Response().Writer, "data: [DONE]\n\n")
	flusher.Flush()

	return nil
}

func (c *Chat) handleListTools(ctx echo.Context) error {
	tools := c.Tools()
	type toolInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	out := make([]toolInfo, len(tools))
	for i, t := range tools {
		out[i] = toolInfo{Name: t.Name, Description: t.Description}
	}
	return ctx.JSON(http.StatusOK, out)
}

func (c *Chat) handleChatPage(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, chatHTML)
}

func (c *Chat) getOrCreateSession(id string) *Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s, ok := c.sessions[id]; ok {
		return s
	}
	s := NewSession(id)
	c.sessions[id] = s
	return s
}

func generateID() string {
	return fmt.Sprintf("chat-%d", time.Now().UnixNano())
}

const chatHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Chat Agent</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Courier New', monospace; background: #0a0a0a; color: #0f0; height: 100vh; display: flex; flex-direction: column; }
        #header { padding: 10px 20px; border-bottom: 1px solid #0f0; background: #111; }
        #chat { flex: 1; overflow-y: auto; padding: 20px; }
        .msg { margin: 8px 0; line-height: 1.4; }
        .user { color: #0ff; }
        .user::before { content: '> '; color: #888; }
        .assistant { color: #0f0; }
        .system { color: #888; font-style: italic; }
        #footer { padding: 10px 20px; border-top: 1px solid #0f0; background: #111; display: flex; gap: 10px; }
        #input { flex: 1; background: #000; color: #0f0; border: 1px solid #333; padding: 10px; font-family: inherit; font-size: 14px; }
        #input:focus { outline: none; border-color: #0f0; }
        #send { background: #0f0; color: #000; border: none; padding: 10px 20px; cursor: pointer; font-family: inherit; font-weight: bold; }
        #send:hover { background: #0a0; }
        #status { color: #888; font-size: 12px; padding: 5px 20px; }
    </style>
</head>
<body>
    <div id="header">
        <h2>Conversational Agent</h2>
    </div>
    <div id="chat"></div>
    <div id="status">Ready</div>
    <div id="footer">
        <input type="text" id="input" placeholder="Type a message..." autofocus />
        <button id="send" onclick="send()">Send</button>
    </div>

    <script>
        let sessionId = '';
        const chat = document.getElementById('chat');
        const input = document.getElementById('input');
        const status = document.getElementById('status');

        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                send();
            }
        });

        async function send() {
            const msg = input.value.trim();
            if (!msg) return;

            appendMessage('user', msg);
            input.value = '';
            status.textContent = 'Thinking...';

            try {
                const res = await fetch('/api/chat', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({message: msg, session: sessionId})
                });

                const data = await res.json();
                if (data.error) {
                    appendMessage('system', 'Error: ' + data.error);
                } else {
                    sessionId = data.session;
                    appendMessage('assistant', data.message);
                    if (data.tools_used && data.tools_used.length > 0) {
                        appendMessage('system', 'Tools used: ' + data.tools_used.join(', '));
                    }
                }
            } catch (err) {
                appendMessage('system', 'Connection error: ' + err.message);
            }

            status.textContent = 'Ready';
        }

        function appendMessage(role, text) {
            const div = document.createElement('div');
            div.className = 'msg ' + role;
            div.textContent = text;
            chat.appendChild(div);
            chat.scrollTop = chat.scrollHeight;
        }
    </script>
</body>
</html>`
