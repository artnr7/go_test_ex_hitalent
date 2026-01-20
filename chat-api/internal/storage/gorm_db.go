// Package storage implements GORM PostgreSQL database processing
package storage

import (
	"chat-api/internal/models"
	"errors"

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
) (*models.ChatWithMessages, error) {
	if limit > 100 {
		limit = 100
	}

	var chatWMessages models.ChatWithMessages
	err := s.db.Where("id = ", chatID).First(&chatWMessages.Chat).Error
	if err != nil {
		return &models.ChatWithMessages{}, err
	}

	err = s.db.Where("chatID = ?", chatWMessages.Chat.ID).
		Order("created_at DESC").
		Find(&chatWMessages.Messages).
		Error

	return &chatWMessages, err
}

func (s *chatStorage) SendMessage(chatID uint, message *models.Message) error {
	message.ChatID = chatID
	return s.db.Create(message).Error
}

func (s *chatStorage) DeleteChat(chatID uint) error {
	res := s.db.Where("id = ?", chatID).Delete(&models.Chat{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected <= 0 {
		return errors.New("chat not found")
	}

	return nil
}
