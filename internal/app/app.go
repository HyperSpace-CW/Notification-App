package app

import (
	"github.com/HyperSpace-CW/Notification-App/config"
	"github.com/HyperSpace-CW/Notification-App/pkg/logger"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Logger.Level)

	log.Info("Starting the application")
}
