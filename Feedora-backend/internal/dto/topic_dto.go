package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// TopicSummary 话题简要信息。
type TopicSummary struct {
	TopicID int64  `json:"topicId"`
	Name    string `json:"name"`
}

// Topic 话题完整信息。
type Topic struct {
	TopicID          int64  `json:"topicId"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	CoverImage       string `json:"coverImage"`
	PostCount        int64  `json:"postCount"`
	ParticipantCount int64  `json:"participantCount"`
	IsOfficial       bool   `json:"isOfficial"`
	IsRecommended    bool   `json:"isRecommended"`
	Status           string `json:"status"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
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
