package channel

import (
	"context"
	"notification_service/internal/models"
)

type NotificationChannel interface {
	Send(ctx context.Context, notification *models.Notification) error
	GetType() models.NotificationType
}
