package service

import (
	"log"
	"notification_service/internal/models"

	"gorm.io/gorm"
)

type NotesService struct {
	db *gorm.DB
}

func NewNotesService(db *gorm.DB) *NotesService {
	return &NotesService{db: db}
}

func (s *NotesService) CreateNote(note *models.Notification) error {
	if err := s.db.Create(note).Error; err != nil {
		log.Printf("Error creating note: %v", err)
		return err
	}

	log.Println("Note created successfully")
	return nil
}

func (s *NotesService) GetAllNotes() ([]*models.Notification, error) {
	var notes []*models.Notification
	if err := s.db.Find(&notes).Error; err != nil {
		log.Printf("Error getting notes: %v", err)
		return nil, err
	}

	return notes, nil
}

func (s *NotesService) GetNotesById(id int) (*models.Notification, error) {
	var note *models.Notification
	if err := s.db.Find(&note).Where("id = ?", id).Error; err != nil {
		log.Printf("Error getting note: %v", err)
		return nil, err
	}

	return note, nil
}

func (s *NotesService) UpdateNote(id int, note *models.Notification) error {
	if err := s.db.Where("id = ?", id).Updates(note).Error; err != nil {
		log.Printf("Error updating note: %v", err)
		return err
	}

	log.Println("Note updated successfully")
	return nil
}

func (s *NotesService) DeleteNote(id int) error {
	if err := s.db.Delete(&models.Notification{}, id).Error; err != nil {
		log.Printf("Error deleting note: %v", err)
		return err
	}

	log.Println("Note deleted successfully")
	return nil
}
