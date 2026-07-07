package repository

import (
	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 后台管理数据访问。
type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// ListPosts 分页查询帖子（后台，可按状态过滤）。
func (r *AdminRepository) ListPosts(status string, offset, limit int) ([]model.Post, int64) {
	q := r.db.Model(&model.Post{})
	if status != "" && status != "all" {
		q = q.Where("status = ?", status)
	}
	var total int64
	q.Count(&total)
	var rows []model.Post
	q.Order("id DESC").Offset(offset).Limit(limit).Find(&rows)
	return rows, total
}

// ListLogs 分页查询操作日志。
func (r *AdminRepository) ListLogs(offset, limit int) ([]model.OperationLog, int64) {
	var total int64
	r.db.Model(&model.OperationLog{}).Count(&total)
	var rows []model.OperationLog
	r.db.Order("id DESC").Offset(offset).Limit(limit).Find(&rows)
	return rows, total
}

// Stats 统计各实体数量。
func (r *AdminRepository) Stats() (users, posts, comments, circles int64) {
	r.db.Model(&model.User{}).Count(&users)
	r.db.Model(&model.Post{}).Count(&posts)
	r.db.Model(&model.Comment{}).Count(&comments)
	r.db.Model(&model.Circle{}).Count(&circles)
	return
}
