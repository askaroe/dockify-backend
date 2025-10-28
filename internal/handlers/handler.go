package handlers

import (
	"fmt"
	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s      *services.Service
	logger *utils.Logger
}

func NewHandler(logger *utils.Logger, s *services.Service) *Handler {
	return &Handler{
		s:      s,
		logger: logger,
	}
}

// HealthCheck godoc
// @Summary Health Check (Live)
// @Description Returns the live status of the service
// @Tags Health
// @Produce json
// @Success 200 {string} string "health"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "health")
}

func (h *Handler) LoginHuawei(c *gin.Context) {
	scopes := []string{
		"openid",
		"https://www.huawei.com/healthkit/step.read",
		"https://www.huawei.com/healthkit/heart_rate.read",
	}

	authUrl := fmt.Sprintf(
		"https://oauth-login.cloud.huawei.com/oauth2/v3/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s",
		url.QueryEscape(entity.ClientID),
		url.QueryEscape(entity.RedirectURI),
		url.QueryEscape(strings.Join(scopes, " ")),
		entity.State,
	)

	c.Redirect(http.StatusFound, authUrl)
}

func (h *Handler) HuaweiCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if state != entity.State {
		h.logger.Info(state)
		h.logger.Info(entity.State)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	tokens, err := h.s.GetAuthToken(c.Request.Context(), code, state)
	if err != nil {
		h.logger.Errorf("failed to get auth token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auth token " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
