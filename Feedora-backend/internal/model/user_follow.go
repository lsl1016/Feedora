package model

import "time"

// UserFollow 用户关注关系（follower 关注 followee）。
type UserFollow struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	FollowerID int64     `gorm:"column:follower_id;uniqueIndex:uk_user_follow,priority:1"` // 关注发起方。
	FolloweeID int64     `gorm:"column:followee_id;uniqueIndex:uk_user_follow,priority:2"` // 被关注方。
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (UserFollow) TableName() string { return "user_follows" }
