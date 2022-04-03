package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Posts struct {
		ID          int        `gorm:"primary_key" json:"id"`
		PostTitle   string     `gorm:"not null,colomn:title" json:"post_title"`
		PostContent string     `gorm:"not null" json:"post_content"`
		Publish     bool       `gorm:"not null" json:"publish"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		UserID      uint       `gorm:"not null" json:"user_id"`
		User        User       `gorm:"references:user_id" json:"user"`
		CategoryID  int        `gorm:"not null" json:"category_id"`
		Category    Category   `gorm:"foreignKey:category_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category"`
		Comments    []Comments `gorm:"many2many:post_comments;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments"`
		Meta        []PostMeta `gorm:"foreignKey:post_id;" json:"meta"`
	}
	PostMeta struct {
		MetaID    int    `gorm:"primaryKey,index" json:"meta_id"`
		PostID    int    `json:"post_id"`
		MetaKey   string `json:"meta_key"`
		MetaValue string `json:"meta_value"`
	}
)

func (p *Posts) PSave(postMeta map[string]string, db *gorm.DB) (*Posts, error) {
	ps := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			ps.Rollback()
		}
	}()

	if err := ps.Error; err != nil {
		return &Posts{}, err
	}
	var err error = ps.Create(&p).Error
	if err != nil {
		ps.Rollback()
		return &Posts{}, err
	}
	for key, value := range postMeta {
		metaInput := PostMeta{
			PostID:    p.ID,
			MetaKey:   key,
			MetaValue: value,
		}
		var err error = ps.Create(&metaInput).Error
		if err != nil {
			ps.Rollback()
			return &Posts{}, err
		}
	}

	return p, ps.Commit().Error
}
