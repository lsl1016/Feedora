package repository

import (
	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// FileRepository 文件记录数据访问。
type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create 写入文件记录。
func (r *FileRepository) Create(f *model.File) error {
	return r.db.Create(f).Error
}
