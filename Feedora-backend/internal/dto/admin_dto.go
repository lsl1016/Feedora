package dto

// AdminUserListQuery 后台用户列表请求。
type AdminUserListQuery struct {
	PageRequest
	Keyword string `form:"keyword" example:"feedora"` // 搜索关键词。
}

// AdminPostListQuery 后台帖子列表请求。
type AdminPostListQuery struct {
	PageRequest
	Status string `form:"status" example:"published"` // 帖子状态。
}

// AdminCommentListQuery 后台评论列表请求。
type AdminCommentListQuery struct {
	PageRequest
}

// AdminTopicListQuery 后台话题列表请求。
type AdminTopicListQuery struct {
	PageRequest
}

// AdminCircleListQuery 后台圈子列表请求。
type AdminCircleListQuery struct {
	PageRequest
}

// AdminLogListQuery 后台操作日志列表请求。
type AdminLogListQuery struct {
	PageRequest
}

// AdminStats 后台概览统计。
type AdminStats struct {
	UserCount    int64 `json:"userCount"`    // 用户数量。
	PostCount    int64 `json:"postCount"`    // 帖子数量。
	CommentCount int64 `json:"commentCount"` // 评论数量。
	CircleCount  int64 `json:"circleCount"`  // 圈子数量。
}

// AdminPostItem 后台帖子列表项。
type AdminPostItem struct {
	ID            int64   `json:"ID"`            // 帖子主键 ID。
	AuthorID      int64   `json:"AuthorID"`      // 作者 ID。
	Title         string  `json:"Title"`         // 帖子标题。
	ContentMD     string  `json:"ContentMD"`     // Markdown 内容。
	Summary       string  `json:"Summary"`       // 摘要内容。
	PostType      string  `json:"PostType"`      // 帖子类型。
	SourcePostID  *int64  `json:"SourcePostID"`  // 来源帖子 ID。
	RepostComment string  `json:"RepostComment"` // 转发文案。
	CircleID      *int64  `json:"CircleID"`      // 圈子 ID。
	Visibility    string  `json:"Visibility"`    // 可见性。
	Status        string  `json:"Status"`        // 帖子状态。
	CoverURL      string  `json:"CoverURL"`      // 封面图地址。
	ViewCount     int64   `json:"ViewCount"`     // 浏览数。
	LikeCount     int64   `json:"LikeCount"`     // 点赞数。
	CommentCount  int64   `json:"CommentCount"`  // 评论数。
	FavoriteCount int64   `json:"FavoriteCount"` // 收藏数。
	ShareCount    int64   `json:"ShareCount"`    // 分享数。
	RepostCount   int64   `json:"RepostCount"`   // 转发数。
	HotScore      int64   `json:"HotScore"`      // 热度值。
	ScheduledAt   *string `json:"ScheduledAt"`   // 预约发布时间。
	PublishedAt   *string `json:"PublishedAt"`   // 实际发布时间。
	CreatedAt     string  `json:"CreatedAt"`     // 创建时间。
	UpdatedAt     string  `json:"UpdatedAt"`     // 更新时间。
}

// AdminCommentItem 后台评论列表项。
type AdminCommentItem struct {
	ID            int64  `json:"ID"`            // 评论主键 ID。
	PostID        int64  `json:"PostID"`        // 帖子 ID。
	UserID        int64  `json:"UserID"`        // 用户 ID。
	ParentID      int64  `json:"ParentID"`      // 父评论 ID。
	RootID        int64  `json:"RootID"`        // 根评论 ID。
	ReplyToUserID *int64 `json:"ReplyToUserID"` // 被回复用户 ID。
	Content       string `json:"Content"`       // 评论内容。
	Status        string `json:"Status"`        // 评论状态。
	LikeCount     int64  `json:"LikeCount"`     // 点赞数。
	CreatedAt     string `json:"CreatedAt"`     // 创建时间。
	UpdatedAt     string `json:"UpdatedAt"`     // 更新时间。
}

// AdminCircleItem 后台圈子列表项。
type AdminCircleItem struct {
	ID             int64  `json:"ID"`             // 圈子主键 ID。
	OwnerID        int64  `json:"OwnerID"`        // 圈主 ID。
	Name           string `json:"Name"`           // 圈子名称。
	Avatar         string `json:"Avatar"`         // 圈子头像。
	Description    string `json:"Description"`    // 圈子描述。
	Category       string `json:"Category"`       // 圈子分类。
	JoinType       string `json:"JoinType"`       // 加入方式。
	PostPermission string `json:"PostPermission"` // 发帖权限。
	Rules          string `json:"Rules"`          // 圈规说明。
	Status         string `json:"Status"`         // 圈子状态。
	MemberCount    int64  `json:"MemberCount"`    // 成员数。
	PostCount      int64  `json:"PostCount"`      // 帖子数。
	FeaturedCount  int64  `json:"FeaturedCount"`  // 精选数。
	IsRecommended  bool   `json:"IsRecommended"`  // 是否推荐。
	CreatedAt      string `json:"CreatedAt"`      // 创建时间。
	UpdatedAt      string `json:"UpdatedAt"`      // 更新时间。
}

// AdminOperationLogItem 后台操作日志列表项。
type AdminOperationLogItem struct {
	ID         int64  `json:"ID"`         // 日志主键 ID。
	AdminID    int64  `json:"AdminID"`    // 管理员 ID。
	AdminName  string `json:"AdminName"`  // 管理员名称。
	Action     string `json:"Action"`     // 操作动作。
	TargetType string `json:"TargetType"` // 目标类型。
	TargetID   int64  `json:"TargetID"`   // 目标 ID。
	Detail     string `json:"Detail"`     // 详情描述。
	CreatedAt  string `json:"CreatedAt"`  // 创建时间。
}

// CreateTagRequest 后台创建标签请求。
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required" example:"Go"` // 标签名称。
	Description string `json:"description" example:"Golang 相关内容"`    // 标签描述。
}

// UpdateTagRequest 后台更新标签请求。
type UpdateTagRequest struct {
	Name        *string `json:"name" example:"Go"`               // 标签名称。
	Description *string `json:"description" example:"Golang 内容"` // 标签描述。
	Status      *string `json:"status" example:"enabled"`        // 标签状态。
}

// CreateTopicRequest 后台创建话题请求。
type CreateTopicRequest struct {
	Name        string `json:"name" binding:"required" example:"Feedora 官方话题"` // 话题名称。
	Description string `json:"description" example:"话题描述"`                     // 话题描述。
	CoverImage  string `json:"coverImage" example:"/static/topic-cover.png"`   // 封面图片。
	IsOfficial  bool   `json:"isOfficial" example:"true"`                      // 是否官方话题。
}

// UpdateTopicRequest 后台更新话题请求。
type UpdateTopicRequest struct {
	Name          *string `json:"name" example:"Feedora 官方话题"`            // 话题名称。
	Description   *string `json:"description" example:"更新后的话题描述"`         // 话题描述。
	CoverImage    *string `json:"coverImage" example:"/static/topic.png"` // 封面图片。
	IsOfficial    *bool   `json:"isOfficial" example:"true"`              // 是否官方话题。
	IsRecommended *bool   `json:"isRecommended" example:"true"`           // 是否推荐。
	Status        *string `json:"status" example:"enabled"`               // 话题状态。
}

// AdminUserStatusRequest 后台变更用户状态请求（封禁 / 解禁）。
type AdminUserStatusRequest struct {
	Status string `json:"status" binding:"required" example:"banned"` // 目标状态：normal / banned。
}

// AdminPostStatusRequest 后台变更帖子状态请求（上下架）。
type AdminPostStatusRequest struct {
	Status string `json:"status" binding:"required" example:"takedown"` // 目标状态：published / hidden / takedown。
}

// AdminStatsResponse 后台统计响应。
type AdminStatsResponse struct {
	TraceEnvelope
	Data AdminStats `json:"data"` // 业务数据。
}

// AdminPostPageResponse 后台帖子分页响应。
type AdminPostPageResponse struct {
	TraceEnvelope
	Data PageResult[AdminPostItem] `json:"data"` // 业务数据。
}

// AdminCommentPageResponse 后台评论分页响应。
type AdminCommentPageResponse struct {
	TraceEnvelope
	Data PageResult[AdminCommentItem] `json:"data"` // 业务数据。
}

// AdminCirclePageResponse 后台圈子分页响应。
type AdminCirclePageResponse struct {
	TraceEnvelope
	Data PageResult[AdminCircleItem] `json:"data"` // 业务数据。
}

// AdminOperationLogPageResponse 后台操作日志分页响应。
type AdminOperationLogPageResponse struct {
	TraceEnvelope
	Data PageResult[AdminOperationLogItem] `json:"data"` // 业务数据。
}
