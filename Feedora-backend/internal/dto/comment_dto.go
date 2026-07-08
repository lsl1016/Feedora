package dto

// Comment 评论（含子回复）。
type Comment struct {
	CommentID int64       `json:"commentId"` // 评论 ID。
	PostID    int64       `json:"postId"`    // 所属帖子 ID。
	UserID    int64       `json:"userId"`    // 评论用户 ID。
	User      UserSummary `json:"user"`      // 评论用户信息。
	Content   string      `json:"content"`   // 评论内容。
	LikeCount int64       `json:"likeCount"` // 点赞数量。
	Liked     bool        `json:"liked"`     // 当前用户是否已点赞。
	Status    string      `json:"status"`    // 评论状态。
	Replies   []Comment   `json:"replies"`   // 子回复列表。
	CreatedAt string      `json:"createdAt"` // 创建时间。
	UpdatedAt string      `json:"updatedAt"` // 更新时间。
}

// MyCommentItem 我的评论列表项。
type MyCommentItem struct {
	CommentID int64  `json:"commentId"` // 评论 ID。
	PostID    int64  `json:"postId"`    // 帖子 ID。
	PostTitle string `json:"postTitle"` // 帖子标题。
	Content   string `json:"content"`   // 评论内容。
	LikeCount int64  `json:"likeCount"` // 点赞数量。
	Status    string `json:"status"`    // 评论状态。
	CreatedAt string `json:"createdAt"` // 创建时间。
}

// CreateCommentRequest 发表评论请求。
type CreateCommentRequest struct {
	PostID  int64  `json:"postId" binding:"required,min=1" example:"1"`   // 帖子 ID。
	Content string `json:"content" binding:"required" example:"这是一条评论内容"` // 评论内容。
}

// ReplyCommentRequest 回复评论请求。
type ReplyCommentRequest struct {
	Content       string `json:"content" binding:"required" example:"这是一条回复内容"` // 回复内容。
	ReplyToUserID *int64 `json:"replyToUserId" example:"2"`                     // 被回复用户 ID。
}

// CommentListResponse 评论列表响应。
type CommentListResponse struct {
	TraceEnvelope
	Data []Comment `json:"data"` // 业务数据。
}

// CommentResponse 单条评论响应。
type CommentResponse struct {
	TraceEnvelope
	Data Comment `json:"data"` // 业务数据。
}

// MyCommentPageResponse 我的评论分页响应。
type MyCommentPageResponse struct {
	TraceEnvelope
	Data PageResult[MyCommentItem] `json:"data"` // 业务数据。
}
