package repository

import (
	"errors"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// TopicRepository 话题数据访问。
type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

// List 话题广场，tab 支持 official/hot/latest。
func (r *TopicRepository) List(tab string, offset, limit int) ([]model.Topic, int64, error) {
	q := r.db.Model(&model.Topic{}).Where("status = ?", model.StatusEnabled)
	order := "created_at DESC"
	switch tab {
	case "official":
		q = q.Where("is_official = ?", true)
	case "hot":
		order = "participant_count DESC, post_count DESC"
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Topic
	if err := q.Order(order).Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListAll 查询全部话题（后台分页）。
func (r *TopicRepository) ListAll(offset, limit int) ([]model.Topic, int64, error) {
	var total int64
	r.db.Model(&model.Topic{}).Count(&total)
	var rows []model.Topic
	err := r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// FindByID 按 ID 查询话题，不存在返回 (nil, nil)。
func (r *TopicRepository) FindByID(id int64) (*model.Topic, error) {
	var t model.Topic
	err := r.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create 创建话题。
func (r *TopicRepository) Create(t *model.Topic) error {
	return r.db.Create(t).Error
}

// Update 按 ID 更新话题。
func (r *TopicRepository) Update(id int64, updates map[string]any) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", id).Updates(updates).Error
}

// FindByPostIDs 批量查询帖子话题，返回 postID->话题列表。
func (r *TopicRepository) FindByPostIDs(postIDs []int64) map[int64][]model.Topic {
	res := map[int64][]model.Topic{}
	if len(postIDs) == 0 {
		return res
	}
	type row struct {
		PostID int64
		ID     int64
		Name   string
	}
	var rs []row
	r.db.Table("post_topics").
		Select("post_topics.post_id as post_id, topics.id as id, topics.name as name").
		Joins("JOIN topics ON topics.id = post_topics.topic_id").
		Where("post_topics.post_id IN ?", postIDs).Scan(&rs)
	for _, x := range rs {
		res[x.PostID] = append(res[x.PostID], model.Topic{ID: x.ID, Name: x.Name})
	}
	return res
}

// RecomputeParticipantCountsByPost 按帖子重算其所属话题的参与人数
// （统计未删除帖子的去重作者数），幂等，供 worker 在帖子创建 / 删除后调用。
func (r *TopicRepository) RecomputeParticipantCountsByPost(postID int64) {
	sql := `
UPDATE topics t
SET participant_count = (
	SELECT COUNT(DISTINCT p.author_id)
	FROM post_topics pt
	JOIN posts p ON p.id = pt.post_id
		AND p.deleted_at IS NULL
		AND p.status <> 'deleted'
	WHERE pt.topic_id = t.id
)
WHERE t.id IN (SELECT topic_id FROM post_topics WHERE post_id = ?)`
	r.db.Exec(sql, postID)
}
