package mindspore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/askaroe/dockify-backend/config"
	"github.com/askaroe/dockify-backend/pkg/utils"
)

type MindSpore interface {
	PredictLifestyle(ctx context.Context, body PredictLifestyleRequest) (PredictLifestyleResponse, error)
	PredictSleep(ctx context.Context, body PredictSleepRequest) (PredictSleepResponse, error)
}

type mindspore struct {
	cfg *config.Config
}

func NewMindSporeService(cfg *config.Config) MindSpore {
	return &mindspore{cfg: cfg}
}

func (m *mindspore) PredictLifestyle(ctx context.Context, body PredictLifestyleRequest) (PredictLifestyleResponse, error) {
	headers := utils.Headers{{"Content-Type", "application/json"}}
	b := new(bytes.Buffer)

	if err := json.NewEncoder(b).Encode(body); err != nil {
		return PredictLifestyleResponse{}, fmt.Errorf("could not encode request for jira stats: %w", err)
	}

	req := utils.Request{
		Method:  http.MethodPost,
		URL:     m.cfg.MindsporeModelURL + PredictLifestyleEndpoint,
		Headers: headers,
		Body:    b,
	}

	response, err := utils.SendRequest(ctx, &req)
	if err != nil {
		return PredictLifestyleResponse{}, fmt.Errorf("could not send request to mindspore model: %w", err)
	}

	var predictResponse PredictLifestyleResponse
	if err := json.Unmarshal(response, &predictResponse); err != nil {
		return PredictLifestyleResponse{}, fmt.Errorf("could not unmarshal response from mindspore model: %w", err)
	}

	return predictResponse, nil
}

func (m *mindspore) PredictSleep(ctx context.Context, body PredictSleepRequest) (PredictSleepResponse, error) {
	headers := utils.Headers{{"Content-Type", "application/json"}}
	b := new(bytes.Buffer)

	if err := json.NewEncoder(b).Encode(body); err != nil {
		return PredictSleepResponse{}, fmt.Errorf("could not encode request for jira stats: %w", err)
	}

	req := utils.Request{
		Method:  http.MethodPost,
		URL:     m.cfg.MindsporeModelURL + PredictSleepEndpoint,
		Headers: headers,
		Body:    b,
	}

	response, err := utils.SendRequest(ctx, &req)
	if err != nil {
		return PredictSleepResponse{}, fmt.Errorf("could not send request to mindspore model: %w", err)
	}

	var predictResponse PredictSleepResponse
	if err := json.Unmarshal(response, &predictResponse); err != nil {
		return PredictSleepResponse{}, fmt.Errorf("could not unmarshal response from mindspore model: %w", err)
	}

	return predictResponse, nil
}
