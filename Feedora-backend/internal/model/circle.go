package model

import (
	"time"

	"gorm.io/gorm"
)

// Circle 圈子表。
type Circle struct {
	ID             int64          `gorm:"primaryKey;column:id"`
	OwnerID        int64          `gorm:"column:owner_id;index"`
	Name           string         `gorm:"column:name;size:128;uniqueIndex:uk_circle_name"`
	Avatar         string         `gorm:"column:avatar;size:512"`
	Description    string         `gorm:"column:description;size:500"`
	Category       string         `gorm:"column:category;size:64"`
	JoinType       string         `gorm:"column:join_type;size:32;default:direct"`
	PostPermission string         `gorm:"column:post_permission;size:32;default:all"`
	Rules          string         `gorm:"column:rules;type:text"`
	Status         string         `gorm:"column:status;size:32;default:normal;index"`
	MemberCount    int64          `gorm:"column:member_count;default:0"`
	PostCount      int64          `gorm:"column:post_count;default:0"`
	FeaturedCount  int64          `gorm:"column:featured_count;default:0"`
	IsRecommended  bool           `gorm:"column:is_recommended;default:false;index"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Circle) TableName() string { return "circles" }

// CircleMember 圈子成员表。
type CircleMember struct {
	ID         int64      `gorm:"primaryKey;column:id"`
	CircleID   int64      `gorm:"column:circle_id;uniqueIndex:uk_circle_user;index:idx_circle_role;index:idx_circle_status"`
	UserID     int64      `gorm:"column:user_id;uniqueIndex:uk_circle_user;index"`
	Role       string     `gorm:"column:role;size:32;default:member;index:idx_circle_role"`
	Status     string     `gorm:"column:status;size:32;default:normal;index:idx_circle_status"`
	MuteReason string     `gorm:"column:mute_reason;size:255"`
	MutedUntil *time.Time `gorm:"column:muted_until"`
	JoinedAt   time.Time  `gorm:"column:joined_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

func (CircleMember) TableName() string { return "circle_members" }
