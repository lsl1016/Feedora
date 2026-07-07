package repository

import (
	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// NotificationRepository 通知数据访问。
type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Insert 写入一条通知。
func (r *NotificationRepository) Insert(n *model.Notification) error {
	return r.db.Create(n).Error
}

// List 分页查询用户通知（按时间倒序）。
func (r *NotificationRepository) List(userID int64, category string, offset, limit int) ([]model.Notification, int64) {
	q := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if category != "" && category != "all" {
		q = q.Where("type = ?", category)
	}
	var total int64
	q.Count(&total)
	var rows []model.Notification
	q.Order("id DESC").Offset(offset).Limit(limit).Find(&rows)
	return rows, total
}

// MarkRead 标记单条已读，返回是否有更新。
func (r *NotificationRepository) MarkRead(id, userID int64) bool {
	res := r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ? AND read_status = ?", id, userID, "unread").
		Update("read_status", "read")
	return res.RowsAffected > 0
}

// MarkAllRead 标记用户全部通知已读。
func (r *NotificationRepository) MarkAllRead(userID int64) {
	r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read_status = ?", userID, "unread").
		Update("read_status", "read")
}

// UnreadCount 统计未读数。
func (r *NotificationRepository) UnreadCount(userID int64) int64 {
	var n int64
	r.db.Model(&model.Notification{}).Where("user_id = ? AND read_status = ?", userID, "unread").Count(&n)
	return n
}
