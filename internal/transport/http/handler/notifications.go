package handler

import (
	"fmt"
	_ "github.com/HyperSpace-CW/Notification-App/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func (h *Handler) initNotificationsRoutes(api fiber.Router) {
	serviceRoute := api.Group("/notifications")
	{
		serviceRoute.Post("/send", h.SendCodeToEmail)
	}
	serviceRoute.Get("/swagger/*", fiberSwagger.WrapHandler)
}

// HTTPError представляет ошибку HTTP-ответа
type HTTPError struct {
	Message string `json:"message"`
}

type SendCodeToEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// SendCodeToEmail отправляет код подтверждения на email
// @Summary Отправить код на email
// @Tags notifications
// @Description Отправляет пользователю код подтверждения по email
// @Accept json
// @Produce json
// @Param data body SendCodeToEmailRequest true "Email и код подтверждения"
// @Success 200 {string} string "OK"
// @Failure 400 {object} HTTPError "Невалидный запрос или ошибка отправки письма"
// @Router /notifications/send [post]
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
	err := h.notificationService.SendCodeToEmail(ctx.Context(), req.Email, req.Code)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Errorf("h.notificationService.SendCodeToEmail: %w", err).Error(),
		)
	}

	return nil
}
