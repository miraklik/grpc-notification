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
	queue    *queue.Queue
}

func NewNotificationHandler(services *service.NotificationsService, queue *queue.Queue) *NotificationHandler {
	return &NotificationHandler{
		services: services,
		queue:    queue,
	}
}

func (n *NotificationHandler) SendNotification(ctx context.Context, userID, notType, message string, priority int32, scheduledAt int64) (string, error) {
	if userID == "" || message == "" || notType == "" {
		log.Println("Missing required fields: userID, message, type")
		return "", nil
	}

	if priority < 1 || priority > 5 {
		log.Println("Priority must be between 1 and 5")
		return "", errors.New("priority must be between 1 and 5")
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
		return "", errors.New("invalid notification type")
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
		return "", err
	}

	n.queue.Push(notification)

	return notification.ID, nil
}

func (n *NotificationHandler) GetStatus(ctx context.Context, notificationID string) (*models.Notification, error) {
	if notificationID == "" {
		log.Println("Missing required fields: notificationID")
		return nil, errors.New("missing required fields: notificationID")
	}

	notification, err := n.services.GetNotesById(notificationID)
	if err != nil {
		log.Println("Error getting note: ", err)
		return nil, err
	}

	return notification, nil
}
