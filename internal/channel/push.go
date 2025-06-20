package channel

import (
	"context"
	"log"
	"notification_service/internal/models"
	"time"
)

type PushChannel struct{}

func NewPushChannel() *PushChannel {
	return &PushChannel{}
}

func (p *PushChannel) Send(ctx context.Context, notification *models.Notification) error {
	time.Sleep(500 * time.Microsecond)

	log.Printf("Push sent to user %s: %s\n", notification.UserID, notification.Message)
	return nil
}

func (p *PushChannel) GetType() models.NotificationType {
	return models.TypePush
}
