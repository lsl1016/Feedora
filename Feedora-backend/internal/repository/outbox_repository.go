package repository

import (
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// OutboxRepository event_outbox 事件表数据访问。
type OutboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

// Insert 写入一条待投递事件。
func (r *OutboxRepository) Insert(e *model.EventOutbox) error {
	return r.db.Create(e).Error
}

// FetchDispatchable 取出可投递的事件（pending 且到达重试时间）。
func (r *OutboxRepository) FetchDispatchable(limit int) ([]model.EventOutbox, error) {
	var rows []model.EventOutbox
	now := time.Now()
	err := r.db.Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", "pending", now).
		Order("id ASC").Limit(limit).Find(&rows).Error
	return rows, err
}

// MarkDispatched 标记为已投递。
func (r *OutboxRepository) MarkDispatched(id int64) {
	now := time.Now()
	r.db.Model(&model.EventOutbox{}).Where("id = ?", id).Updates(map[string]any{
		"status": "dispatched", "dispatched_at": now, "updated_at": now,
	})
}

// MarkRetry 记录一次投递失败并安排下次重试。
func (r *OutboxRepository) MarkRetry(id int64, retryCount int, nextRetryAt time.Time) {
	r.db.Model(&model.EventOutbox{}).Where("id = ?", id).Updates(map[string]any{
		"retry_count": retryCount, "next_retry_at": nextRetryAt, "updated_at": time.Now(),
	})
}

// MarkFailed 超过最大重试次数，标记为失败。
func (r *OutboxRepository) MarkFailed(id int64) {
	r.db.Model(&model.EventOutbox{}).Where("id = ?", id).Updates(map[string]any{
		"status": "failed", "updated_at": time.Now(),
	})
}
