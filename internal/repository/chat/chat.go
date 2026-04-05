package chat

import (
	"context"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Chat interface {
	Save(ctx context.Context, msg models.ChatMessage) error
	GetHistory(ctx context.Context, userID int, docID *string, limit int) ([]models.ChatMessage, error)
	DeleteHistory(ctx context.Context, userID int, docID *string) error
}

type chat struct {
	db *psql.Client
}

func NewChatRepository(db *psql.Client) Chat {
	return &chat{db: db}
}

func (c *chat) Save(ctx context.Context, msg models.ChatMessage) error {
	query := `INSERT INTO chat_messages (user_id, doc_id, role, content) VALUES ($1, $2, $3, $4)`
	_, err := c.db.Exec(ctx, query, msg.UserID, msg.DocID, msg.Role, msg.Content)
	return err
}

func (c *chat) GetHistory(ctx context.Context, userID int, docID *string, limit int) ([]models.ChatMessage, error) {
	var (
		messages []models.ChatMessage
	)

	// Query with DESC to get most recent, then reverse for chronological order
	var query string
	var args []interface{}
	if docID != nil {
		query = `SELECT id, user_id, doc_id, role, content, created_at
		         FROM chat_messages WHERE user_id = $1 AND doc_id = $2
		         ORDER BY created_at DESC LIMIT $3`
		args = []interface{}{userID, *docID, limit}
	} else {
		query = `SELECT id, user_id, doc_id, role, content, created_at
		         FROM chat_messages WHERE user_id = $1 AND doc_id IS NULL
		         ORDER BY created_at DESC LIMIT $2`
		args = []interface{}{userID, limit}
	}

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.DocID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse so oldest comes first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (c *chat) DeleteHistory(ctx context.Context, userID int, docID *string) error {
	if docID != nil {
		_, err := c.db.Exec(ctx, `DELETE FROM chat_messages WHERE user_id = $1 AND doc_id = $2`, userID, *docID)
		return err
	}
	_, err := c.db.Exec(ctx, `DELETE FROM chat_messages WHERE user_id = $1 AND doc_id IS NULL`, userID)
	return err
}
