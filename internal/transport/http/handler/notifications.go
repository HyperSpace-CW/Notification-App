package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (h *Handler) initNotificationsRoutes(api fiber.Router) {
	serviceRoute := api.Group("/notifications")
	{
		serviceRoute.Post("/send", h.SendCodeToEmail)
	}
}

type SendCodeToEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *Handler) SendCodeToEmail(ctx *fiber.Ctx) error {
	log.Info("Parse request")
	var req SendCodeToEmailRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Errorf("ctx.BodyParser: %w", err).Error(),
		)
	}

	log.Info("Sending code to email")
	err := h.notificationService.SendCodeToEmail(req.Email, req.Code)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Errorf("h.notificationService.SendCodeToEmail: %w", err).Error(),
		)
	}

	return nil
}
