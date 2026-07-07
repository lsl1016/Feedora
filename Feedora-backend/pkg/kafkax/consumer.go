package kafkax

import (
	"context"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// EnsureTopics 预创建 Topic（单分区单副本），避免消费组在 Topic 不存在时拿不到分区分配。
func EnsureTopics(brokers []string, topics []string) error {
	if len(brokers) == 0 {
		return nil
	}
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	cc, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer cc.Close()
	cfgs := make([]kafka.TopicConfig, 0, len(topics))
	for _, t := range topics {
		cfgs = append(cfgs, kafka.TopicConfig{Topic: t, NumPartitions: 1, ReplicationFactor: 1})
	}
	return cc.CreateTopics(cfgs...)
}

// Consumer 基于消费组订阅多个 Topic 的消费者。
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer 创建消费组消费者，订阅给定的多个 Topic。
func NewConsumer(brokers []string, groupID string, topics []string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupID,
			GroupTopics: topics,
			MinBytes:    1,
			MaxBytes:    10e6,
		}),
	}
}

// Message 消费到的消息（对上层屏蔽 kafka-go 细节）。
type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

// Fetch 拉取一条消息（不自动提交），阻塞直到有消息或 ctx 取消。
func (c *Consumer) Fetch(ctx context.Context) (*Message, kafka.Message, error) {
	m, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, m, err
	}
	return &Message{Topic: m.Topic, Key: m.Key, Value: m.Value}, m, nil
}

// Commit 提交已处理的消息位点。
func (c *Consumer) Commit(ctx context.Context, raw kafka.Message) error {
	return c.reader.CommitMessages(ctx, raw)
}

// Close 关闭消费者。
func (c *Consumer) Close() error {
	return c.reader.Close()
}
