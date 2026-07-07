package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/repository"
	"github.com/feedora/backend/internal/search"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/kafkax"
	"github.com/feedora/backend/pkg/logger"
	"github.com/feedora/backend/pkg/observability"
	"gorm.io/gorm"
)

// Runner Worker 运行器：负责 Outbox 投递与各类事件的异步消费。
type Runner struct {
	cfg      *config.Config
	db       *gorm.DB
	cch      *cache.Cache
	sc       *search.Client
	producer *kafkax.Producer
	consumer *kafkax.Consumer

	posts    *repository.PostRepository
	users    *repository.UserRepository
	tags     *repository.TagRepository
	topics   *repository.TopicRepository
	circles  *repository.CircleRepository
	comments *repository.CommentRepository
	outbox   *repository.OutboxRepository
	idem     *repository.WorkerRepository
	notifs   *repository.NotificationRepository
	growth   *repository.GrowthRepository
}

// New 构建 Worker 运行器。sc 可为 nil（未启用 ES 时跳过索引）。
func New(cfg *config.Config, db *gorm.DB, cch *cache.Cache, sc *search.Client) *Runner {
	brokers := cfg.Kafka.Brokers
	return &Runner{
		cfg:      cfg,
		db:       db,
		cch:      cch,
		sc:       sc,
		producer: kafkax.NewProducer(brokers),
		consumer: kafkax.NewConsumer(brokers, cfg.Kafka.ConsumerGroup, topicsWithPrefix(cfg.Kafka.TopicPrefix)),
		posts:    repository.NewPostRepository(db),
		users:    repository.NewUserRepository(db),
		tags:     repository.NewTagRepository(db),
		topics:   repository.NewTopicRepository(db),
		circles:  repository.NewCircleRepository(db),
		comments: repository.NewCommentRepository(db),
		outbox:   repository.NewOutboxRepository(db),
		idem:     repository.NewWorkerRepository(db),
		notifs:   repository.NewNotificationRepository(db),
		growth:   repository.NewGrowthRepository(db),
	}
}

// topicsWithPrefix 生成需消费的全部业务 Topic（含前缀）。
func topicsWithPrefix(prefix string) []string {
	names := []string{event.TopicUser, event.TopicPost, event.TopicComment, event.TopicInteraction, event.TopicCircle}
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, prefix+"."+n)
	}
	return out
}

// Run 启动 Outbox Dispatcher 与消费循环，阻塞直到 ctx 取消。
func (r *Runner) Run(ctx context.Context) {
	// 预创建 Topic，确保消费组能拿到分区分配。
	if err := kafkax.EnsureTopics(r.cfg.Kafka.Brokers, topicsWithPrefix(r.cfg.Kafka.TopicPrefix)); err != nil {
		logger.Warnf("预创建 Topic 失败（可能已存在）: %v", err)
	}

	if r.sc != nil {
		if err := r.sc.EnsureIndices(ctx); err != nil {
			logger.Errorf("创建 ES 索引失败: %v", err)
		} else {
			r.reindexAll(ctx)
		}
	}

	go r.dispatchLoop(ctx)

	logger.Infof("worker 开始消费 topics=%v group=%s", topicsWithPrefix(r.cfg.Kafka.TopicPrefix), r.cfg.Kafka.ConsumerGroup)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, raw, err := r.consumer.Fetch(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Errorf("拉取消息失败: %v", err)
			time.Sleep(time.Second)
			continue
		}
		var m event.Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			logger.Errorf("解析事件失败: %v", err)
			_ = r.consumer.Commit(ctx, raw)
			continue
		}
		r.handle(ctx, &m)
		observability.KafkaConsume(msg.Topic, "success")
		_ = r.consumer.Commit(ctx, raw)
	}
}

// dispatchLoop Outbox Dispatcher：定时把待投递事件发送到 Kafka。
func (r *Runner) dispatchLoop(ctx context.Context) {
	interval := time.Duration(r.cfg.Worker.OutboxIntervalS) * time.Second
	if interval <= 0 {
		interval = 2 * time.Second
	}
	batch := r.cfg.Worker.BatchSize
	if batch <= 0 {
		batch = 100
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.dispatchOnce(ctx, batch)
		}
	}
}

// backoffSeconds 重试退避（秒）。
var backoffSeconds = []int{10, 30, 60, 300, 600}

func (r *Runner) dispatchOnce(ctx context.Context, batch int) {
	rows, err := r.outbox.FetchDispatchable(batch)
	if err != nil || len(rows) == 0 {
		return
	}
	for i := range rows {
		row := &rows[i]
		msg := event.Message{
			EventID: row.EventID, EventType: row.EventType, AggregateID: row.AggregateID,
			UserID: row.UserID, TraceID: row.TraceID, Payload: json.RawMessage(row.Payload), CreatedAt: row.CreatedAt,
		}
		b, _ := json.Marshal(msg)
		if err := r.producer.Publish(ctx, row.Topic, []byte(row.EventID), b); err != nil {
			observability.KafkaProduce(row.Topic, "fail")
			next := backoffSeconds[len(backoffSeconds)-1]
			if row.RetryCount < len(backoffSeconds) {
				next = backoffSeconds[row.RetryCount]
			}
			rc := row.RetryCount + 1
			if rc > r.cfg.Worker.MaxRetry {
				r.outbox.MarkFailed(row.ID)
				logger.Errorf("事件投递失败超过最大重试, eventId=%s", row.EventID)
			} else {
				r.outbox.MarkRetry(row.ID, rc, time.Now().Add(time.Duration(next)*time.Second))
				logger.Warnf("事件投递失败将重试, eventId=%s retry=%d", row.EventID, rc)
			}
			continue
		}
		r.outbox.MarkDispatched(row.ID)
		observability.KafkaProduce(row.Topic, "success")
		logger.Infof("事件已投递 topic=%s eventType=%s eventId=%s traceId=%s", row.Topic, row.EventType, row.EventID, row.TraceID)
	}
}

// Close 释放资源。
func (r *Runner) Close() {
	if r.consumer != nil {
		_ = r.consumer.Close()
	}
	if r.producer != nil {
		_ = r.producer.Close()
	}
}
