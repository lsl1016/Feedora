package dto

// GrowthTaskQuery 成长任务列表请求。
type GrowthTaskQuery struct {
	Type string `form:"type" example:"daily"` // 任务类型。
}

// GrowthRankingQuery 成长排行榜请求。
type GrowthRankingQuery struct {
	PageRequest
	Type  string `form:"type" example:"active"` // 排行类型。
	Range string `form:"range" example:"all"`   // 时间范围。
}

// CheckInResult 签到响应。
type CheckInResult struct {
	Points         int64 `json:"points"`         // 当前积分。
	ContinuousDays int   `json:"continuousDays"` // 连续签到天数。
	AwardedPoints  int   `json:"awardedPoints"`  // 本次发放积分。
}

// Task 成长任务项（对应前端 Task）。
type Task struct {
	TaskID       int64  `json:"taskId"`       // 任务 ID。
	Title        string `json:"title"`        // 任务标题。
	Description  string `json:"description"`  // 任务描述。
	Type         string `json:"type"`         // 任务类型。
	RewardPoints int    `json:"rewardPoints"` // 奖励积分。
	TargetValue  int    `json:"targetValue"`  // 目标值。
	CurrentValue int    `json:"currentValue"` // 当前进度值。
	Status       string `json:"status"`       // 任务状态。
	ActionText   string `json:"actionText"`   // 操作文案。
	ActionURL    string `json:"actionUrl"`    // 操作跳转链接。
}

// CheckInResponse 签到响应。
type CheckInResponse struct {
	TraceEnvelope
	Data CheckInResult `json:"data"` // 业务数据。
}

// TaskListResponse 任务列表响应。
type TaskListResponse struct {
	TraceEnvelope
	Data []Task `json:"data"` // 业务数据。
}

// RankingListResponse 排行列表响应。
type RankingListResponse struct {
	TraceEnvelope
	Data []RankingItem `json:"data"` // 业务数据。
}
