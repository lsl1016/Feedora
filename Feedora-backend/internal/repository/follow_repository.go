package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// FollowRepository 关注关系数据访问。
type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Follow 建立关注关系，已存在时返回 false（幂等）。
func (r *FollowRepository) Follow(followerID, followeeID int64) (bool, error) {
	f := &model.UserFollow{FollowerID: followerID, FolloweeID: followeeID, CreatedAt: time.Now()}
	if err := r.db.Create(f).Error; err != nil {
		if isDuplicateEntry(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Unfollow 取消关注，关系不存在时返回 false。
func (r *FollowRepository) Unfollow(followerID, followeeID int64) (bool, error) {
	res := r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&model.UserFollow{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// IsFollowing 查询 follower 是否关注了 followee。
func (r *FollowRepository) IsFollowing(followerID, followeeID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserFollow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count).Error
	return count > 0, err
}

// ListFollowing 关注的人列表。
func (r *FollowRepository) ListFollowing(followerID int64, offset, limit int) ([]model.User, int64, error) {
	sub := r.db.Model(&model.UserFollow{}).Select("followee_id").Where("follower_id = ?", followerID)
	return r.pageUsers(sub, offset, limit)
}

// ListFollowers 粉丝列表。
func (r *FollowRepository) ListFollowers(followeeID int64, offset, limit int) ([]model.User, int64, error) {
	sub := r.db.Model(&model.UserFollow{}).Select("follower_id").Where("followee_id = ?", followeeID)
	return r.pageUsers(sub, offset, limit)
}

func (r *FollowRepository) pageUsers(idSub *gorm.DB, offset, limit int) ([]model.User, int64, error) {
	var total int64
	if err := r.db.Model(&model.User{}).Where("id IN (?) AND status = ?", idSub, model.UserNormal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []model.User
	err := r.db.Where("id IN (?) AND status = ?", idSub, model.UserNormal).
		Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// JoinedCircles 用户加入的圈子（关注中心「关注的圈子」数据源）。
func (r *FollowRepository) JoinedCircles(userID int64, offset, limit int) ([]model.Circle, int64, error) {
	sub := r.db.Model(&model.CircleMember{}).Select("circle_id").
		Where("user_id = ? AND status = ?", userID, model.CircleMemberNormal)
	var total int64
	if err := r.db.Model(&model.Circle{}).Where("id IN (?) AND status = ?", sub, model.CircleMemberNormal).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var circles []model.Circle
	err := r.db.Where("id IN (?) AND status = ?", sub, model.CircleMemberNormal).
		Order("id DESC").Offset(offset).Limit(limit).Find(&circles).Error
	return circles, total, err
}

// ParticipatedTopics 用户发帖涉及的话题（「关注的话题」数据源）。
func (r *FollowRepository) ParticipatedTopics(userID int64, offset, limit int) ([]model.Topic, int64, error) {
	sub := r.db.Model(&model.PostTopic{}).Select("topic_id").Where(
		"post_id IN (?)", r.db.Model(&model.Post{}).Select("id").Where("author_id = ? AND deleted_at IS NULL", userID),
	)
	var total int64
	if err := r.db.Model(&model.Topic{}).Where("id IN (?) AND status = ?", sub, model.StatusEnabled).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var topics []model.Topic
	err := r.db.Where("id IN (?) AND status = ?", sub, model.StatusEnabled).
		Order("id DESC").Offset(offset).Limit(limit).Find(&topics).Error
	return topics, total, err
}

// UsedTags 用户发帖使用的标签（「关注的标签」数据源）。
func (r *FollowRepository) UsedTags(userID int64, offset, limit int) ([]model.Tag, int64, error) {
	sub := r.db.Model(&model.PostTag{}).Select("tag_id").Where(
		"post_id IN (?)", r.db.Model(&model.Post{}).Select("id").Where("author_id = ? AND deleted_at IS NULL", userID),
	)
	var total int64
	if err := r.db.Model(&model.Tag{}).Where("id IN (?) AND status = ?", sub, model.StatusEnabled).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var tags []model.Tag
	err := r.db.Where("id IN (?) AND status = ?", sub, model.StatusEnabled).
		Order("use_count DESC, id DESC").Offset(offset).Limit(limit).Find(&tags).Error
	return tags, total, err
}

// isDuplicateEntry 判断 MySQL 唯一键冲突（不同驱动错误形态不同，双保险）。
func isDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "Error 1062")
}
