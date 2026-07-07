package model

import (
	"time"

	"gorm.io/gorm"
)

// Post 帖子表。
type Post struct {
	ID            int64          `gorm:"primaryKey;column:id"`
	AuthorID      int64          `gorm:"column:author_id;index:idx_author_status"`
	Title         string         `gorm:"column:title;size:200"`
	ContentMD     string         `gorm:"column:content_md;type:longtext"`
	Summary       string         `gorm:"column:summary;size:500"`
	PostType      string         `gorm:"column:post_type;size:32;default:original"`
	SourcePostID  *int64         `gorm:"column:source_post_id"`
	RepostComment string         `gorm:"column:repost_comment;size:1000"`
	CircleID      *int64         `gorm:"column:circle_id;index:idx_circle_status"`
	Visibility    string         `gorm:"column:visibility;size:32;default:public"`
	Status        string         `gorm:"column:status;size:32;default:published;index:idx_author_status;index:idx_circle_status"`
	CoverURL      string         `gorm:"column:cover_url;size:512"`
	ViewCount     int64          `gorm:"column:view_count;default:0"`
	LikeCount     int64          `gorm:"column:like_count;default:0"`
	CommentCount  int64          `gorm:"column:comment_count;default:0"`
	FavoriteCount int64          `gorm:"column:favorite_count;default:0"`
	ShareCount    int64          `gorm:"column:share_count;default:0"`
	RepostCount   int64          `gorm:"column:repost_count;default:0"`
	HotScore      int64          `gorm:"column:hot_score;default:0;index"`
	ScheduledAt   *time.Time     `gorm:"column:scheduled_at"`
	PublishedAt   *time.Time     `gorm:"column:published_at"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Post) TableName() string { return "posts" }

// PostImage 帖子图片表。
type PostImage struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	PostID    int64     `gorm:"column:post_id;index"`
	ImageURL  string    `gorm:"column:image_url;size:512"`
	SortOrder int       `gorm:"column:sort_order;default:0"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PostImage) TableName() string { return "post_images" }

// PostTag 帖子标签关系表。
type PostTag struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	PostID    int64     `gorm:"column:post_id;uniqueIndex:uk_post_tag"`
	TagID     int64     `gorm:"column:tag_id;uniqueIndex:uk_post_tag;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PostTag) TableName() string { return "post_tags" }

// PostTopic 帖子话题关系表。
type PostTopic struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	PostID    int64     `gorm:"column:post_id;uniqueIndex:uk_post_topic"`
	TopicID   int64     `gorm:"column:topic_id;uniqueIndex:uk_post_topic;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (PostTopic) TableName() string { return "post_topics" }
