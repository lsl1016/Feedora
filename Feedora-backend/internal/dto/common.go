package dto

import "github.com/feedora/backend/pkg/response"

// 本包定义与前端 src/types.ts 对应的请求 / 响应结构，作为前后端 JSON 契约的唯一来源。

// EmptyData 表示无业务数据的成功响应体。
type EmptyData struct{}

// TraceEnvelope 表示统一响应的公共字段。
type TraceEnvelope struct {
	Code      int    `json:"code" example:"0"`                         // 业务状态码。
	Message   string `json:"message" example:"success"`                // 业务提示信息。
	TraceID   string `json:"traceId" example:"trace-1234567890"`       // 链路追踪 ID。
	Timestamp string `json:"timestamp" example:"2026-07-08T12:00:00Z"` // 响应时间。
}

// EmptyResponse 表示无业务数据的成功响应。
type EmptyResponse struct {
	TraceEnvelope
	Data EmptyData `json:"data"` // 业务数据。
}

// CountData 表示数量型响应数据。
type CountData struct {
	Count int64 `json:"count" example:"3"` // 统计数量。
}

// CountResponse 表示数量型成功响应。
type CountResponse struct {
	TraceEnvelope
	Data CountData `json:"data"` // 业务数据。
}

// ClaimedData 表示任务领取结果。
type ClaimedData struct {
	Claimed bool `json:"claimed" example:"true"` // 是否领取成功。
}

// ClaimedResponse 表示任务领取成功响应。
type ClaimedResponse struct {
	TraceEnvelope
	Data ClaimedData `json:"data"` // 业务数据。
}

// PageRequest 分页请求参数。
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"omitempty,min=1" example:"1"`                  // 页码，从 1 开始。
	PageSize int `json:"pageSize" form:"pageSize" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量，最大 100。
}

// IDURI 表示单个路径 ID 参数。
type IDURI struct {
	ID int64 `uri:"id" binding:"required,min=1" example:"1"` // 资源 ID。
}

// PostIDURI 表示帖子路径参数。
type PostIDURI struct {
	PostID int64 `uri:"postId" binding:"required,min=1" example:"1"` // 帖子 ID。
}

// CommentIDURI 表示评论路径参数。
type CommentIDURI struct {
	CommentID int64 `uri:"commentId" binding:"required,min=1" example:"1"` // 评论 ID。
}

// UserIDURI 表示用户路径参数。
type UserIDURI struct {
	UserID int64 `uri:"userId" binding:"required,min=1" example:"1"` // 用户 ID。
}

// CircleIDURI 表示圈子路径参数。
type CircleIDURI struct {
	CircleID int64 `uri:"circleId" binding:"required,min=1" example:"1"` // 圈子 ID。
}

// TopicIDURI 表示话题路径参数。
type TopicIDURI struct {
	TopicID int64 `uri:"topicId" binding:"required,min=1" example:"1"` // 话题 ID。
}

// TagIDURI 表示标签路径参数。
type TagIDURI struct {
	TagID int64 `uri:"tagId" binding:"required,min=1" example:"1"` // 标签 ID。
}

// NotificationIDURI 表示通知路径参数。
type NotificationIDURI struct {
	NotificationID int64 `uri:"notificationId" binding:"required,min=1" example:"1"` // 通知 ID。
}

// TaskIDURI 表示任务路径参数。
type TaskIDURI struct {
	TaskID int64 `uri:"taskId" binding:"required,min=1" example:"1"` // 任务 ID。
}

// UserAndCircleURI 表示圈子成员路径参数。
type UserAndCircleURI struct {
	CircleID int64 `uri:"circleId" binding:"required,min=1" example:"1"` // 圈子 ID。
	UserID   int64 `uri:"userId" binding:"required,min=1" example:"2"`   // 用户 ID。
}

// SortPageRequest 表示带排序的分页请求。
type SortPageRequest struct {
	PageRequest
	Sort string `form:"sort" example:"latest"` // 排序方式。
}

// KeywordPageRequest 表示带关键字的分页请求。
type KeywordPageRequest struct {
	PageRequest
	Keyword string `form:"keyword" example:"feedora"` // 搜索关键词。
}

// PageResult 分页响应（与前端 PageResult<T> 对应）。
type PageResult[T any] struct {
	List     []T   `json:"list"`     // 列表数据。
	Total    int64 `json:"total"`    // 总记录数。
	Page     int   `json:"page"`     // 当前页码。
	PageSize int   `json:"pageSize"` // 每页数量。
}

// PostPageData 表示帖子分页数据。
type PostPageData = response.PageData

// CommentPageData 表示评论分页数据。
type CommentPageData = response.PageData

// CirclePageData 表示圈子分页数据。
type CirclePageData = response.PageData

// UserPageData 表示用户分页数据。
type UserPageData = response.PageData

// TopicPageData 表示话题分页数据。
type TopicPageData = response.PageData

// TagPageData 表示标签分页数据。
type TagPageData = response.PageData

// NotificationPageData 表示通知分页数据。
type NotificationPageData = response.PageData

// RankingPageData 表示排行分页数据。
type RankingPageData = response.PageData

// OperationLogPageData 表示操作日志分页数据。
type OperationLogPageData = response.PageData
