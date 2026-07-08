package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// NotificationListQuery 通知列表请求。
type NotificationListQuery struct {
	PageRequest
	Category string `form:"category" example:"system"` // 通知分类。
}

// NotificationItem 通知列表项（对应前端 NotificationItem）。
type NotificationItem struct {
	NotificationID int64  `json:"notificationId"`      // 通知 ID。
	Title          string `json:"title"`               // 通知标题。
	Content        string `json:"content"`             // 通知内容。
	Category       string `json:"category"`            // 通知分类。
	TargetURL      string `json:"targetUrl,omitempty"` // 跳转链接。
	ReadStatus     string `json:"readStatus"`          // 已读状态。
	CreatedAt      string `json:"createdAt"`           // 创建时间。
}

// NotificationListResponse 通知列表响应。
type NotificationListResponse struct {
	TraceEnvelope
	Data []NotificationItem `json:"data"` // 业务数据。
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
