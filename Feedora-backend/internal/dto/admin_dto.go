package dto

// AdminStats 后台概览统计。
type AdminStats struct {
	UserCount    int64 `json:"userCount"`
	PostCount    int64 `json:"postCount"`
	CommentCount int64 `json:"commentCount"`
	CircleCount  int64 `json:"circleCount"`
}

// CreateTagRequest 后台创建标签请求。
type CreateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateTagRequest 后台更新标签请求。
type UpdateTagRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// CreateTopicRequest 后台创建话题请求。
type CreateTopicRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	IsOfficial  bool   `json:"isOfficial"`
}
