package model

import "time"

// Notification 通知表（阶段二由 Notification Worker 异步生成）。
// 阶段一预留表结构，暂不产生数据。
type Notification struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	UserID     int64     `gorm:"column:user_id;index:idx_user_read_created;index:idx_user_created"`
	ActorID    int64     `gorm:"column:actor_id;default:0"`
	Type       string    `gorm:"column:type;size:64"`
	Title      string    `gorm:"column:title;size:200"`
	Content    string    `gorm:"column:content;size:1000"`
	TargetType string    `gorm:"column:target_type;size:64"`
	TargetID   int64     `gorm:"column:target_id"`
	TargetURL  string    `gorm:"column:target_url;size:512"`
	ReadStatus string    `gorm:"column:read_status;size:32;default:unread;index:idx_user_read_created"`
	CreatedAt  time.Time `gorm:"column:created_at;index:idx_user_read_created;index:idx_user_created"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Notification) TableName() string { return "notifications" }
