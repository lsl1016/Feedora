package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// TagSummary 标签简要信息。
type TagSummary struct {
	TagID   int64  `json:"tagId"`
	TagName string `json:"tagName"`
}

// ContentTag 标签完整信息。
type ContentTag struct {
	TagID       int64  `json:"tagId"`
	TagName     string `json:"tagName"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UseCount    int64  `json:"useCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
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
