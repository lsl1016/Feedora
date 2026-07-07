package dto

// Comment 评论（含子回复）。
type Comment struct {
	CommentID int64       `json:"commentId"`
	PostID    int64       `json:"postId"`
	UserID    int64       `json:"userId"`
	User      UserSummary `json:"user"`
	Content   string      `json:"content"`
	LikeCount int64       `json:"likeCount"`
	Liked     bool        `json:"liked"`
	Status    string      `json:"status"`
	Replies   []Comment   `json:"replies"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

// MyCommentItem 我的评论列表项。
type MyCommentItem struct {
	CommentID int64  `json:"commentId"`
	PostID    int64  `json:"postId"`
	PostTitle string `json:"postTitle"`
	Content   string `json:"content"`
	LikeCount int64  `json:"likeCount"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// CreateCommentRequest 发表评论请求。
type CreateCommentRequest struct {
	PostID  int64  `json:"postId"`
	Content string `json:"content"`
}

// ReplyCommentRequest 回复评论请求。
type ReplyCommentRequest struct {
	Content       string `json:"content"`
	ReplyToUserID *int64 `json:"replyToUserId"`
}
