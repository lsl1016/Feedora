package dto

// InteractionResponse 点赞或收藏响应。
type InteractionResponse struct {
	TraceEnvelope
	Data InteractionResult `json:"data"` // 业务数据。
}

// InteractionResult 点赞 / 收藏后返回的最新状态，便于前端即时更新。
type InteractionResult struct {
	Liked         bool  `json:"liked"`         // 是否已点赞。
	Favorited     bool  `json:"favorited"`     // 是否已收藏。
	LikeCount     int64 `json:"likeCount"`     // 点赞数。
	FavoriteCount int64 `json:"favoriteCount"` // 收藏数。
}
