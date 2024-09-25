package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/auth_service/configs"
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/internal/service"
	"github.com/magmaheat/auth_service/pkg/hasher"
	"github.com/magmaheat/auth_service/pkg/httpserver"
	"github.com/magmaheat/auth_service/pkg/postgres"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(configPath string) {
	cfg, err := configs.New(configPath)
	if err != nil {
		log.Fatalf("error setup config: %w", err)
	}

	setupLogger(cfg.Log.Level)

	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.URL, postgres.MaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()

	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	log.Info("Initializing service...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewBCRYTHasher(),
		SignKey:  cfg.SignKey,
		TokenTTL: cfg.TokenTTL,
	}

	services := service.NewServices(deps)

	log.Info("Initializing handlers and routes...")
	handler := echo.New()

	//TODO init routes

	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Port))

	log.Info("Configuring grace shutdown...")
	interrupt := make(chan os.Signal, 1)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Info(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Info(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
