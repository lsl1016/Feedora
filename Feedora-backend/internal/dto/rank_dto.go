package dto

// HotRankQuery 热门榜单请求。
type HotRankQuery struct {
	PageRequest
	RankType  string `form:"rankType" example:"post"`   // 榜单类型。
	TimeRange string `form:"timeRange" example:"today"` // 时间范围。
}

// HotRankItem 热门榜单项，覆盖帖子 / 圈子 / 话题所需字段（多余字段前端忽略）。
type HotRankItem struct {
	Rank              int     `json:"rank"`                        // 排名。
	Score             float64 `json:"score"`                       // 排行得分。
	HotScore          int64   `json:"hotScore"`                    // 热度值。
	TargetID          int64   `json:"targetId"`                    // 目标 ID。
	PostID            int64   `json:"postId,omitempty"`            // 帖子 ID。
	CircleID          int64   `json:"circleId,omitempty"`          // 圈子 ID。
	TopicID           int64   `json:"topicId,omitempty"`           // 话题 ID。
	Title             string  `json:"title,omitempty"`             // 帖子标题。
	Name              string  `json:"name,omitempty"`              // 目标名称。
	AuthorName        string  `json:"authorName,omitempty"`        // 作者名称。
	LikeCount         int64   `json:"likeCount,omitempty"`         // 点赞数。
	CommentCount      int64   `json:"commentCount,omitempty"`      // 评论数。
	MemberCount       int64   `json:"memberCount,omitempty"`       // 成员数。
	PostCount         int64   `json:"postCount,omitempty"`         // 帖子数。
	FeaturedPostCount int64   `json:"featuredPostCount,omitempty"` // 精选帖子数。
	ParticipantCount  int64   `json:"participantCount,omitempty"`  // 参与人数。
}

// RankingItem 用户 / 圈子排行榜项（对应前端 RankingItem）。
type RankingItem struct {
	Rank              int     `json:"rank"`                        // 排名。
	TargetID          int64   `json:"targetId"`                    // 目标 ID。
	TargetType        string  `json:"targetType"`                  // 目标类型。
	Name              string  `json:"name"`                        // 展示名称。
	Avatar            string  `json:"avatar"`                      // 展示头像。
	Level             int     `json:"level,omitempty"`             // 用户等级。
	LevelName         string  `json:"levelName,omitempty"`         // 等级名称。
	Points            int64   `json:"points,omitempty"`            // 积分。
	PostCount         int64   `json:"postCount,omitempty"`         // 发帖数。
	LikeReceivedCount int64   `json:"likeReceivedCount,omitempty"` // 获赞数。
	MemberCount       int64   `json:"memberCount,omitempty"`       // 成员数。
	Score             float64 `json:"score"`                       // 排行得分。
	IsCurrentUser     bool    `json:"isCurrentUser,omitempty"`     // 是否当前用户。
}

// HotRankListResponse 热门榜单响应。
type HotRankListResponse struct {
	TraceEnvelope
	Data []HotRankItem `json:"data"` // 业务数据。
}
