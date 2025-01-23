package http

import (
	"errors"
	"fmt"
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	"github.com/HyperSpace-CW/Notification-App/internal/transport/http/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

type Server struct {
	addr string

	notificationService services.NotificationService

	logger *zap.Logger
	app    *fiber.App
}

type Config struct {
	Addr string

	NotificationService services.NotificationService

	Logger *zap.Logger
}

func NewServer(cfg Config) *Server {
	server := &Server{
		addr:                cfg.Addr,
		logger:              cfg.Logger,
		notificationService: cfg.NotificationService,
		app:                 nil,
	}

	server.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	server.init()

	return server
}

func (s *Server) Run() error {
	if err := s.app.Listen(s.addr); err != nil {
		return fmt.Errorf("listening HTTP server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	return nil
}

func (s *Server) init() {
	s.app.Use(cors.New())
	s.app.Use(requestid.New())
	s.app.Use(func(ctx *fiber.Ctx) error {
		ctx.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return ctx.Next()
	})

	s.setHandlers()
}

// TODO: set path
func (s *Server) setHandlers() {
	handlerV1 := handler.NewHandler(handler.Config{
		NotificationService: s.notificationService,
		Logger:              s.logger,
	})
	handlerV1.Init(s.app)
}

type errResp struct {
	Error bool   `json:"error"`
	Data  string `json:"data"`
}

func (s *Server) errorHandler(ctx *fiber.Ctx, err error) error {
	//requestID := utils.GetRequestID(ctx)

	var statusCode = fiber.StatusInternalServerError

	resp := errResp{
		Error: true,
		Data:  err.Error(),
	}

	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &fiberErr):
		statusCode = fiberErr.Code
		resp.Data = fiberErr.Message
	}

	s.logger.Error(
		err.Error(),
		//zap.String("request_id", requestID),
		zap.String("method", ctx.Method()),
		zap.String("path", ctx.Path()),
		zap.Int("status", statusCode),
	)

	body, _ := jsoniter.Marshal(resp)

	if respondErr := ctx.Status(statusCode).Send(body); respondErr != nil {
		s.logger.Error(
			"sending error response",
			zap.String("error", err.Error()),
			//zap.String("request_id", requestID),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.Int("status", statusCode),
		)
	}

	return nil
}
