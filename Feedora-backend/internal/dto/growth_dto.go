package dto

// CheckInResult 签到响应。
type CheckInResult struct {
	Points         int64 `json:"points"`
	ContinuousDays int   `json:"continuousDays"`
	AwardedPoints  int   `json:"awardedPoints"`
}

// Task 成长任务项（对应前端 Task）。
type Task struct {
	TaskID       int64  `json:"taskId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	RewardPoints int    `json:"rewardPoints"`
	TargetValue  int    `json:"targetValue"`
	CurrentValue int    `json:"currentValue"`
	Status       string `json:"status"`
	ActionText   string `json:"actionText"`
	ActionURL    string `json:"actionUrl"`
}
