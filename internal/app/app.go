package app

import (
	"github.com/labstack/echo/v4"
	"github.com/magmaheat/auth_service/configs"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	cfg, err := configs.New(configPath)
	if err != nil {
		log.Fatalf("error setup config: %w", err)
	}

	setupLogger(cfg.Log.Level)

	handler := echo.New()
	if err = handler.Start(":" + cfg.HTTP.Port); err != nil {
		log.Fatalf("error start server: %w", err)
	}
}
