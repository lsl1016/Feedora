package dto

// SearchResult 搜索聚合结果项。
type SearchResult struct {
	PostID      *int64  `json:"postId,omitempty"`   // 帖子 ID。
	UserID      *int64  `json:"userId,omitempty"`   // 用户 ID。
	TopicID     *int64  `json:"topicId,omitempty"`  // 话题 ID。
	CircleID    *int64  `json:"circleId,omitempty"` // 圈子 ID。
	Post        *Post   `json:"post,omitempty"`     // 帖子信息。
	User        *User   `json:"user,omitempty"`     // 用户信息。
	Topic       *Topic  `json:"topic,omitempty"`    // 话题信息。
	Circle      *Circle `json:"circle,omitempty"`   // 圈子信息。
	ResultType  string  `json:"resultType"`         // 结果类型。
	DisplayName string  `json:"displayName"`        // 展示标题。
}

// SearchQuery 综合搜索请求。
type SearchQuery struct {
	PageRequest
	Keyword string `form:"keyword" example:"feedora"` // 搜索关键词。
	Type    string `form:"type" example:"all"`        // 搜索类型。
}

// SearchSuggestQuery 搜索联想请求。
type SearchSuggestQuery struct {
	Keyword string `form:"keyword" example:"feedora"` // 搜索关键词。
}

// SuggestItem 搜索联想项。
type SuggestItem struct {
	Type      string `json:"type"`      // 联想类型。
	Title     string `json:"title"`     // 联想标题。
	TargetID  int64  `json:"targetId"`  // 目标 ID。
	TargetURL string `json:"targetUrl"` // 跳转链接。
}

// SearchResultPageResponse 综合搜索分页响应。
type SearchResultPageResponse struct {
	TraceEnvelope
	Data PageResult[SearchResult] `json:"data"` // 业务数据。
}

// SearchSuggestResponse 搜索联想响应。
type SearchSuggestResponse struct {
	TraceEnvelope
	Data []SuggestItem `json:"data"` // 业务数据。
}

// HotKeywordResponse 热门关键词响应。
type HotKeywordResponse struct {
	TraceEnvelope
	Data []string `json:"data"` // 业务数据。
}
