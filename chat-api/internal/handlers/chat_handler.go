// Package handlers implement chat handler methods
package handlers

import (
	"chat-api/internal/models"
	"chat-api/internal/storage"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
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

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title string `json:"title" validate:"requires,min=1,max=200"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON file", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(&request); err != nil {
		http.Error(w, "Title must be 1 - 200 chars", http.StatusBadRequest)
		return
	}

	chat := models.Chat{Title: request.Title}

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
