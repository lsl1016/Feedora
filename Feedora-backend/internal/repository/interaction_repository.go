package repository

import (
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// InteractionRepository 点赞 / 收藏数据访问。
type InteractionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

// AddPostLike 新增点赞，返回是否实际插入（幂等）。
func (r *InteractionRepository) AddPostLike(postID, userID int64) bool {
	res := r.db.Create(&model.PostLike{PostID: postID, UserID: userID, CreatedAt: time.Now()})
	return res.Error == nil && res.RowsAffected > 0
}

// RemovePostLike 取消点赞，返回是否实际删除。
func (r *InteractionRepository) RemovePostLike(postID, userID int64) bool {
	res := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostLike{})
	return res.RowsAffected > 0
}

// AddPostFavorite 新增收藏，返回是否实际插入。
func (r *InteractionRepository) AddPostFavorite(postID, userID int64) bool {
	res := r.db.Create(&model.PostFavorite{PostID: postID, UserID: userID, CreatedAt: time.Now()})
	return res.Error == nil && res.RowsAffected > 0
}

// RemovePostFavorite 取消收藏，返回是否实际删除。
func (r *InteractionRepository) RemovePostFavorite(postID, userID int64) bool {
	res := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostFavorite{})
	return res.RowsAffected > 0
}

// HasLiked 判断用户是否已点赞。
func (r *InteractionRepository) HasLiked(postID, userID int64) bool {
	var count int64
	r.db.Model(&model.PostLike{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count)
	return count > 0
}

// HasFavorited 判断用户是否已收藏。
func (r *InteractionRepository) HasFavorited(postID, userID int64) bool {
	var count int64
	r.db.Model(&model.PostFavorite{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count)
	return count > 0
}

// LikedSet 返回用户在给定帖子集合中已点赞的帖子。
func (r *InteractionRepository) LikedSet(userID int64, postIDs []int64) map[int64]bool {
	res := map[int64]bool{}
	if userID <= 0 || len(postIDs) == 0 {
		return res
	}
	var ls []model.PostLike
	r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&ls)
	for _, l := range ls {
		res[l.PostID] = true
	}
	return res
}

// FavoritedSet 返回用户在给定帖子集合中已收藏的帖子。
func (r *InteractionRepository) FavoritedSet(userID int64, postIDs []int64) map[int64]bool {
	res := map[int64]bool{}
	if userID <= 0 || len(postIDs) == 0 {
		return res
	}
	var fs []model.PostFavorite
	r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&fs)
	for _, f := range fs {
		res[f.PostID] = true
	}
	return res
}

// PagePostIDsByUser 分页查询用户点赞 / 收藏的帖子 ID（按关系时间倒序），避免全量加载后内存分页。
func (r *InteractionRepository) PagePostIDsByUser(table string, userID int64, offset, limit int) ([]int64, int64) {
	var total int64
	r.db.Table(table).Where("user_id = ?", userID).Count(&total)
	if total == 0 {
		return nil, 0
	}
	var ids []int64
	r.db.Table(table).Where("user_id = ?", userID).Order("id DESC").
		Offset(offset).Limit(limit).Pluck("post_id", &ids)
	return ids, total
}

// AddCommentLike 评论点赞，返回是否实际插入。
func (r *InteractionRepository) AddCommentLike(commentID, userID int64) bool {
	res := r.db.Create(&model.CommentLike{CommentID: commentID, UserID: userID, CreatedAt: time.Now()})
	return res.Error == nil && res.RowsAffected > 0
}

// RemoveCommentLike 取消评论点赞，返回是否实际删除。
func (r *InteractionRepository) RemoveCommentLike(commentID, userID int64) bool {
	res := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&model.CommentLike{})
	return res.RowsAffected > 0
}

// CommentLikedSet 返回用户在给定评论集合中已点赞的评论。
func (r *InteractionRepository) CommentLikedSet(userID int64, commentIDs []int64) map[int64]bool {
	res := map[int64]bool{}
	if userID <= 0 || len(commentIDs) == 0 {
		return res
	}
	var ls []model.CommentLike
	r.db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&ls)
	for _, l := range ls {
		res[l.CommentID] = true
	}
	return res
}
