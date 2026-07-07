package dto

// HotRankItem 热门榜单项，覆盖帖子 / 圈子 / 话题所需字段（多余字段前端忽略）。
type HotRankItem struct {
	Rank              int     `json:"rank"`
	Score             float64 `json:"score"`
	HotScore          int64   `json:"hotScore"`
	TargetID          int64   `json:"targetId"`
	PostID            int64   `json:"postId,omitempty"`
	CircleID          int64   `json:"circleId,omitempty"`
	TopicID           int64   `json:"topicId,omitempty"`
	Title             string  `json:"title,omitempty"`
	Name              string  `json:"name,omitempty"`
	AuthorName        string  `json:"authorName,omitempty"`
	LikeCount         int64   `json:"likeCount,omitempty"`
	CommentCount      int64   `json:"commentCount,omitempty"`
	MemberCount       int64   `json:"memberCount,omitempty"`
	PostCount         int64   `json:"postCount,omitempty"`
	FeaturedPostCount int64   `json:"featuredPostCount,omitempty"`
	ParticipantCount  int64   `json:"participantCount,omitempty"`
}

// RankingItem 用户 / 圈子排行榜项（对应前端 RankingItem）。
type RankingItem struct {
	Rank              int     `json:"rank"`
	TargetID          int64   `json:"targetId"`
	TargetType        string  `json:"targetType"`
	Name              string  `json:"name"`
	Avatar            string  `json:"avatar"`
	Level             int     `json:"level,omitempty"`
	LevelName         string  `json:"levelName,omitempty"`
	Points            int64   `json:"points,omitempty"`
	PostCount         int64   `json:"postCount,omitempty"`
	LikeReceivedCount int64   `json:"likeReceivedCount,omitempty"`
	MemberCount       int64   `json:"memberCount,omitempty"`
	Score             float64 `json:"score"`
	IsCurrentUser     bool    `json:"isCurrentUser,omitempty"`
}
