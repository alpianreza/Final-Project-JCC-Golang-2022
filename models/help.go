package models

import "time"

type Help struct {
	ID        int       `gorm:"primary_key" json:"id"`
	User      Users     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
