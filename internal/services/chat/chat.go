package chat

import (
	"context"
	"fmt"

	"github.com/askaroe/dockify-backend/internal/gateway"
	"github.com/askaroe/dockify-backend/internal/gateway/deepseek"
	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/internal/repository"
)

type Chat interface {
	SendMessage(ctx context.Context, userID int, docID *string, userMessage string) (string, error)
	SendMessageStream(ctx context.Context, userID int, docID *string, userMessage string, out chan<- string) error
	GetHistory(ctx context.Context, userID int, docID *string) ([]models.ChatMessage, error)
	ClearHistory(ctx context.Context, userID int, docID *string) error
}

type chat struct {
	repo *repository.Repository
	gw   *gateway.Gateway
}

func NewChatService(repo *repository.Repository, gw *gateway.Gateway) Chat {
	return &chat{repo: repo, gw: gw}
}

func (c *chat) buildMessages(ctx context.Context, userID int, docID *string, userMessage string) ([]deepseek.ChatMessage, error) {
	systemPrompt := "You are a medical assistant helping users understand their medical documents and health data. Be concise and clear."

	// If chatting about a specific document, include its summary as context
	if docID != nil {
		doc, err := c.repo.Document.GetByID(ctx, *docID)
		if err == nil && doc.Summary != "" {
			systemPrompt += "\n\nDocument summary:\n" + doc.Summary
		}
	}

	messages := []deepseek.ChatMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Load recent chat history for context
	history, _ := c.repo.Chat.GetHistory(ctx, userID, docID, 20)
	for _, h := range history {
		messages = append(messages, deepseek.ChatMessage{Role: h.Role, Content: h.Content})
	}

	messages = append(messages, deepseek.ChatMessage{Role: "user", Content: userMessage})
	return messages, nil
}

func (c *chat) SendMessage(ctx context.Context, userID int, docID *string, userMessage string) (string, error) {
	messages, err := c.buildMessages(ctx, userID, docID, userMessage)
	if err != nil {
		return "", fmt.Errorf("build messages: %w", err)
	}

	reply, err := c.gw.DeepSeek.Chat(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("deepseek chat: %w", err)
	}

	// Persist both turns
	c.repo.Chat.Save(ctx, models.ChatMessage{UserID: userID, DocID: docID, Role: "user", Content: userMessage})
	c.repo.Chat.Save(ctx, models.ChatMessage{UserID: userID, DocID: docID, Role: "assistant", Content: reply})

	return reply, nil
}

func (c *chat) SendMessageStream(ctx context.Context, userID int, docID *string, userMessage string, out chan<- string) error {
	messages, err := c.buildMessages(ctx, userID, docID, userMessage)
	if err != nil {
		close(out)
		return fmt.Errorf("build messages: %w", err)
	}

	// Collect full reply while streaming
	collector := make(chan string, 64)
	go func() {
		var fullReply string
		for chunk := range collector {
			fullReply += chunk
			out <- chunk
		}
		close(out)

		// Persist both turns after stream completes
		c.repo.Chat.Save(ctx, models.ChatMessage{UserID: userID, DocID: docID, Role: "user", Content: userMessage})
		c.repo.Chat.Save(ctx, models.ChatMessage{UserID: userID, DocID: docID, Role: "assistant", Content: fullReply})
	}()

	return c.gw.DeepSeek.ChatStream(ctx, messages, collector)
}

func (c *chat) GetHistory(ctx context.Context, userID int, docID *string) ([]models.ChatMessage, error) {
	return c.repo.Chat.GetHistory(ctx, userID, docID, 100)
}

func (c *chat) ClearHistory(ctx context.Context, userID int, docID *string) error {
	return c.repo.Chat.DeleteHistory(ctx, userID, docID)
}
