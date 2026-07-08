package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// TopicListQuery 话题列表请求。
type TopicListQuery struct {
	PageRequest
	Tab string `form:"tab" example:"all"` // 分类标签。
}

// TopicPostsQuery 话题帖子列表请求。
type TopicPostsQuery struct {
	PageRequest
	Sort string `form:"sort" example:"latest"` // 排序方式。
}

// TopicSummary 话题简要信息。
type TopicSummary struct {
	TopicID int64  `json:"topicId"` // 话题 ID。
	Name    string `json:"name"`    // 话题名称。
}

// Topic 话题完整信息。
type Topic struct {
	TopicID          int64  `json:"topicId"`          // 话题 ID。
	Name             string `json:"name"`             // 话题名称。
	Description      string `json:"description"`      // 话题描述。
	CoverImage       string `json:"coverImage"`       // 封面图片。
	PostCount        int64  `json:"postCount"`        // 帖子数量。
	ParticipantCount int64  `json:"participantCount"` // 参与人数。
	IsOfficial       bool   `json:"isOfficial"`       // 是否官方。
	IsRecommended    bool   `json:"isRecommended"`    // 是否推荐。
	Status           string `json:"status"`           // 话题状态。
	CreatedAt        string `json:"createdAt"`        // 创建时间。
	UpdatedAt        string `json:"updatedAt"`        // 更新时间。
}

// TopicResponse 单个话题响应。
type TopicResponse struct {
	TraceEnvelope
	Data Topic `json:"data"` // 业务数据。
}

// TopicPageResponse 话题分页响应。
type TopicPageResponse struct {
	TraceEnvelope
	Data PageResult[Topic] `json:"data"` // 业务数据。
}

// ToTopic 转换话题完整信息。
func ToTopic(t *model.Topic) Topic {
	return Topic{
		TopicID:          t.ID,
		Name:             t.Name,
		Description:      t.Description,
		CoverImage:       t.CoverURL,
		PostCount:        t.PostCount,
		ParticipantCount: t.ParticipantCount,
		IsOfficial:       t.IsOfficial,
		IsRecommended:    t.IsRecommended,
		Status:           t.Status,
		CreatedAt:        utils.FormatTime(t.CreatedAt),
		UpdatedAt:        utils.FormatTime(t.UpdatedAt),
	}
}

// ToTopicSummary 转换话题简要信息。
func ToTopicSummary(t *model.Topic) TopicSummary {
	return TopicSummary{TopicID: t.ID, Name: t.Name}
}
