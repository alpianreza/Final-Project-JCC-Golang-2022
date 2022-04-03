package models

import "time"

type Help struct {
	ID        int       `gorm:"primary_key" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
