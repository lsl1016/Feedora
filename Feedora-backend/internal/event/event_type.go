package event

// 事件类型常量（阶段一核心集合，阶段二继续扩展）。
const (
	UserRegistered = "UserRegistered"
	PostCreated    = "PostCreated"
	PostUpdated    = "PostUpdated"
	PostDeleted    = "PostDeleted"
	PostHidden     = "PostHidden"
	CommentCreated = "CommentCreated"
	PostLiked      = "PostLiked"
	PostUnliked    = "PostUnliked"
	PostFavorited  = "PostFavorited"
	CircleCreated  = "CircleCreated"
	CircleJoined   = "CircleJoined"
)
