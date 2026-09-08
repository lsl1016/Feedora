package repository

import (
	"time"

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

// GetPostByID 按 ID 查询帖子，不存在返回 nil。
func (r *AdminRepository) GetPostByID(id int64) *model.Post {
	var p model.Post
	if err := r.db.First(&p, id).Error; err != nil {
		return nil
	}
	return &p
}

// UpdatePostStatus 更新帖子状态（上下架）。
func (r *AdminRepository) UpdatePostStatus(id int64, status string) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateUserStatus 更新用户状态（封禁 / 解禁）。
func (r *AdminRepository) UpdateUserStatus(id int64, status string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// AddLog 写入后台操作日志。
func (r *AdminRepository) AddLog(l *model.OperationLog) {
	l.CreatedAt = time.Now()
	r.db.Create(l)
}
