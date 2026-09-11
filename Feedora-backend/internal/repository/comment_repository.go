package repository

import (
	"errors"
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问。
type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// FindByID 按 ID 查询评论，不存在返回 (nil, nil)。
func (r *CommentRepository) FindByID(id int64) (*model.Comment, error) {
	var m model.Comment
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByPost 查询帖子的可见评论（含删除占位节点），按创建时间升序。
func (r *CommentRepository) ListByPost(postID int64) ([]model.Comment, error) {
	var rows []model.Comment
	err := r.db.Where("post_id = ? AND status IN ?", postID, []string{model.CommentNormal, model.CommentDeleted}).
		Order("created_at ASC").Find(&rows).Error
	return rows, err
}

// ListMine 分页查询用户的评论。
func (r *CommentRepository) ListMine(userID int64, offset, limit int) ([]model.Comment, int64, error) {
	q := r.db.Model(&model.Comment{}).Where("user_id = ? AND status = ?", userID, model.CommentNormal)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Comment
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// CreateWithCounters 在事务中创建评论并递增帖子、用户评论数。
func (r *CommentRepository) CreateWithCounters(c *model.Comment) error {
	return tx(r.db, func(t *gorm.DB) error {
		if err := t.Create(c).Error; err != nil {
			return err
		}
		t.Model(&model.Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
		t.Model(&model.User{}).Where("id = ?", c.UserID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))
		return nil
	})
}

// HasChildren 判断评论是否仍有未物理删除的子回复。
func (r *CommentRepository) HasChildren(id int64) bool {
	var n int64
	r.db.Model(&model.Comment{}).Where("parent_id = ?", id).Count(&n)
	return n > 0
}

// DeleteWithCounters 在事务中删除评论，并递减帖子、用户评论数。
// placeholder 为 true 时（有子回复的评论）保留行作为占位节点，仅置 status，不软删。
func (r *CommentRepository) DeleteWithCounters(c *model.Comment, placeholder bool) {
	_ = tx(r.db, func(t *gorm.DB) error {
		if placeholder {
			if err := t.Model(&model.Comment{}).Where("id = ?", c.ID).
				Updates(map[string]any{"status": model.CommentDeleted, "updated_at": time.Now()}).Error; err != nil {
				return err
			}
		} else {
			if err := t.Model(&model.Comment{}).Where("id = ?", c.ID).Update("status", model.CommentDeleted).Error; err != nil {
				return err
			}
			if err := t.Delete(&model.Comment{}, c.ID).Error; err != nil {
				return err
			}
		}
		t.Model(&model.Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
		t.Model(&model.User{}).Where("id = ?", c.UserID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
		return nil
	})
}

// IncLikeCount 对评论点赞数做增量（可为负）。
func (r *CommentRepository) IncLikeCount(id int64, delta int) {
	if delta >= 0 {
		r.db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta))
	} else {
		r.db.Model(&model.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - ?, 0)", -delta))
	}
}

// Count 统计评论总数（后台）。
func (r *CommentRepository) Count() int64 {
	var total int64
	r.db.Model(&model.Comment{}).Count(&total)
	return total
}

// ListAll 分页查询全部评论（后台）。
func (r *CommentRepository) ListAll(offset, limit int) []model.Comment {
	var rows []model.Comment
	r.db.Order("id DESC").Offset(offset).Limit(limit).Find(&rows)
	return rows
}
