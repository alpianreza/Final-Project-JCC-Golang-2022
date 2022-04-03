package models

import "time"

type Category struct {
	ID        uint      `gorm:"primary_key;index:,unique" json:"id"`
	Category  string    `gorm:"not null" json:"category_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     *[]Posts  `json:"post" gorm:"foreignKey:category_id;references:id"`
}
