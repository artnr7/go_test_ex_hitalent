// Package handlers implement chat handler methods
package handlers

import (
	"chat-api/internal/models"
	"chat-api/internal/storage"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type ChatHandler struct {
	storage   storage.ChatStorage
	validator *validator.Validate
}

func NewChatHandler(storage storage.ChatStorage) *ChatHandler {
	return &ChatHandler{
		storage:   storage,
		validator: validator.New(),
	}
}

func decodeRequestBodyJSON(
	w http.ResponseWriter,
	r *http.Request,
	requestBody any,
) error {
	if json.NewDecoder(r.Body).Decode(requestBody) != nil {
		http.Error(w, "Invalid JSON file", http.StatusBadRequest)
		return errors.New("bad request")
	}
	return nil
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Title string `json:"title" validate:"required,min=1,max=200"`
	}

	if decodeRequestBodyJSON(w, r, &requestBody) != nil {
		return
	}

	if err := h.validator.Struct(&requestBody); err != nil {
		http.Error(w, "Title must be 1 - 200 chars", http.StatusBadRequest)
		return
	}

	chat := models.Chat{Title: requestBody.Title}

	// db process
	if err := h.storage.CreateChat(&chat); err != nil {
		http.Error(
			w,
			"Failed to create chat in storage",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func chatIDParse(w http.ResponseWriter, r *http.Request, chatID *int) error {
	var err error
	*chatID, err = strconv.Atoi(r.PathValue("id"))
	if err != nil || *chatID < 0 {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return errors.New("bad parse")
	}
	return nil
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	// url
	var chatID int
	if chatIDParse(w, r, &chatID) != nil {
		return
	}

	// request body
	var requestBody struct {
		Text string `json:"text" validate:"required,min=1,max=5000"`
	}

	if decodeRequestBodyJSON(w, r, &requestBody) != nil {
		return
	}

	if h.validator.Struct(&requestBody) != nil {
		http.Error(w, "Text must be 1-5000 chars", http.StatusBadRequest)
		return
	}

	// db process
	message := models.Message{Text: requestBody.Text}
	if error := h.storage.SendMessage(uint(chatID), &message); error != nil {
		if strings.Contains(error.Error(), "doesn't exist") {
			http.Error(w, "Chat not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
		}
		return
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (h *ChatHandler) GetChatWithMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	// url
	var chatID int
	if chatIDParse(w, r, &chatID) != nil {
		return
	}

	limitStr := r.URL.Query().Get("limit")
	var limit int
	if limitStr == "" {
		limit = 20
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			http.Error(w, "Limit must be 1-100 messages", http.StatusBadRequest)
			return
		}
	}

	chatWMessages, err := h.storage.GetChatWithMessages(
		uint(chatID),
		uint(limit),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Chat not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get chat", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatWMessages)
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	// url
	var chatID int
	if chatIDParse(w, r, &chatID) != nil {
		return
	}

	err := h.storage.DeleteChat(uint(chatID))
	if err != nil {
		if strings.Contains(err.Error(), "chat not found") {
			http.Error(w, "Chat not found", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
