package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// TagPostsQuery 标签帖子列表请求。
type TagPostsQuery struct {
	PageRequest
	Sort string `form:"sort" example:"latest"` // 排序方式。
}

// TagSummary 标签简要信息。
type TagSummary struct {
	TagID   int64  `json:"tagId"`   // 标签 ID。
	TagName string `json:"tagName"` // 标签名称。
}

// ContentTag 标签完整信息。
type ContentTag struct {
	TagID       int64  `json:"tagId"`       // 标签 ID。
	TagName     string `json:"tagName"`     // 标签名称。
	Description string `json:"description"` // 标签描述。
	Status      string `json:"status"`      // 标签状态。
	UseCount    int64  `json:"useCount"`    // 使用次数。
	CreatedAt   string `json:"createdAt"`   // 创建时间。
	UpdatedAt   string `json:"updatedAt"`   // 更新时间。
}

// TagListResponse 标签列表响应。
type TagListResponse struct {
	TraceEnvelope
	Data []ContentTag `json:"data"` // 业务数据。
}

// TagResponse 单个标签响应。
type TagResponse struct {
	TraceEnvelope
	Data ContentTag `json:"data"` // 业务数据。
}

// ToContentTag 转换标签完整信息。
func ToContentTag(t *model.Tag) ContentTag {
	return ContentTag{
		TagID:       t.ID,
		TagName:     t.Name,
		Description: t.Description,
		Status:      t.Status,
		UseCount:    t.UseCount,
		CreatedAt:   utils.FormatTime(t.CreatedAt),
		UpdatedAt:   utils.FormatTime(t.UpdatedAt),
	}
}

// ToTagSummary 转换标签简要信息。
func ToTagSummary(t *model.Tag) TagSummary {
	return TagSummary{TagID: t.ID, TagName: t.Name}
}
