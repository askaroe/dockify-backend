package server

import (
	"context"
	"errors"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/askaroe/dockify-backend/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	logger     *utils.Logger
}

func New(cfg *config.Config, router *gin.Engine, logger *utils.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: router,
		},
		logger: logger,
	}
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("starting server")
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("error starting server: %v", err)
		}
	}()
}

func (s *Server) HandleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatalf("error shutting down server: %v", err)
	}

	s.logger.Info("server exiting")
}
