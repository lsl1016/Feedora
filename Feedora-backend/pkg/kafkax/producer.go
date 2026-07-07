// Package kafkax 封装 Kafka 生产者 / 消费者基础能力（基于 segmentio/kafka-go）。
package kafkax

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer Kafka 生产者，支持按消息指定 Topic。
type Producer struct {
	writer *kafka.Writer
}

// NewProducer 创建生产者，允许自动创建 Topic。
func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
			BatchTimeout:           50 * time.Millisecond,
		},
	}
}

// Publish 向指定 Topic 发送一条消息。
func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

// Close 关闭生产者。
func (p *Producer) Close() error {
	return p.writer.Close()
}
