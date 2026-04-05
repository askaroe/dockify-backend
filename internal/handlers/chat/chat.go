package chat

import (
	"io"
	"net/http"
	"strconv"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Chat interface {
	SendMessage(c *gin.Context)
	SendMessageStream(c *gin.Context)
	GetHistory(c *gin.Context)
	ClearHistory(c *gin.Context)
}

type chatHandler struct {
	s      *services.Service
	logger *utils.Logger
}

func NewChatHandler(s *services.Service, logger *utils.Logger) Chat {
	return &chatHandler{s: s, logger: logger}
}

// SendMessage godoc
// @Summary Send a chat message
// @Description Send a message and get a response from the AI assistant
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body entity.ChatRequest true "Chat message"
// @Success 200 {object} entity.ChatResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/chat [post]
func (h *chatHandler) SendMessage(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	reply, err := h.s.Chat.SendMessage(ctx, req.UserID, req.DocID, req.Message)
	if err != nil {
		h.logger.Errorf("Chat SendMessage error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to get response"})
		return
	}

	c.JSON(http.StatusOK, entity.ChatResponse{Reply: reply})
}

// SendMessageStream godoc
// @Summary Send a chat message (streaming)
// @Description Send a message and receive a streaming SSE response from the AI assistant
// @Tags Chat
// @Accept json
// @Produce text/event-stream
// @Param request body entity.ChatRequest true "Chat message"
// @Success 200 {string} string "SSE stream"
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/chat/stream [post]
func (h *chatHandler) SendMessageStream(c *gin.Context) {
	ctx := c.Request.Context()

	var req entity.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid request"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	chunks := make(chan string, 64)

	go func() {
		if err := h.s.Chat.SendMessageStream(ctx, req.UserID, req.DocID, req.Message, chunks); err != nil {
			h.logger.Errorf("Chat stream error: %v", err)
		}
	}()

	c.Stream(func(w io.Writer) bool {
		chunk, ok := <-chunks
		if !ok {
			return false
		}
		c.SSEvent("message", chunk)
		return true
	})
}

// GetHistory godoc
// @Summary Get chat history
// @Description Returns chat message history for a user, optionally filtered by document
// @Tags Chat
// @Produce json
// @Param user_id query integer true "User ID"
// @Param doc_id query string false "Document ID"
// @Success 200 {array} entity.ChatMessageResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/chat [get]
func (h *chatHandler) GetHistory(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	var docID *string
	if d := c.Query("doc_id"); d != "" {
		docID = &d
	}

	messages, err := h.s.Chat.GetHistory(ctx, userID, docID)
	if err != nil {
		h.logger.Errorf("Chat GetHistory error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to get history"})
		return
	}

	var resp []entity.ChatMessageResponse
	for _, msg := range messages {
		createdAt := ""
		if msg.CreatedAt != nil {
			createdAt = msg.CreatedAt.Format("2006-01-02 15:04:05")
		}
		resp = append(resp, entity.ChatMessageResponse{
			ID:        msg.ID,
			UserID:    msg.UserID,
			DocID:     msg.DocID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: createdAt,
		})
	}

	if resp == nil {
		resp = []entity.ChatMessageResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// ClearHistory godoc
// @Summary Clear chat history
// @Description Delete chat history for a user, optionally filtered by document
// @Tags Chat
// @Produce json
// @Param user_id query integer true "User ID"
// @Param doc_id query string false "Document ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/chat [delete]
func (h *chatHandler) ClearHistory(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	var docID *string
	if d := c.Query("doc_id"); d != "" {
		docID = &d
	}

	if err := h.s.Chat.ClearHistory(ctx, userID, docID); err != nil {
		h.logger.Errorf("Chat ClearHistory error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to clear history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "history cleared"})
}
