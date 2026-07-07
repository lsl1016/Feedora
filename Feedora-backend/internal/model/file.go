package model

import "time"

// File 文件表。
type File struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	UserID    int64     `gorm:"column:user_id;index"`
	BizType   string    `gorm:"column:biz_type;size:64;index"`
	Filename  string    `gorm:"column:filename;size:255"`
	ObjectKey string    `gorm:"column:object_key;size:512"`
	URL       string    `gorm:"column:url;size:512"`
	Size      int64     `gorm:"column:size;default:0"`
	MimeType  string    `gorm:"column:mime_type;size:128"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (File) TableName() string { return "files" }
