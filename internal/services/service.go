package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/repository"
	"net/http"
	"net/url"
	"strings"
)

type Service struct {
}

func NewService(
	repo *repository.Repository) *Service {
	return &Service{}
}

func (s *Service) GetAuthToken(ctx context.Context, code, state string) (entity.TokenResp, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {entity.ClientID},
		"client_secret": {entity.ClientSecret},
		"code":          {code},
		"redirect_uri":  {entity.RedirectURI},
	}

	resp, err := http.Post(
		"https://oauth-login.cloud.huawei.com/oauth2/v3/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return entity.TokenResp{}, fmt.Errorf("post token request: %w", err)
	}
	defer resp.Body.Close()

	var tokens entity.TokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return entity.TokenResp{}, fmt.Errorf("decode token response: %w", err)
	}

	// Save tokens per user (DB or cache)
	fmt.Println("AccessToken:", tokens.AccessToken)
	fmt.Println("RefreshToken:", tokens.RefreshToken)

	// Redirect or respond
	return tokens, nil
}
