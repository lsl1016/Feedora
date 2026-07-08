package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// Post 帖子完整信息。
type Post struct {
	PostID         int64          `json:"postId"`                  // 帖子 ID。
	PostType       string         `json:"postType"`                // 帖子类型。
	AuthorID       int64          `json:"authorId"`                // 作者 ID。
	Author         UserSummary    `json:"author"`                  // 作者信息。
	Title          string         `json:"title"`                   // 帖子标题。
	Content        string         `json:"content"`                 // Markdown 内容。
	Summary        string         `json:"summary"`                 // 摘要内容。
	Images         []string       `json:"images"`                  // 图片列表。
	Tags           []TagSummary   `json:"tags"`                    // 标签列表。
	Topics         []TopicSummary `json:"topics"`                  // 话题列表。
	CircleID       *int64         `json:"circleId,omitempty"`      // 所属圈子 ID。
	Circle         *CircleBrief   `json:"circle,omitempty"`        // 所属圈子信息。
	Visibility     string         `json:"visibility"`              // 可见性。
	Status         string         `json:"status"`                  // 帖子状态。
	IsTop          bool           `json:"isTop"`                   // 是否置顶。
	IsFeatured     bool           `json:"isFeatured"`              // 是否精选。
	IsSelected     bool           `json:"isSelected"`              // 是否已入选。
	ScheduledAt    string         `json:"scheduledAt,omitempty"`   // 预约发布时间。
	SourcePostID   *int64         `json:"sourcePostId,omitempty"`  // 来源帖子 ID。
	RepostComment  string         `json:"repostComment,omitempty"` // 转发文案。
	ViewCount      int64          `json:"viewCount"`               // 浏览数。
	LikeCount      int64          `json:"likeCount"`               // 点赞数。
	CommentCount   int64          `json:"commentCount"`            // 评论数。
	FavoriteCount  int64          `json:"favoriteCount"`           // 收藏数。
	ShareCount     int64          `json:"shareCount"`              // 分享数。
	RepostCount    int64          `json:"repostCount"`             // 转发数。
	HotScore       int64          `json:"hotScore"`                // 热度分。
	Liked          bool           `json:"liked"`                   // 当前用户是否点赞。
	Favorited      bool           `json:"favorited"`               // 当前用户是否收藏。
	FollowedAuthor bool           `json:"followedAuthor"`          // 当前用户是否关注作者。
	CreatedAt      string         `json:"createdAt"`               // 创建时间。
	PublishedAt    string         `json:"publishedAt,omitempty"`   // 发布时间。
	UpdatedAt      string         `json:"updatedAt"`               // 更新时间。
}

// CreatePostRequest 创建帖子请求。
type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required" example:"Feedora 使用体验"` // 帖子标题。
	Content     string   `json:"content" binding:"required" example:"这里是正文内容"`    // 帖子内容。
	Images      []string `json:"images" example:"/static/post-1.png"`             // 图片地址列表。
	TagIDs      []int64  `json:"tagIds" example:"1"`                              // 标签 ID 列表。
	TopicIDs    []int64  `json:"topicIds" example:"1"`                            // 话题 ID 列表。
	CircleID    *int64   `json:"circleId" example:"1"`                            // 圈子 ID。
	Visibility  string   `json:"visibility" example:"public"`                     // 可见范围。
	PublishMode string   `json:"publishMode" example:"immediate"`                 // 发布模式。
	ScheduledAt *string  `json:"scheduledAt" example:"2026-07-08T20:00:00Z"`      // 预约发布时间。
}

// UpdatePostRequest 编辑帖子请求。
type UpdatePostRequest struct {
	Title      *string  `json:"title" example:"更新后的标题"`              // 帖子标题。
	Content    *string  `json:"content" example:"更新后的正文内容"`          // 帖子内容。
	Visibility *string  `json:"visibility" example:"private"`        // 可见范围。
	Images     []string `json:"images" example:"/static/post-1.png"` // 图片地址列表。
}

// RepostPostRequest 转发帖子请求。
type RepostPostRequest struct {
	RepostComment string `json:"repostComment" example:"转发一下这篇内容"` // 转发附言。
}

// PostListQuery 帖子列表请求。
type PostListQuery struct {
	PageRequest
	FeedType      string `form:"feedType" example:"recommend"`  // 信息流类型。
	Sort          string `form:"sort" example:"latest"`         // 排序方式。
	Status        string `form:"status" example:"published"`    // 帖子状态。
	Keyword       string `form:"keyword" example:"feedora"`     // 搜索关键词。
	TagID         int64  `form:"tagId" example:"1"`             // 标签 ID。
	CircleID      int64  `form:"circleId" example:"1"`          // 圈子 ID。
	TopicID       int64  `form:"topicId" example:"1"`           // 话题 ID。
	AuthorID      int64  `form:"authorId" example:"1"`          // 作者 ID。
	IncludeHidden bool   `form:"includeHidden" example:"false"` // 是否包含隐藏帖子。
}

// PostListByTopicQuery 话题或标签下的帖子列表请求。
type PostListByTopicQuery struct {
	PageRequest
	Sort string `form:"sort" example:"latest"` // 排序方式。
}

// PostResponse 单个帖子响应。
type PostResponse struct {
	TraceEnvelope
	Data Post `json:"data"` // 业务数据。
}

// PostPageResponse 帖子分页响应。
type PostPageResponse struct {
	TraceEnvelope
	Data PageResult[Post] `json:"data"` // 业务数据。
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
