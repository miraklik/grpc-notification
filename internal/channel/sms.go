package channel

import (
	"context"
	"log"
	"notification_service/internal/models"
	"time"
)

type SmsChannel struct{}

func NewSmsChannel() *SmsChannel {
	return &SmsChannel{}
}

func (s *SmsChannel) Send(ctx context.Context, notification *models.Notification) error {
	time.Sleep(500 * time.Millisecond)

	log.Printf("Sms sent to user %s: %s\n", notification.UserID, notification.Message)
	return nil
}

func (s *SmsChannel) GetType() models.NotificationType {
	return models.TypeSMS
}
