package event

import (
	"encoding/json"
	"time"

	"github.com/feedora/backend/pkg/logger"
	"github.com/feedora/backend/pkg/utils"
)

// Producer 发布领域事件。阶段一使用 NoopProducer；
// 阶段二可替换为基于 Outbox + Kafka 的可靠投递实现，接口保持不变。
type Producer interface {
	Publish(topic, eventType string, bizID, userID int64, payload interface{}) error
}

// NoopProducer 仅打印事件日志，让核心闭环在不依赖 Kafka 的情况下运行。
type NoopProducer struct {
	topicPrefix string
}

func NewNoopProducer(topicPrefix string) *NoopProducer {
	if topicPrefix == "" {
		topicPrefix = "community"
	}
	return &NoopProducer{topicPrefix: topicPrefix}
}

func (p *NoopProducer) Publish(topic, eventType string, bizID, userID int64, payload interface{}) error {
	e := Event{
		EventID:   utils.UUID(),
		EventType: eventType,
		BizID:     bizID,
		UserID:    userID,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
	b, _ := json.Marshal(e)
	logger.Infof("event (noop) topic=%s.%s %s", p.topicPrefix, topic, string(b))
	return nil
}
