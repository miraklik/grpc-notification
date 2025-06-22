package models

import (
	"time"
)

type NotificationStatus string

const (
	StatusPending    NotificationStatus = "PENDING"
	StatusInProgress NotificationStatus = "IN_PROGRESS"
	StatusRetrying   NotificationStatus = "RETRYING"
	StatusSent       NotificationStatus = "SENT"
	StatusFailed     NotificationStatus = "FAILED"
)

type NotificationType string

const (
	TypeEmail NotificationType = "email"
	TypePush  NotificationType = "push"
	TypeSMS   NotificationType = "sms"
)

type Notification struct {
	ID          int64              `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      int64              `json:"user_id"`
	Type        NotificationType   `json:"type"`
	Message     string             `json:"message"`
	Priority    int32              `json:"priority"`
	ScheduledAt time.Time          `json:"scheduled_at"`
	Status      NotificationStatus `json:"status"`
	Attempts    int32              `json:"attempts"`
	LastError   string             `json:"last_error"`
	DeliveredAt *time.Time         `json:"delivered_at"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
