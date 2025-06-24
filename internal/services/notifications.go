package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/HyperSpace-CW/Notification-App/config"
	"gopkg.in/gomail.v2"
)

type NotificationService interface {
	SendCodeToEmail(ctx context.Context, email, code string) error
}

type notificationService struct {
	cfg *config.Config
}

func NewNotificationService(config *config.Config) NotificationService {
	return &notificationService{
		cfg: config,
	}
}

func (s *notificationService) SendCodeToEmail(ctx context.Context, email, code string) error {
	log.Printf("Sending email to %s with code %s", email, code)
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.cfg.Email.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Verify Code")
	msg.SetBody("text/html", fmt.Sprintf(`
        <html>
            <body>
                <h2>Код подтверждения</h2>
                <p>Ваш код подтверждения: <strong>%s</strong></p>
                <p>Спасибо за использование нашего сервиса.</p>
            </body>
        </html>`, code))

	n := gomail.NewDialer("smtp.gmail.com", 587, s.cfg.Email.Username, s.cfg.Email.Password)
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := n.DialAndSend(msg); err != nil {
		return fmt.Errorf("s.SendCodeToEmail err: %w", err)
	}
	return nil
}
