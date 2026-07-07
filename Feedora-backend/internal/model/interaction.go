package model

import "time"

// PostLike 帖子点赞表。
type PostLike struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	PostID    int64     `gorm:"column:post_id;uniqueIndex:uk_post_like_user"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex:uk_post_like_user;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PostLike) TableName() string { return "post_likes" }

// PostFavorite 帖子收藏表。
type PostFavorite struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	PostID    int64     `gorm:"column:post_id;uniqueIndex:uk_post_fav_user"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex:uk_post_fav_user;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PostFavorite) TableName() string { return "post_favorites" }

// CommentLike 评论点赞表。
type CommentLike struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	CommentID int64     `gorm:"column:comment_id;uniqueIndex:uk_comment_like_user"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex:uk_comment_like_user;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (CommentLike) TableName() string { return "comment_likes" }
