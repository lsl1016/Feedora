package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论表。
type Comment struct {
	ID            int64          `gorm:"primaryKey;column:id"`
	PostID        int64          `gorm:"column:post_id;index:idx_post_parent;index:idx_post_created"`
	UserID        int64          `gorm:"column:user_id;index:idx_user_created"`
	ParentID      int64          `gorm:"column:parent_id;default:0;index:idx_post_parent"`
	RootID        int64          `gorm:"column:root_id;default:0;index"`
	ReplyToUserID *int64         `gorm:"column:reply_to_user_id"`
	Content       string         `gorm:"column:content;size:2000"`
	Status        string         `gorm:"column:status;size:32;default:normal"`
	LikeCount     int64          `gorm:"column:like_count;default:0"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Comment) TableName() string { return "comments" }
