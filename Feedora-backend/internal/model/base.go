package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 通用基础字段，供带软删除的模型内嵌复用。
type BaseModel struct {
	ID        int64          `gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
