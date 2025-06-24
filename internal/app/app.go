package app

import (
	"fmt"
	"github.com/HyperSpace-CW/Notification-App/config"
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	"github.com/HyperSpace-CW/Notification-App/internal/transport/http"
	"github.com/HyperSpace-CW/Notification-App/pkg/logger"
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

	//brokers := []string{"kafka:9092"}
	//topic := "notifications"
	//group := "notification-group"
	//
	//// Инициализация Kafka Consumer
	//consumer, err := consumer.NewKafkaConsumer(brokers, topic, group, notificationService)
	//if err != nil {
	//	log.Fatal("Failed to create KafkaConsumer: %v", zap.Error(err))
	//}
	//
	//// Инициализация Kafka Producer
	//producer, err := producer.NewKafkaProducer(brokers, topic)
	//if err != nil {
	//	log.Fatal("Failed to create KafkaProducer: %v", zap.Error(err))
	//}

	//ctx, cancel := context.WithCancel(context.Background())
	//go consumer.Start(ctx)

	go func() {
		if err := httpServer.Run(); err != nil {
			log.Fatal(fmt.Sprintf("error occurred while running HTTP server: %v", err))
		}
	}()

	log.Info("Starting the application")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	<-quit

	//cancel()
	//if err := consumer.Close(); err != nil {
	//	log.Info("Error closing KafkaConsumer: %v", zap.Error(err))
	//}
	//if err := producer.Close(); err != nil {
	//	log.Info("Error closing KafkaProducer: %v", zap.Error(err))
	//}

	log.Info("shutdown HTTP server...")
	if err := httpServer.Shutdown(); err != nil {
		log.Error(fmt.Sprintf("failed to shutdown HTTP server: %v", err))
	} else {
		log.Info("HTTP server successfully shutdown")
	}
}
