package event

import (
	"encoding/json"
	"time"

	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	"github.com/feedora/backend/pkg/utils"
)

// Message 是投递到 Kafka 的事件消息线格式，Worker 按此结构反序列化。
type Message struct {
	EventID     string          `json:"eventId"`
	EventType   string          `json:"eventType"`
	AggregateID int64           `json:"aggregateId"`
	UserID      int64           `json:"userId"`
	TraceID     string          `json:"traceId"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// OutboxProducer 通过写 event_outbox 表实现事件的可靠投递（阶段二默认实现）。
// 业务写库成功后写入 outbox，由 Outbox Dispatcher 异步投递 Kafka，避免事件丢失。
type OutboxProducer struct {
	outbox      *repository.OutboxRepository
	topicPrefix string
}

func NewOutboxProducer(outbox *repository.OutboxRepository, topicPrefix string) *OutboxProducer {
	if topicPrefix == "" {
		topicPrefix = "feedora"
	}
	return &OutboxProducer{outbox: outbox, topicPrefix: topicPrefix}
}

// Publish 将事件写入 event_outbox 表。
func (p *OutboxProducer) Publish(topic, eventType string, bizID, userID int64, payload interface{}) error {
	raw := []byte("{}")
	if payload != nil {
		if b, err := json.Marshal(payload); err == nil {
			raw = b
		}
	}
	now := time.Now()
	return p.outbox.Insert(&model.EventOutbox{
		EventID:     utils.UUID(),
		Topic:       p.topicPrefix + "." + topic,
		EventType:   eventType,
		AggregateID: bizID,
		UserID:      userID,
		Payload:     string(raw),
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}
