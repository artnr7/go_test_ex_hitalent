package main

import (
	"database/sql"
	"log"
	"net/http"

	"chat-api/internal/handlers"
	"chat-api/internal/storage"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// db
	dsn := "host=postgresql user=chat_api_user password=123 dbname=chat_api_db port=5432 sslmode=disable"
	log.Println("Database opening")
	psqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Fail to connect to database:", err)
	}
	defer psqlDB.Close()

	log.Println("Start migrating database")
	err = goose.Up(psqlDB, "./migrations")
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to connect to database:", err)
	}

	// dependencies
	chatStorage := storage.NewChatStorage(db)
	chatHandler := handlers.NewChatHandler(chatStorage)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /chats/", chatHandler.CreateChat)
	mux.HandleFunc("POST /chats/{id}/messages/", chatHandler.SendMessage)
	mux.HandleFunc("GET /chats/{id}", chatHandler.GetChatWithMessages)
	mux.HandleFunc("DELETE /chats/{id}", chatHandler.DeleteChat)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
