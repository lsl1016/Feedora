package repository

import (
	"errors"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// TagRepository 标签数据访问。
type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// ListEnabled 查询启用中的标签，按使用次数倒序。
func (r *TagRepository) ListEnabled() ([]model.Tag, error) {
	var rows []model.Tag
	err := r.db.Where("status = ?", model.StatusEnabled).Order("use_count DESC, id ASC").Find(&rows).Error
	return rows, err
}

// ListAll 查询全部标签（后台）。
func (r *TagRepository) ListAll() ([]model.Tag, error) {
	var rows []model.Tag
	err := r.db.Order("id ASC").Find(&rows).Error
	return rows, err
}

// FindByID 按 ID 查询标签，不存在返回 (nil, nil)。
func (r *TagRepository) FindByID(id int64) (*model.Tag, error) {
	var t model.Tag
	err := r.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create 创建标签。
func (r *TagRepository) Create(t *model.Tag) error {
	return r.db.Create(t).Error
}

// Update 按 ID 更新标签。
func (r *TagRepository) Update(id int64, updates map[string]any) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).Updates(updates).Error
}

// FindByPostIDs 批量查询帖子标签，返回 postID->标签列表。
func (r *TagRepository) FindByPostIDs(postIDs []int64) map[int64][]model.Tag {
	res := map[int64][]model.Tag{}
	if len(postIDs) == 0 {
		return res
	}
	type row struct {
		PostID int64
		ID     int64
		Name   string
	}
	var rs []row
	r.db.Table("post_tags").
		Select("post_tags.post_id as post_id, tags.id as id, tags.name as name").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id IN ?", postIDs).Scan(&rs)
	for _, x := range rs {
		res[x.PostID] = append(res[x.PostID], model.Tag{ID: x.ID, Name: x.Name})
	}
	return res
}
