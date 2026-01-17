package main

import (
	"chat-api/internal/storage"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=chatapidb port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to connect to database:", err)
	}

	chatStorage := storage.NewChatStorage(db)
	chat:= 

}
