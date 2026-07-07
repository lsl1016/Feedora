package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// Post 帖子完整信息。
type Post struct {
	PostID         int64          `json:"postId"`
	PostType       string         `json:"postType"`
	AuthorID       int64          `json:"authorId"`
	Author         UserSummary    `json:"author"`
	Title          string         `json:"title"`
	Content        string         `json:"content"`
	Summary        string         `json:"summary"`
	Images         []string       `json:"images"`
	Tags           []TagSummary   `json:"tags"`
	Topics         []TopicSummary `json:"topics"`
	CircleID       *int64         `json:"circleId,omitempty"`
	Circle         *CircleBrief   `json:"circle,omitempty"`
	Visibility     string         `json:"visibility"`
	Status         string         `json:"status"`
	IsTop          bool           `json:"isTop"`
	IsFeatured     bool           `json:"isFeatured"`
	IsSelected     bool           `json:"isSelected"`
	ScheduledAt    string         `json:"scheduledAt,omitempty"`
	SourcePostID   *int64         `json:"sourcePostId,omitempty"`
	RepostComment  string         `json:"repostComment,omitempty"`
	ViewCount      int64          `json:"viewCount"`
	LikeCount      int64          `json:"likeCount"`
	CommentCount   int64          `json:"commentCount"`
	FavoriteCount  int64          `json:"favoriteCount"`
	ShareCount     int64          `json:"shareCount"`
	RepostCount    int64          `json:"repostCount"`
	HotScore       int64          `json:"hotScore"`
	Liked          bool           `json:"liked"`
	Favorited      bool           `json:"favorited"`
	FollowedAuthor bool           `json:"followedAuthor"`
	CreatedAt      string         `json:"createdAt"`
	PublishedAt    string         `json:"publishedAt,omitempty"`
	UpdatedAt      string         `json:"updatedAt"`
}

// CreatePostRequest 创建帖子请求。
type CreatePostRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Images      []string `json:"images"`
	TagIDs      []int64  `json:"tagIds"`
	TopicIDs    []int64  `json:"topicIds"`
	CircleID    *int64   `json:"circleId"`
	Visibility  string   `json:"visibility"`
	PublishMode string   `json:"publishMode"`
	ScheduledAt *string  `json:"scheduledAt"`
}

// UpdatePostRequest 编辑帖子请求。
type UpdatePostRequest struct {
	Title      *string  `json:"title"`
	Content    *string  `json:"content"`
	Visibility *string  `json:"visibility"`
	Images     []string `json:"images"`
}

// PostRelations 组装帖子 DTO 所需的关联数据。
type PostRelations struct {
	Author    *model.User
	Images    []string
	Tags      []TagSummary
	Topics    []TopicSummary
	Circle    *model.Circle
	Liked     bool
	Favorited bool
}

// ToPost 转换帖子完整信息。
func ToPost(p *model.Post, r PostRelations) Post {
	tags := r.Tags
	if tags == nil {
		tags = []TagSummary{}
	}
	topics := r.Topics
	if topics == nil {
		topics = []TopicSummary{}
	}
	images := r.Images
	if images == nil {
		images = []string{}
	}
	return Post{
		PostID:        p.ID,
		PostType:      p.PostType,
		AuthorID:      p.AuthorID,
		Author:        ToUserSummary(r.Author),
		Title:         p.Title,
		Content:       p.ContentMD,
		Summary:       p.Summary,
		Images:        images,
		Tags:          tags,
		Topics:        topics,
		CircleID:      p.CircleID,
		Circle:        ToCircleBrief(r.Circle),
		Visibility:    p.Visibility,
		Status:        p.Status,
		ScheduledAt:   utils.FormatTimePtr(p.ScheduledAt),
		SourcePostID:  p.SourcePostID,
		RepostComment: p.RepostComment,
		ViewCount:     p.ViewCount,
		LikeCount:     p.LikeCount,
		CommentCount:  p.CommentCount,
		FavoriteCount: p.FavoriteCount,
		ShareCount:    p.ShareCount,
		RepostCount:   p.RepostCount,
		HotScore:      p.HotScore,
		Liked:         r.Liked,
		Favorited:     r.Favorited,
		CreatedAt:     utils.FormatTime(p.CreatedAt),
		PublishedAt:   utils.FormatTimePtr(p.PublishedAt),
		UpdatedAt:     utils.FormatTime(p.UpdatedAt),
	}
}
