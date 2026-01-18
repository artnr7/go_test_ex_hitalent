package storage

import (
	"chat-api/internal/models"

	"gorm.io/gorm"
)

type chatStorage struct {
	db *gorm.DB
}

func NewChatStorage(db *gorm.DB) ChatStorage {
	return &chatStorage{db: db}
}

func (s *chatStorage) CreateChat(chat *models.Chat) error {
	return s.db.Create(chat).Error
}

func (s *chatStorage) GetChatWithMessages(
	chatID uint,
	limit uint,
) (*models.Chat, error) {
	if limit > 100 {
		limit = 100
	}

	var chat models.Chat

	err := s.db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc").Limit(int(limit))
	}).First(&chat, chatID).Error

	return &chat, err
}

func (s *chatStorage) SendMessage(chatID uint, message *models.Message) error {
	message.ChatID = chatID
	return s.db.Create(message).Error
}

func (s *chatStorage) DeleteChat(chatID uint) error {
	return s.db.Where("id = ?", chatID).Delete(&models.Chat{}).Error
}
