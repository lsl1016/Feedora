package model

import "time"

// Topic 话题表。
type Topic struct {
	ID               int64     `gorm:"primaryKey;column:id"`
	Name             string    `gorm:"column:name;size:128;uniqueIndex:uk_topic_name"`
	Description      string    `gorm:"column:description;size:500"`
	CoverURL         string    `gorm:"column:cover_url;size:512"`
	IsOfficial       bool      `gorm:"column:is_official;default:false"`
	IsRecommended    bool      `gorm:"column:is_recommended;default:false;index"`
	ParticipantCount int64     `gorm:"column:participant_count;default:0"`
	PostCount        int64     `gorm:"column:post_count;default:0"`
	Status           string    `gorm:"column:status;size:32;default:enabled;index"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (Topic) TableName() string { return "topics" }
