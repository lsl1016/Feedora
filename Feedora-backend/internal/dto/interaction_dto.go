package dto

// InteractionResult 点赞 / 收藏后返回的最新状态，便于前端即时更新。
type InteractionResult struct {
	Liked         bool  `json:"liked"`
	Favorited     bool  `json:"favorited"`
	LikeCount     int64 `json:"likeCount"`
	FavoriteCount int64 `json:"favoriteCount"`
}
