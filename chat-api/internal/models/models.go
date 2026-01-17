// Package models is implement of models structs
package models

import "time"

type Chat struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:200;not null"`
	CreatedAt time.Time
}

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	ChatID    uint   `gorm:"constraint:OnDelete:CASCADE"`
	Text      string `gorm:"size:5000;not null"`
	CreatedAt time.Time
}
