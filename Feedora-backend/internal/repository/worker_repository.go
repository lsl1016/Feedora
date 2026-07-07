package repository

import (
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// WorkerRepository Worker 幂等消费记录数据访问。
type WorkerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

// Claim 尝试占用一次消费记录。返回 true 表示首次消费可继续，
// false 表示该 (eventID, workerName) 已处理过（幂等拦截）。
func (r *WorkerRepository) Claim(eventID, workerName string) bool {
	rec := &model.WorkerEventRecord{
		EventID: eventID, WorkerName: workerName, Status: "success", CreatedAt: time.Now(),
	}
	// 依赖 uk_event_worker 唯一索引：重复插入会失败。
	res := r.db.Create(rec)
	return res.Error == nil && res.RowsAffected > 0
}

// Release 处理失败时删除占用记录，允许后续重试。
func (r *WorkerRepository) Release(eventID, workerName string) {
	r.db.Where("event_id = ? AND worker_name = ?", eventID, workerName).Delete(&model.WorkerEventRecord{})
}
