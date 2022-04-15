package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Posts struct {
		gorm.Model
		ID          int        `json:"id"`
		PostTitle   string     `json:"post_title"`
		PostContent string     `json:"post_content"`
		Publish     bool       `json:"publish"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		UserID      uint       `json:"user_id"`
		CategoryID  int        `json:"category_id"`
		Users       Users      `json:"-"`
		Category    Category   `json:"-"`
		Comments    []Comments `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		Meta        PostMeta   `json:"-"`
	}
	PostMeta struct {
		MetaID    int    `json:"meta_id"`
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
