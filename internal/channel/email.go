package channel

import (
	"context"
	"notification_service/internal/models"
	"time"

	"log"
)

type EmailChannel struct{}

func NewEmailChannel() *EmailChannel {
	return &EmailChannel{}
}

func (n *EmailChannel) Send(ctx context.Context, notification *models.Notification) error {
	time.Sleep(500 * time.Microsecond)

	log.Printf("Email sent to user %s: %s\n", notification.UserID, notification.Message)
	return nil
}

func (n *EmailChannel) GetType() models.NotificationType {
	return models.TypeEmail
}
