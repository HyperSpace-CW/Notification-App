package kafka

type EmailNotification struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
