package models

import (
	"time"

	"gorm.io/gorm"
)

type Help struct {
	gorm.Model
	ID        int       `json:"id"`
	User      Users     `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
