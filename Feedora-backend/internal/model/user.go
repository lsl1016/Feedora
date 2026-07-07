package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表。
type User struct {
	ID             int64          `gorm:"primaryKey;column:id"`
	Account        string         `gorm:"column:account;size:64;uniqueIndex:uk_account"`
	PasswordHash   string         `gorm:"column:password_hash;size:255"`
	Nickname       string         `gorm:"column:nickname;size:64"`
	Avatar         string         `gorm:"column:avatar;size:512"`
	Bio            string         `gorm:"column:bio;size:500"`
	Role           string         `gorm:"column:role;size:32;default:user"`
	Status         string         `gorm:"column:status;size:32;default:normal;index"`
	FollowerCount  int64          `gorm:"column:follower_count;default:0"`
	FollowingCount int64          `gorm:"column:following_count;default:0"`
	PostCount      int64          `gorm:"column:post_count;default:0"`
	CommentCount   int64          `gorm:"column:comment_count;default:0"`
	LikeCount      int64          `gorm:"column:like_count;default:0"`
	PointCount     int64          `gorm:"column:point_count;default:0"`
	Level          int            `gorm:"column:level;default:1"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string { return "users" }
