package handlers

import (
	"context"
	"errors"
	"log"
	"notification_service/internal/models"
	"notification_service/internal/queue"
	"notification_service/internal/service"
	"time"
)

type NotificationHandler struct {
	services *service.NotificationsService
	queue    *queue.NatsQueue
}

func NewNotificationHandler(services *service.NotificationsService, queue *queue.NatsQueue) *NotificationHandler {
	return &NotificationHandler{
		services: services,
		queue:    queue,
	}
}

func (n *NotificationHandler) SendNotification(ctx context.Context, userID int64, notType, message string, priority int32, scheduledAt int64) (int64, error) {
	if message == "" || notType == "" {
		log.Println("Missing required fields: userID, message, type")
		return 0, nil
	}

	if priority < 1 || priority > 5 {
		log.Println("Priority must be between 1 and 5")
		return 0, errors.New("priority must be between 1 and 5")
	}

	var TypeNotes models.NotificationType

	switch notType {
	case "email":
		TypeNotes = models.TypeEmail
	case "push":
		TypeNotes = models.TypePush
	case "sms":
		TypeNotes = models.TypeSMS
	default:
		log.Println("Invalid notification type")
		return 0, errors.New("invalid notification type")
	}

	notification := &models.Notification{
		UserID:      userID,
		Type:        TypeNotes,
		Message:     message,
		Priority:    priority,
		ScheduledAt: time.Unix(scheduledAt, 0),
		Status:      models.StatusPending,
		Attempts:    0,
		LastError:   "",
		DeliveredAt: nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if scheduledAt == 0 {
		notification.ScheduledAt = time.Now()
	}

	err := n.services.CreateNote(notification)
	if err != nil {
		log.Println("Error creating note: ", err)
		return 0, err
	}

	if err := n.queue.Push(notification); err != nil {
		log.Println("Error pushing note to queue: ", err)
		return 0, err
	}

	return notification.ID, nil
}

func (n *NotificationHandler) GetStatus(ctx context.Context, notificationID int64) (*models.Notification, error) {
	notification, err := n.services.GetNotesById(notificationID)
	if err != nil {
		log.Println("Error getting note: ", err)
		return nil, err
	}

	return notification, nil
}
