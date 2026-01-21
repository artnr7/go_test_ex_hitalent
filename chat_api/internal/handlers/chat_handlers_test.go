package handlers

import (
	"httptest"
	"net/http"
	"testing"
)

func TestCreateChat(t *testing.T){
	mux:= http.NewServeMux()

	mux.HandleFunc("POST /chats/", ChatHandler.)
}
