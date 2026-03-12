package main

import (
	"github.com/askaroe/dockify-backend/config"
	_ "github.com/askaroe/dockify-backend/docs"
	"github.com/askaroe/dockify-backend/internal/gateway"
	"github.com/askaroe/dockify-backend/internal/handlers"
	"github.com/askaroe/dockify-backend/internal/repository"
	"github.com/askaroe/dockify-backend/internal/router"
	"github.com/askaroe/dockify-backend/internal/scheduler"
	"github.com/askaroe/dockify-backend/internal/server"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/psql"
	"github.com/askaroe/dockify-backend/pkg/utils"
)

// @title Dockify Backend API
// @version 1.0
// @description API for Dockify backend.
// @schemes http https
func main() {
	logger := utils.NewLogger("dockify-backend")
	cfg, err := config.GetConfig()

	if err != nil {
		logger.Fatalf("failed to get config: %v", err)
	}

	db, err := psql.New(*cfg)
	if err != nil {
		logger.Fatalf("failed to initialize database: %v", err)
		return
	}

	repo := repository.NewRepository(db)

	gw := gateway.NewGateway(cfg)
	s := services.NewService(repo, gw)

	handler := handlers.NewHandler(logger, s)

	r := router.NewRouter(handler)

	// Start background scheduler: clears documents every 3 hours
	scheduler.StartDocumentCleanup(db, logger)

	srv := server.New(cfg, r, logger)
	srv.Start()
	srv.HandleShutdown()

}

