package channel

import (
	"context"
	"errors"
	"notification_service/internal/models"
	"time"

	"log"
)

type NotesChannel struct{}

func NewNotesChannel() *NotesChannel {
	return &NotesChannel{}
}

func (n *NotesChannel) Send(ctx context.Context, notification *models.Notification) error {
	time.Sleep(500 * time.Microsecond)

	if notification.Priority < 3 {
		log.Printf("Note: %s", notification.Message)
		return errors.New("note sent")
	}

	log.Printf("Email sent to user %s: %s\n", notification.UserID, notification.Message)
	return nil
}

func (n *NotesChannel) GetType() models.NotificationType {
	return models.TypeEmail
}
