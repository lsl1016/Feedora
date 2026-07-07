package model

// 用户角色与状态。
const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	UserNormal = "normal"
	UserMuted  = "muted"
	UserBanned = "banned"
)

// 帖子状态。
const (
	PostDraft     = "draft"
	PostScheduled = "scheduled"
	PostReviewing = "reviewing"
	PostPublished = "published"
	PostHidden    = "hidden"
	PostRejected  = "rejected"
	PostDeleted   = "deleted"
	PostTakedown  = "takedown"
)

// 帖子可见性与类型。
const (
	VisibilityPublic     = "public"
	VisibilityCircleOnly = "circle_only"

	PostTypeOriginal = "original"
	PostTypeRepost   = "repost"
)

// 评论状态。
const (
	CommentNormal   = "normal"
	CommentDeleted  = "deleted"
	CommentRejected = "rejected"
)

// 圈子成员角色与状态。
const (
	CircleRoleOwner     = "owner"
	CircleRoleModerator = "moderator"
	CircleRoleReviewer  = "reviewer"
	CircleRoleMember    = "member"

	CircleMemberNormal  = "normal"
	CircleMemberMuted   = "muted"
	CircleMemberRemoved = "removed"
)

// 通用启用状态。
const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"
)
