package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// NotificationItem 通知列表项（对应前端 NotificationItem）。
type NotificationItem struct {
	NotificationID int64  `json:"notificationId"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Category       string `json:"category"`
	TargetURL      string `json:"targetUrl,omitempty"`
	ReadStatus     string `json:"readStatus"`
	CreatedAt      string `json:"createdAt"`
}

// ToNotificationItem 转换通知。type 字段即前端 category。
func ToNotificationItem(n *model.Notification) NotificationItem {
	return NotificationItem{
		NotificationID: n.ID,
		Title:          n.Title,
		Content:        n.Content,
		Category:       n.Type,
		TargetURL:      n.TargetURL,
		ReadStatus:     n.ReadStatus,
		CreatedAt:      utils.FormatTime(n.CreatedAt),
	}
}
