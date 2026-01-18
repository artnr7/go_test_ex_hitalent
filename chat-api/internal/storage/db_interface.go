package storage

import (
	"chat-api/internal/models"
)

type ChatStorage interface {
	CreateChat(chat *models.Chat) error
	GetChatWithMessages(chatID uint, limit uint) error
	SendMessages(chatID uint, message *models.Message) error
	DeleteChat(chatID uint) error
}
