package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// CircleBrief 圈子简要信息。
type CircleBrief struct {
	CircleID    int64  `json:"circleId"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	MemberCount int64  `json:"memberCount"`
	PostCount   int64  `json:"postCount"`
}

// Circle 圈子完整信息。
type Circle struct {
	CircleID          int64        `json:"circleId"`
	Name              string       `json:"name"`
	Avatar            string       `json:"avatar"`
	Description       string       `json:"description"`
	Category          string       `json:"category"`
	Tags              []TagSummary `json:"tags"`
	OwnerID           int64        `json:"ownerId"`
	Owner             UserSummary  `json:"owner"`
	JoinType          string       `json:"joinType"`
	PostPermission    string       `json:"postPermission"`
	MemberCount       int64        `json:"memberCount"`
	PostCount         int64        `json:"postCount"`
	FeaturedPostCount int64        `json:"featuredPostCount"`
	IsJoined          bool         `json:"isJoined"`
	MyRole            string       `json:"myRole,omitempty"`
	MyStatus          string       `json:"myStatus,omitempty"`
	IsRecommended     bool         `json:"isRecommended"`
	Status            string       `json:"status"`
	Rules             string       `json:"rules"`
	CreatedAt         string       `json:"createdAt"`
	UpdatedAt         string       `json:"updatedAt"`
}

// CircleMember 圈子成员。
type CircleMember struct {
	ID         int64       `json:"id"`
	CircleID   int64       `json:"circleId"`
	UserID     int64       `json:"userId"`
	User       UserSummary `json:"user"`
	Role       string      `json:"role"`
	Status     string      `json:"status"`
	MuteReason string      `json:"muteReason,omitempty"`
	MutedUntil string      `json:"mutedUntil,omitempty"`
	JoinedAt   string      `json:"joinedAt"`
	UpdatedAt  string      `json:"updatedAt"`
}

// CreateCircleRequest 创建圈子请求。
type CreateCircleRequest struct {
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	JoinType       string `json:"joinType"`
	PostPermission string `json:"postPermission"`
	Rules          string `json:"rules"`
}

// ToCircleBrief 转换圈子简要信息。
func ToCircleBrief(c *model.Circle) *CircleBrief {
	if c == nil {
		return nil
	}
	return &CircleBrief{
		CircleID:    c.ID,
		Name:        c.Name,
		Avatar:      c.Avatar,
		MemberCount: c.MemberCount,
		PostCount:   c.PostCount,
	}
}

// ToCircle 转换圈子完整信息。member 为当前用户在该圈子的成员记录，可为 nil。
func ToCircle(c *model.Circle, owner *model.User, member *model.CircleMember) Circle {
	res := Circle{
		CircleID:          c.ID,
		Name:              c.Name,
		Avatar:            c.Avatar,
		Description:       c.Description,
		Category:          c.Category,
		Tags:              []TagSummary{},
		OwnerID:           c.OwnerID,
		Owner:             ToUserSummary(owner),
		JoinType:          c.JoinType,
		PostPermission:    c.PostPermission,
		MemberCount:       c.MemberCount,
		PostCount:         c.PostCount,
		FeaturedPostCount: c.FeaturedCount,
		IsRecommended:     c.IsRecommended,
		Status:            c.Status,
		Rules:             c.Rules,
		CreatedAt:         utils.FormatTime(c.CreatedAt),
		UpdatedAt:         utils.FormatTime(c.UpdatedAt),
	}
	if member != nil && member.Status != model.CircleMemberRemoved {
		res.IsJoined = true
		res.MyRole = member.Role
		res.MyStatus = member.Status
	}
	return res
}

// ToCircleMember 转换圈子成员。
func ToCircleMember(m *model.CircleMember, user *model.User) CircleMember {
	return CircleMember{
		ID:         m.ID,
		CircleID:   m.CircleID,
		UserID:     m.UserID,
		User:       ToUserSummary(user),
		Role:       m.Role,
		Status:     m.Status,
		MuteReason: m.MuteReason,
		MutedUntil: utils.FormatTimePtr(m.MutedUntil),
		JoinedAt:   utils.FormatTime(m.JoinedAt),
		UpdatedAt:  utils.FormatTime(m.UpdatedAt),
	}
}
