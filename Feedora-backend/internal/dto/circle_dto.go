package dto

import (
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/utils"
)

// CircleBrief 圈子简要信息。
type CircleBrief struct {
	CircleID    int64  `json:"circleId"`    // 圈子 ID。
	Name        string `json:"name"`        // 圈子名称。
	Avatar      string `json:"avatar"`      // 圈子头像。
	MemberCount int64  `json:"memberCount"` // 成员数量。
	PostCount   int64  `json:"postCount"`   // 帖子数量。
}

// Circle 圈子完整信息。
type Circle struct {
	CircleID          int64        `json:"circleId"`           // 圈子 ID。
	Name              string       `json:"name"`               // 圈子名称。
	Avatar            string       `json:"avatar"`             // 圈子头像。
	Description       string       `json:"description"`        // 圈子描述。
	Category          string       `json:"category"`           // 圈子分类。
	Tags              []TagSummary `json:"tags"`               // 标签列表。
	OwnerID           int64        `json:"ownerId"`            // 圈主用户 ID。
	Owner             UserSummary  `json:"owner"`              // 圈主信息。
	JoinType          string       `json:"joinType"`           // 加入方式。
	PostPermission    string       `json:"postPermission"`     // 发帖权限。
	MemberCount       int64        `json:"memberCount"`        // 成员数量。
	PostCount         int64        `json:"postCount"`          // 帖子数量。
	FeaturedPostCount int64        `json:"featuredPostCount"`  // 精选帖子数量。
	IsJoined          bool         `json:"isJoined"`           // 当前用户是否已加入。
	MyRole            string       `json:"myRole,omitempty"`   // 当前用户在圈子内角色。
	MyStatus          string       `json:"myStatus,omitempty"` // 当前用户在圈子内状态。
	IsRecommended     bool         `json:"isRecommended"`      // 是否推荐。
	Status            string       `json:"status"`             // 圈子状态。
	Rules             string       `json:"rules"`              // 圈规说明。
	CreatedAt         string       `json:"createdAt"`          // 创建时间。
	UpdatedAt         string       `json:"updatedAt"`          // 更新时间。
}

// CircleMember 圈子成员。
type CircleMember struct {
	ID         int64       `json:"id"`                   // 成员记录 ID。
	CircleID   int64       `json:"circleId"`             // 圈子 ID。
	UserID     int64       `json:"userId"`               // 用户 ID。
	User       UserSummary `json:"user"`                 // 用户信息。
	Role       string      `json:"role"`                 // 成员角色。
	Status     string      `json:"status"`               // 成员状态。
	MuteReason string      `json:"muteReason,omitempty"` // 禁言原因。
	MutedUntil string      `json:"mutedUntil,omitempty"` // 禁言结束时间。
	JoinedAt   string      `json:"joinedAt"`             // 加入时间。
	UpdatedAt  string      `json:"updatedAt"`            // 更新时间。
}

// CreateCircleRequest 创建圈子请求。
type CreateCircleRequest struct {
	Name           string `json:"name" binding:"required" example:"Feedora 圈子"` // 圈子名称。
	Avatar         string `json:"avatar" example:"/static/circle.png"`          // 圈子头像。
	Description    string `json:"description" example:"圈子简介"`                   // 圈子描述。
	Category       string `json:"category" example:"technology"`                // 圈子分类。
	JoinType       string `json:"joinType" example:"direct"`                    // 加入方式。
	PostPermission string `json:"postPermission" example:"all"`                 // 发帖权限。
	Rules          string `json:"rules" example:"请遵守社区规范"`                      // 圈规说明。
}

// CircleListQuery 圈子列表请求。
type CircleListQuery struct {
	PageRequest
	Scope    string `form:"scope" example:"all"`           // 范围。
	Keyword  string `form:"keyword" example:"feedora"`     // 搜索关键词。
	Category string `form:"category" example:"technology"` // 分类。
	Sort     string `form:"sort" example:"hot"`            // 排序方式。
}

// SetCircleMemberRoleRequest 设置成员角色请求。
type SetCircleMemberRoleRequest struct {
	Role string `json:"role" binding:"required" example:"admin"` // 成员角色。
}

// MuteCircleMemberRequest 禁言成员请求。
type MuteCircleMemberRequest struct {
	Duration int    `json:"duration" binding:"required,min=1" example:"7"` // 禁言天数。
	Reason   string `json:"reason" example:"违反圈规"`                         // 禁言原因。
}

// CircleResponse 单个圈子响应。
type CircleResponse struct {
	TraceEnvelope
	Data Circle `json:"data"` // 业务数据。
}

// CirclePageResponse 圈子分页响应。
type CirclePageResponse struct {
	TraceEnvelope
	Data PageResult[Circle] `json:"data"` // 业务数据。
}

// CircleMemberPageResponse 圈子成员分页响应。
type CircleMemberPageResponse struct {
	TraceEnvelope
	Data PageResult[CircleMember] `json:"data"` // 业务数据。
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
