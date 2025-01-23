package handler

import (
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	notificationService services.NotificationService
	logger              *zap.Logger
}

type Config struct {
	NotificationService services.NotificationService
	Logger              *zap.Logger
}

func NewHandler(cfg Config) *Handler {
	return &Handler{
		notificationService: cfg.NotificationService,
		logger:              cfg.Logger,
	}
}

func (h *Handler) Init(routeV1 fiber.Router) {
	h.initNotificationsRoutes(routeV1)
}
