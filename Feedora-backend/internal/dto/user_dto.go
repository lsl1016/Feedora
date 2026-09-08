package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// UserSummary 用户简要信息。
type UserSummary struct {
	UserID    int64  `json:"userId"`              // 用户 ID。
	Nickname  string `json:"nickname"`            // 用户昵称。
	Avatar    string `json:"avatar"`              // 头像地址。
	Bio       string `json:"bio,omitempty"`       // 个性签名。
	Level     int    `json:"level,omitempty"`     // 用户等级。
	LevelName string `json:"levelName,omitempty"` // 等级名称。
}

// User 用户完整信息。
type User struct {
	UserSummary
	Account               string `json:"account"`               // 登录账号。
	Role                  string `json:"role"`                  // 用户角色。
	Status                string `json:"status"`                // 用户状态。
	Points                int64  `json:"points"`                // 当前积分。
	Experience            int64  `json:"experience"`            // 当前经验值。
	NextLevelExperience   int64  `json:"nextLevelExperience"`   // 下一等级所需经验。
	BadgeCount            int    `json:"badgeCount"`            // 勋章数量。
	CheckedInToday        bool   `json:"checkedInToday"`        // 今日是否已签到。
	ContinuousCheckInDays int    `json:"continuousCheckInDays"` // 连续签到天数。
	PostCount             int64  `json:"postCount"`             // 发帖数。
	CommentCount          int64  `json:"commentCount"`          // 评论数。
	FollowerCount         int64  `json:"followerCount"`         // 粉丝数。
	FollowingCount        int64  `json:"followingCount"`        // 关注数。
	LikeReceivedCount     int64  `json:"likeReceivedCount"`     // 获赞数。
	CreatedAt             string `json:"createdAt"`             // 创建时间。
	UpdatedAt             string `json:"updatedAt"`             // 更新时间。
}

// UpdateProfileRequest 修改资料请求。
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname" example:"新的昵称"`             // 用户昵称。
	Avatar   *string `json:"avatar" example:"/static/avatar.png"` // 头像地址。
	Bio      *string `json:"bio" example:"热爱分享的 Feedora 用户"`      // 个性签名。
}

// UserListQuery 用户列表请求。
type UserListQuery struct {
	PageRequest
	Keyword string `form:"keyword" example:"feedora"` // 搜索关键词。
}

// UserListResponse 用户列表分页响应。
type UserListResponse struct {
	TraceEnvelope
	Data PageResult[User] `json:"data"` // 业务数据。
}

// UserResponse 单个用户响应。
type UserResponse struct {
	TraceEnvelope
	Data User `json:"data"` // 业务数据。
}

// PostListResponse 帖子列表分页响应。
type PostListResponse struct {
	TraceEnvelope
	Data PageResult[Post] `json:"data"` // 业务数据。
}

// levelName 根据等级返回等级名称。
func levelName(level int) string {
	switch {
	case level >= 10:
		return "资深专家"
	case level >= 6:
		return "活跃达人"
	case level >= 3:
		return "进阶用户"
	default:
		return "新手上路"
	}
}

// LevelNameOf 供其它模块（如排行榜）组装 DTO 时复用等级名称。
func LevelNameOf(level int) string {
	return levelName(level)
}

// ToUserSummary 将用户模型转换为简要信息。
func ToUserSummary(u *model.User) UserSummary {
	if u == nil {
		return UserSummary{}
	}
	return UserSummary{
		UserID:    u.ID,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Level:     u.Level,
		LevelName: levelName(u.Level),
	}
}

// ToUser 将用户模型转换为完整信息。experience/nextLevel 阶段一用积分近似。
func ToUser(u *model.User, checkedInToday bool, continuousDays int) User {
	return User{
		UserSummary:           ToUserSummary(u),
		Account:               u.Account,
		Role:                  u.Role,
		Status:                u.Status,
		Points:                u.PointCount,
		Experience:            u.PointCount,
		NextLevelExperience:   int64(u.Level) * 100,
		BadgeCount:            0,
		CheckedInToday:        checkedInToday,
		ContinuousCheckInDays: continuousDays,
		PostCount:             u.PostCount,
		CommentCount:          u.CommentCount,
		FollowerCount:         u.FollowerCount,
		FollowingCount:        u.FollowingCount,
		LikeReceivedCount:     u.LikeCount,
		CreatedAt:             utils.FormatTime(u.CreatedAt),
		UpdatedAt:             utils.FormatTime(u.UpdatedAt),
	}
}
