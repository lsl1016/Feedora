package repository

import (
	"errors"
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// PostRepository 帖子数据访问。
type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// PostFilter 帖子列表过滤条件（数据访问层）。
type PostFilter struct {
	FeedType      string
	Sort          string
	Status        string
	Keyword       string
	TagID         int64
	CircleID      int64
	TopicID       int64
	AuthorID      int64
	AuthorIDs     []int64 // 作者集合（关注动态等场景）。
	CircleIDs     []int64 // 圈子集合（圈子动态等场景）。
	TopicIDs      []int64 // 话题集合（话题动态等场景）。
	IncludeHidden bool
	ViewerID      int64
	Offset        int
	Limit         int
}

// List 按过滤条件分页查询帖子。
func (r *PostRepository) List(f PostFilter) ([]model.Post, int64, error) {
	q := r.db.Model(&model.Post{})

	switch {
	case f.Status == "all":
		if f.IncludeHidden && f.ViewerID > 0 {
			q = q.Where("status = ? OR (author_id = ? AND status <> ?)",
				model.PostPublished, f.ViewerID, model.PostDeleted)
		} else {
			q = q.Where("status = ?", model.PostPublished)
		}
	case f.Status != "":
		q = q.Where("status = ?", f.Status)
	default:
		q = q.Where("status = ?", model.PostPublished)
	}

	if f.AuthorID > 0 {
		q = q.Where("author_id = ?", f.AuthorID)
	}
	// 关注信息流：只看关注作者的内容，未登录时退化为普通信息流。
	if f.FeedType == "following" && f.ViewerID > 0 {
		q = q.Where("author_id IN (?)", r.db.Model(&model.UserFollow{}).Select("followee_id").Where("follower_id = ?", f.ViewerID))
	}
	if len(f.AuthorIDs) > 0 {
		q = q.Where("author_id IN ?", f.AuthorIDs)
	}
	// 指定圈子（单个 / 集合）时不叠加公开可见性限制，圈内帖对圈子成员可见。
	switch {
	case f.CircleID > 0:
		q = q.Where("circle_id = ?", f.CircleID)
	case len(f.CircleIDs) > 0:
		q = q.Where("circle_id IN ?", f.CircleIDs)
	default:
		q = q.Where("visibility = ? OR circle_id IS NULL", model.VisibilityPublic)
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		q = q.Where("title LIKE ? OR content_md LIKE ?", kw, kw)
	}
	if f.TagID > 0 {
		q = q.Where("id IN (?)", r.db.Model(&model.PostTag{}).Select("post_id").Where("tag_id = ?", f.TagID))
	}
	if f.TopicID > 0 {
		q = q.Where("id IN (?)", r.db.Model(&model.PostTopic{}).Select("post_id").Where("topic_id = ?", f.TopicID))
	}
	if len(f.TopicIDs) > 0 {
		q = q.Where("id IN (?)", r.db.Model(&model.PostTopic{}).Select("post_id").Where("topic_id IN ?", f.TopicIDs))
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	switch {
	case f.FeedType == "hot" || f.Sort == "hot":
		order = "hot_score DESC, created_at DESC"
	case f.Sort == "comment":
		order = "comment_count DESC, created_at DESC"
	case f.Sort == "favorite":
		order = "favorite_count DESC, created_at DESC"
	case f.Sort == "view":
		order = "view_count DESC, created_at DESC"
	}

	var rows []model.Post
	if err := q.Order(order).Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// FindByID 按 ID 查询帖子，不存在返回 (nil, nil)。
func (r *PostRepository) FindByID(id int64) (*model.Post, error) {
	var p model.Post
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindByIDs 批量查询已发布帖子（用于点赞 / 收藏列表）。
func (r *PostRepository) FindByIDs(ids []int64, onlyPublished bool) ([]model.Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	q := r.db.Where("id IN ?", ids)
	if onlyPublished {
		q = q.Where("status = ?", model.PostPublished)
	}
	var rows []model.Post
	err := q.Find(&rows).Error
	return rows, err
}

// CreateWithRelations 在事务中创建帖子及其图片、标签、话题关系并维护计数。
func (r *PostRepository) CreateWithRelations(p *model.Post, images []string, tagIDs, topicIDs []int64) error {
	now := time.Now()
	return tx(r.db, func(t *gorm.DB) error {
		if err := t.Create(p).Error; err != nil {
			return err
		}
		for idx, url := range images {
			if err := t.Create(&model.PostImage{PostID: p.ID, ImageURL: url, SortOrder: idx, CreatedAt: now}).Error; err != nil {
				return err
			}
		}
		for _, tid := range tagIDs {
			if err := t.Create(&model.PostTag{PostID: p.ID, TagID: tid, CreatedAt: now}).Error; err != nil {
				return err
			}
		}
		for _, tid := range topicIDs {
			if err := t.Create(&model.PostTopic{PostID: p.ID, TopicID: tid, CreatedAt: now}).Error; err != nil {
				return err
			}
		}
		// 计数只统计已发布帖子；草稿 / 定时帖在转正时由 PublishScheduled 补齐。
		if p.Status == model.PostPublished {
			for _, tid := range tagIDs {
				t.Model(&model.Tag{}).Where("id = ?", tid).UpdateColumn("use_count", gorm.Expr("use_count + 1"))
			}
			for _, tid := range topicIDs {
				t.Model(&model.Topic{}).Where("id = ?", tid).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
			}
			t.Model(&model.User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
			if p.CircleID != nil {
				t.Model(&model.Circle{}).Where("id = ?", *p.CircleID).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
			}
		}
		return nil
	})
}

// Create 创建单条帖子（用于转发）。
func (r *PostRepository) Create(p *model.Post) error {
	return r.db.Create(p).Error
}

// Update 按 ID 更新字段。
func (r *PostRepository) Update(id int64, updates map[string]any) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

// ReplaceImages 重置帖子图片。
func (r *PostRepository) ReplaceImages(postID int64, images []string) {
	r.db.Where("post_id = ?", postID).Delete(&model.PostImage{})
	now := time.Now()
	for idx, url := range images {
		r.db.Create(&model.PostImage{PostID: postID, ImageURL: url, SortOrder: idx, CreatedAt: now})
	}
}

// SoftDeleteWithCounters 在事务中软删除帖子；原状态为已发布时对称回减用户、话题、圈子、标签计数。
func (r *PostRepository) SoftDeleteWithCounters(p *model.Post) {
	_ = tx(r.db, func(t *gorm.DB) error {
		if err := t.Model(&model.Post{}).Where("id = ?", p.ID).Update("status", model.PostDeleted).Error; err != nil {
			return err
		}
		if err := t.Delete(&model.Post{}, p.ID).Error; err != nil {
			return err
		}
		if p.Status != model.PostPublished {
			return nil
		}
		t.Model(&model.User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("GREATEST(post_count - 1, 0)"))
		if p.CircleID != nil {
			t.Model(&model.Circle{}).Where("id = ?", *p.CircleID).UpdateColumn("post_count", gorm.Expr("GREATEST(post_count - 1, 0)"))
		}
		var topicIDs []int64
		t.Model(&model.PostTopic{}).Where("post_id = ?", p.ID).Pluck("topic_id", &topicIDs)
		for _, tid := range topicIDs {
			t.Model(&model.Topic{}).Where("id = ?", tid).UpdateColumn("post_count", gorm.Expr("GREATEST(post_count - 1, 0)"))
		}
		var tagIDs []int64
		t.Model(&model.PostTag{}).Where("post_id = ?", p.ID).Pluck("tag_id", &tagIDs)
		for _, tid := range tagIDs {
			t.Model(&model.Tag{}).Where("id = ?", tid).UpdateColumn("use_count", gorm.Expr("GREATEST(use_count - 1, 0)"))
		}
		return nil
	})
}

// FindScheduledDue 查询到点待转正的定时帖（未被软删）。
func (r *PostRepository) FindScheduledDue(now time.Time, limit int) ([]model.Post, error) {
	var rows []model.Post
	err := r.db.Where("status = ? AND scheduled_at IS NOT NULL AND scheduled_at <= ?", model.PostScheduled, now).
		Order("scheduled_at ASC").Limit(limit).Find(&rows).Error
	return rows, err
}

// PublishScheduled 事务内把定时帖转正为已发布，并补齐发布计数（用户、圈子、话题、标签）。
// 返回是否实际转正（并发下已被处理时返回 false, nil）。
func (r *PostRepository) PublishScheduled(p *model.Post) (bool, error) {
	now := time.Now()
	published := false
	err := tx(r.db, func(t *gorm.DB) error {
		res := t.Model(&model.Post{}).Where("id = ? AND status = ?", p.ID, model.PostScheduled).
			Updates(map[string]any{"status": model.PostPublished, "published_at": now, "updated_at": now})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}
		published = true
		t.Model(&model.User{}).Where("id = ?", p.AuthorID).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
		if p.CircleID != nil {
			t.Model(&model.Circle{}).Where("id = ?", *p.CircleID).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
		}
		var topicIDs []int64
		t.Model(&model.PostTopic{}).Where("post_id = ?", p.ID).Pluck("topic_id", &topicIDs)
		for _, tid := range topicIDs {
			t.Model(&model.Topic{}).Where("id = ?", tid).UpdateColumn("post_count", gorm.Expr("post_count + 1"))
		}
		var tagIDs []int64
		t.Model(&model.PostTag{}).Where("post_id = ?", p.ID).Pluck("tag_id", &tagIDs)
		for _, tid := range tagIDs {
			t.Model(&model.Tag{}).Where("id = ?", tid).UpdateColumn("use_count", gorm.Expr("use_count + 1"))
		}
		return nil
	})
	return published, err
}

// IncColumn 对计数列做增量（可为负）。
func (r *PostRepository) IncColumn(id int64, column string, delta int) {
	if delta >= 0 {
		r.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(column, gorm.Expr(column+" + ?", delta))
	} else {
		r.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn(column, gorm.Expr("GREATEST("+column+" - ?, 0)", -delta))
	}
}

// IncViewCount 递增浏览量。
func (r *PostRepository) IncViewCount(id int64) {
	r.db.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
}

// Exists 判断帖子是否存在。
func (r *PostRepository) Exists(id int64) bool {
	var count int64
	r.db.Model(&model.Post{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// ImagesByPostIDs 批量查询帖子图片，返回 postID->图片URL列表。
func (r *PostRepository) ImagesByPostIDs(postIDs []int64) map[int64][]string {
	res := map[int64][]string{}
	if len(postIDs) == 0 {
		return res
	}
	var imgs []model.PostImage
	r.db.Where("post_id IN ?", postIDs).Order("sort_order ASC, id ASC").Find(&imgs)
	for _, im := range imgs {
		res[im.PostID] = append(res[im.PostID], im.ImageURL)
	}
	return res
}
