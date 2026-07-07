package model

import "time"

// Tag 标签表。
type Tag struct {
	ID          int64     `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name;size:64;uniqueIndex:uk_tag_name"`
	Description string    `gorm:"column:description;size:255"`
	Status      string    `gorm:"column:status;size:32;default:enabled;index"`
	UseCount    int64     `gorm:"column:use_count;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Tag) TableName() string { return "tags" }
