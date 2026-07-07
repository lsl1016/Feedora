package model

import "time"

// UserPointLog 积分流水表（阶段二由 Growth Worker 异步发放）。
// 唯一索引用于防止重复发放积分。阶段一预留表结构。
type UserPointLog struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex:uk_user_action_biz;index:idx_user_created"`
	Action    string    `gorm:"column:action;size:64;uniqueIndex:uk_user_action_biz"`
	Point     int       `gorm:"column:point"`
	BizType   string    `gorm:"column:biz_type;size:64;uniqueIndex:uk_user_action_biz"`
	BizID     int64     `gorm:"column:biz_id;uniqueIndex:uk_user_action_biz"`
	Remark    string    `gorm:"column:remark;size:255"`
	CreatedAt time.Time `gorm:"column:created_at;index:idx_user_created"`
}

func (UserPointLog) TableName() string { return "user_point_logs" }
