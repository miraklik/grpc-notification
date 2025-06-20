package service

import (
	"log"
	"notification_service/internal/models"

	"gorm.io/gorm"
)

type NotificationsService struct {
	db *gorm.DB
}

func NewNotificationsService(db *gorm.DB) *NotificationsService {
	return &NotificationsService{db: db}
}

func (s *NotificationsService) CreateNote(note *models.Notification) error {
	if err := s.db.Create(note).Error; err != nil {
		log.Printf("Error creating note: %v", err)
		return err
	}

	log.Println("Note created successfully")
	return nil
}

func (s *NotificationsService) GetAllNotes() ([]*models.Notification, error) {
	var notes []*models.Notification
	if err := s.db.Find(&notes).Error; err != nil {
		log.Printf("Error getting notes: %v", err)
		return nil, err
	}

	return notes, nil
}

func (s *NotificationsService) GetNotesById(id string) (*models.Notification, error) {
	var note *models.Notification
	if err := s.db.Find(&note).Where("id = ?", id).Error; err != nil {
		log.Printf("Error getting note: %v", err)
		return nil, err
	}

	return note, nil
}

func (s *NotificationsService) UpdateNote(id int, note *models.Notification) error {
	if err := s.db.Where("id = ?", id).Updates(note).Error; err != nil {
		log.Printf("Error updating note: %v", err)
		return err
	}

	log.Println("Note updated successfully")
	return nil
}

func (s *NotificationsService) DeleteNote(id int) error {
	if err := s.db.Delete(&models.Notification{}, id).Error; err != nil {
		log.Printf("Error deleting note: %v", err)
		return err
	}

	log.Println("Note deleted successfully")
	return nil
}
