// Package worker 阶段二占位：Kafka Worker 总入口。
// 阶段二在此统一启动 Outbox Dispatcher 与各消费者：
//   - dispatcher：Outbox 可靠投递 Kafka
//   - search：同步 ES 索引
//   - notification：生成站内通知
//   - growth：发放积分、更新成长数据
//   - rank：更新热门榜单、排行榜
//   - stat：统计数据处理
//   - counter：Redis 计数回写 MySQL
package worker
