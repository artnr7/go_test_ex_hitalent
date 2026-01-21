package handlers

import (
	"bytes"
	"chat-api/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
)

type MockChatStorage struct {
	dbWasCalled bool
}

func (m *MockChatStorage) CreateChat(chat *models.Chat) error {
	m.dbWasCalled = true
	chat.ID = 1
	return nil
}

func (m *MockChatStorage) GetChatWithMessages(
	chatID uint,
	limit uint,
) (*models.ChatWithMessages, error) {
	return &models.ChatWithMessages{}, nil
}

func (m *MockChatStorage) SendMessage(
	chatID uint,
	message *models.Message,
) error {
	return nil
}
func (m *MockChatStorage) DeleteChat(chatID uint) error { return nil }

func TestCreateChat(t *testing.T) {
	mockStorage := &MockChatStorage{}
	chatHandler := &ChatHandler{
		storage:   mockStorage,
		validator: validator.New(),
	}

	requestBody := []byte(`{"title":"Kazan"}`)
	request := httptest.NewRequest(
		http.MethodPost,
		"/chats/",
		bytes.NewReader(requestBody),
	)
	recorder := httptest.NewRecorder()

	chatHandler.CreateChat(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.True(t, mockStorage.dbWasCalled)

	var response models.Chat
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "Kazan", response.Title)
	require.Equal(t, uint(1), response.ID)
}
