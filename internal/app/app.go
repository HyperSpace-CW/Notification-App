package app

import (
	"fmt"
	"github.com/HyperSpace-CW/Notification-App/config"
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	"github.com/HyperSpace-CW/Notification-App/internal/transport/http"
	"github.com/HyperSpace-CW/Notification-App/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Logger.Level)

	log.Info("Starting the application")

	notificationService := services.NewNotificationService(cfg)

	httpServer := http.NewServer(http.Config{
		Addr:                cfg.Server.Addr,
		Logger:              log,
		NotificationService: notificationService,
	})

	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatal(fmt.Sprintf("error occurred while running HTTP server: %v", err))
		}
	}()
	log.Info("HTTP Server started", zap.String("port", cfg.Server.Addr))

	log.Info("Application successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	<-quit

	log.Info("shutdown HTTP server...")
	if err := httpServer.Shutdown(); err != nil {
		log.Error(fmt.Sprintf("failed to shutdown HTTP server: %v", err))
	} else {
		log.Info("HTTP server successfully shutdown")
	}

	log.Info("Application successfully exited")
}
