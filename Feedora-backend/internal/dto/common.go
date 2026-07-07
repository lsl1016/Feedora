package dto

// 本包定义与前端 src/types.ts 对应的请求 / 响应结构，作为前后端 JSON 契约的唯一来源。

// PageRequest 分页请求参数。
type PageRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// PageResult 分页响应（与前端 PageResult<T> 对应）。
type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
