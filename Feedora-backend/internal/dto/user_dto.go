package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// UserSummary 用户简要信息。
type UserSummary struct {
	UserID    int64  `json:"userId"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio,omitempty"`
	Level     int    `json:"level,omitempty"`
	LevelName string `json:"levelName,omitempty"`
}

// User 用户完整信息。
type User struct {
	UserSummary
	Account               string `json:"account"`
	Role                  string `json:"role"`
	Status                string `json:"status"`
	Points                int64  `json:"points"`
	Experience            int64  `json:"experience"`
	NextLevelExperience   int64  `json:"nextLevelExperience"`
	BadgeCount            int    `json:"badgeCount"`
	CheckedInToday        bool   `json:"checkedInToday"`
	ContinuousCheckInDays int    `json:"continuousCheckInDays"`
	PostCount             int64  `json:"postCount"`
	CommentCount          int64  `json:"commentCount"`
	FollowerCount         int64  `json:"followerCount"`
	FollowingCount        int64  `json:"followingCount"`
	LikeReceivedCount     int64  `json:"likeReceivedCount"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// UpdateProfileRequest 修改资料请求。
type UpdateProfileRequest struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
	Bio      *string `json:"bio"`
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
