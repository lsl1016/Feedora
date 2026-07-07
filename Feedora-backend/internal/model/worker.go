package model

import "time"

// WorkerEventRecord Worker 幂等消费记录表（阶段二保证同一事件不重复处理）。
// 阶段一预留表结构。
type WorkerEventRecord struct {
	ID           int64     `gorm:"primaryKey;column:id"`
	EventID      string    `gorm:"column:event_id;size:64;uniqueIndex:uk_event_worker"`
	WorkerName   string    `gorm:"column:worker_name;size:64;uniqueIndex:uk_event_worker;index:idx_worker_created"`
	Status       string    `gorm:"column:status;size:32;default:success"`
	ErrorMessage string    `gorm:"column:error_message;size:1000"`
	CreatedAt    time.Time `gorm:"column:created_at;index:idx_worker_created"`
}

func (WorkerEventRecord) TableName() string { return "worker_event_records" }
