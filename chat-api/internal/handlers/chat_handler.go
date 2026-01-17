// Package handlers implement chat handler methods
package handlers

import (
	"chat-api/internal/storage"

	"github.com/go-playground/validator"
)

type ChatHandler struct {
	storage  storage.ChatStorage
	validate *validator.Validate
}

func 
