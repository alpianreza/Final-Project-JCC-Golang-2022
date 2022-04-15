package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Comments struct {
		gorm.Model
		ID          uint          `json:"id"`
		PostComment string        `json:"post_comment"`
		UserId      uint          `son:"user_id"`
		Users       Users         `json:"-"`
		Publish     bool          `json:"publish"`
		CreatedAt   time.Time     `json:"created_at"`
		UpdatedAt   time.Time     `json:"updated_at"`
		Meta        []CommentMeta `json:"-"`
	}
	CommentMeta struct {
		MetaID    uint   `json:"meta_id"`
		CommentID uint   `json:"comment_id"`
		MetaKey   string `json:"meta_key"`
		MetaValue string `json:"meta_value"`
	}
)

func (c *Comments) CommentSave(commentMeta map[string]string, db *gorm.DB) (*Comments, error) {
	ps := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			ps.Rollback()
		}
	}()

	if err := ps.Error; err != nil {
		return &Comments{}, err
	}
	var err error = ps.Create(&c).Error
	if err != nil {
		ps.Rollback()
		return &Comments{}, err
	}
	for key, value := range commentMeta {
		metaInput := CommentMeta{
			CommentID: c.ID,
			MetaKey:   key,
			MetaValue: value,
		}
		var err error = ps.Create(&metaInput).Error
		if err != nil {
			ps.Rollback()
			return &Comments{}, err
		}
	}

	return c, ps.Commit().Error
}
