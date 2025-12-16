package gateway

import (
	"github.com/askaroe/dockify-backend/config"
	"github.com/askaroe/dockify-backend/internal/gateway/mindspore"
)

type Gateway struct {
	mindspore.MindSpore
}

func NewGateway(cfg *config.Config) *Gateway {
	return &Gateway{
		MindSpore: mindspore.NewMindSporeService(cfg),
	}
}
