package service

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
)

// NotificationService 通知读取业务逻辑（写入由 Worker 完成）。
type NotificationService struct {
	notifs *repository.NotificationRepository
}

func NewNotificationService(notifs *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifs: notifs}
}

// List 分页查询用户通知。
func (s *NotificationService) List(userID int64, category string, page, size int) ([]dto.NotificationItem, int64) {
	page, size = normPage(page, size)
	rows, total := s.notifs.List(userID, category, offset(page, size), size)
	list := make([]dto.NotificationItem, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToNotificationItem(&rows[i]))
	}
	return list, total
}

// MarkRead 标记单条已读。
func (s *NotificationService) MarkRead(id, userID int64) {
	s.notifs.MarkRead(id, userID)
}

// MarkAllRead 全部已读。
func (s *NotificationService) MarkAllRead(userID int64) {
	s.notifs.MarkAllRead(userID)
}

// UnreadCount 未读数。
func (s *NotificationService) UnreadCount(userID int64) int64 {
	return s.notifs.UnreadCount(userID)
}
