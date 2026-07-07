package model

import "time"

// EventOutbox 事件发件箱表（阶段二用于可靠投递 Kafka）。
// 阶段一预留表结构，业务尚未写入。
type EventOutbox struct {
	ID           int64      `gorm:"primaryKey;column:id"`
	EventID      string     `gorm:"column:event_id;size:64;uniqueIndex:uk_event_id"`
	Topic        string     `gorm:"column:topic;size:128"`
	EventType    string     `gorm:"column:event_type;size:64"`
	AggregateID  int64      `gorm:"column:aggregate_id"`
	UserID       int64      `gorm:"column:user_id;default:0"`
	TraceID      string     `gorm:"column:trace_id;size:64"`
	Payload      string     `gorm:"column:payload;type:json"`
	Status       string     `gorm:"column:status;size:32;default:pending;index:idx_status_next_retry"`
	RetryCount   int        `gorm:"column:retry_count;default:0"`
	NextRetryAt  *time.Time `gorm:"column:next_retry_at;index:idx_status_next_retry"`
	DispatchedAt *time.Time `gorm:"column:dispatched_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;index"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (EventOutbox) TableName() string { return "event_outbox" }
