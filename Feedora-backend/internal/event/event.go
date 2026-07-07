package event

import "time"

// Event 是统一的领域事件结构。
type Event struct {
	EventID   string      `json:"eventId"`
	EventType string      `json:"eventType"`
	BizID     int64       `json:"bizId"`
	UserID    int64       `json:"userId"`
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"createdAt"`
}
