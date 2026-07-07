package model

import "time"

// OperationLog 后台操作日志表。
type OperationLog struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	AdminID    int64     `gorm:"column:admin_id;index"`
	AdminName  string    `gorm:"column:admin_name;size:64"`
	Action     string    `gorm:"column:action;size:64"`
	TargetType string    `gorm:"column:target_type;size:64"`
	TargetID   int64     `gorm:"column:target_id"`
	Detail     string    `gorm:"column:detail;size:1000"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (OperationLog) TableName() string { return "operation_logs" }
