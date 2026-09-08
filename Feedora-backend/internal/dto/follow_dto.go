package dto

// FollowingFeedItem 关注动态条目（对应前端 FollowingFeedItem）。
type FollowingFeedItem struct {
	FeedID       int64  `json:"feedId"`                 // 动态 ID（当前为来源帖子 ID）。
	FeedType     string `json:"feedType"`               // 动态类型：post / circle / topic。
	Title        string `json:"title"`                  // 动态标题。
	Summary      string `json:"summary"`                // 动态摘要。
	TargetID     int64  `json:"targetId"`               // 跳转目标 ID。
	TargetURL    string `json:"targetUrl"`              // 前端跳转路径。
	SourceName   string `json:"sourceName"`             // 来源作者昵称。
	SourceAvatar string `json:"sourceAvatar,omitempty"` // 来源作者头像。
	CreatedAt    string `json:"createdAt"`              // 发生时间。
}

// FollowState 关注状态。
type FollowState struct {
	Following bool `json:"following"` // 当前用户是否已关注目标用户。
}

// FollowingFeedQuery 关注动态查询参数。
type FollowingFeedQuery struct {
	PageRequest
	FeedTab string `form:"feedTab" example:"all"` // 动态范围：all / user / circle / topic。
}

// FollowStateResponse 关注状态响应。
type FollowStateResponse struct {
	TraceEnvelope
	Data FollowState `json:"data"` // 业务数据。
}

// FollowingPageResponse 关注的人分页响应。
type FollowingPageResponse struct {
	TraceEnvelope
	Data PageResult[User] `json:"data"` // 业务数据。
}

// FollowingFeedPageResponse 关注动态分页响应。
type FollowingFeedPageResponse struct {
	TraceEnvelope
	Data PageResult[FollowingFeedItem] `json:"data"` // 业务数据。
}

// TagPageResponse 关注的标签分页响应。
type TagPageResponse struct {
	TraceEnvelope
	Data PageResult[ContentTag] `json:"data"` // 业务数据。
}
