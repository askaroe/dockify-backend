package deepseek

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/askaroe/dockify-backend/config"
)

type DeepSeek interface {
	AnalyzeDocument(ctx context.Context, text string) (string, error)
	Chat(ctx context.Context, messages []ChatMessage) (string, error)
	ChatStream(ctx context.Context, messages []ChatMessage, out chan<- string) error
}

type deepseek struct {
	baseURL string
	apiKey  string
	model   string
}

func NewDeepSeekService(cfg *config.Config) DeepSeek {
	return &deepseek{
		baseURL: cfg.DeepseekBaseURL,
		apiKey:  cfg.DeepseekAPIKey,
		model:   "deepseek-chat",
	}
}

func (d *deepseek) AnalyzeDocument(ctx context.Context, text string) (string, error) {
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: "You are a medical assistant. Analyze the following medical document and provide a brief overview of the key findings, results, and any notable values. Be concise and clear.",
		},
		{
			Role:    "user",
			Content: text,
		},
	}

	return d.Chat(ctx, messages)
}

func (d *deepseek) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	reqBody := ChatRequest{
		Model:    d.model,
		Messages: messages,
		Stream:   false,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.baseURL+ChatCompletionsEndpoint, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body []byte
		body, _ = readAll(resp.Body)
		return "", fmt.Errorf("deepseek returned %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from deepseek")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (d *deepseek) ChatStream(ctx context.Context, messages []ChatMessage, out chan<- string) error {
	defer close(out)

	reqBody := ChatRequest{
		Model:    d.model,
		Messages: messages,
		Stream:   true,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.baseURL+ChatCompletionsEndpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body []byte
		body, _ = readAll(resp.Body)
		return fmt.Errorf("deepseek returned %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			out <- chunk.Choices[0].Delta.Content
		}
	}

	return scanner.Err()
}

func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	return buf.Bytes(), err
}
