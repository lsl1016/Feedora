package service

import (
	"context"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
)

// unreadCacheTTL 未读数缓存有效期。
const unreadCacheTTL = 10 * time.Minute

// NotificationService 通知读取业务逻辑（写入由 Worker 完成）。
type NotificationService struct {
	notifs *repository.NotificationRepository
	cache  *cache.Cache
}

func NewNotificationService(notifs *repository.NotificationRepository, cch *cache.Cache) *NotificationService {
	return &NotificationService{notifs: notifs, cache: cch}
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

// invalidateUnread 失效未读数缓存。
func (s *NotificationService) invalidateUnread(userID int64) {
	s.cache.Del(context.Background(), cache.NotifyUnreadKey(userID))
}

// MarkRead 标记单条已读，并失效未读数缓存。
func (s *NotificationService) MarkRead(id, userID int64) {
	s.notifs.MarkRead(id, userID)
	s.invalidateUnread(userID)
}

// MarkAllRead 全部已读，并失效未读数缓存。
func (s *NotificationService) MarkAllRead(userID int64) {
	s.notifs.MarkAllRead(userID)
	s.invalidateUnread(userID)
}

// UnreadCount 未读数：Redis cache-aside（未命中时回源 DB COUNT 并写回，带 TTL）。
func (s *NotificationService) UnreadCount(userID int64) int64 {
	ctx := context.Background()
	key := cache.NotifyUnreadKey(userID)
	if n, ok := s.cache.GetIntOK(ctx, key); ok {
		return n
	}
	n := s.notifs.UnreadCount(userID)
	s.cache.SetInt(ctx, key, n, unreadCacheTTL)
	return n
}
