# 指标监控设计

使用 Prometheus client_golang，Gin 暴露 `GET /metrics`，Grafana 可视化。指标 label 保持低基数（禁止 userId、traceId 作为 label）。

## 1. 指标命名规范

统一前缀 `community_`（或 `feedora_`，本项目采用 `feedora_`）：

```
feedora_api_request_total
feedora_api_request_duration_seconds
feedora_mysql_query_duration_seconds
feedora_redis_command_duration_seconds
feedora_kafka_produce_total
feedora_kafka_consume_total
feedora_worker_process_total
feedora_es_query_duration_seconds
feedora_business_action_total
```

## 2. API 指标

- `feedora_api_request_total{method,route,status}`：请求总数（Counter）。
- `feedora_api_request_duration_seconds{method,route}`：耗时（Histogram，buckets: 5ms~2s）。
- 错误数：由 `status` label（4xx/5xx）聚合得出。
- 路由维度使用**注册路由模板**（`/posts/:postId`）而非真实路径，避免高基数。
- 不使用 userId 作为 label。

## 3. MySQL 指标

- `feedora_mysql_query_duration_seconds{op}`：查询耗时（op=select/insert/update/delete）。
- 慢查询数：`feedora_mysql_slow_total`（超 200ms +1）。
- 连接池：由 `sql.DBStats` 导出 `feedora_mysql_conns_open`、`feedora_mysql_conns_inuse`、`feedora_mysql_conns_idle`（Gauge，定时刷新）。
- 错误数：`feedora_mysql_error_total`。

## 4. Redis 指标

- `feedora_redis_command_duration_seconds{cmd}`：命令耗时。
- `feedora_redis_error_total`：命令错误数。
- `feedora_cache_hit_total{cache}` / `feedora_cache_miss_total{cache}`：命中/未命中（cache=post_detail/user_profile/circle_detail）。
- 热点 Key：建议由 Redis 侧 `--hotkeys` 或 bigkeys 巡检，不在应用打点。

## 5. Kafka 指标

- `feedora_kafka_produce_total{topic,result}`：生产成功/失败。
- `feedora_kafka_consume_total{topic,result}`：消费成功/失败。
- `feedora_kafka_consume_duration_seconds{topic}`：消费耗时。
- 消费积压 Lag：由 kafka-exporter 采集 `kafka_consumergroup_lag`。
- `feedora_outbox_pending`（Gauge）、`feedora_outbox_failed`（Gauge）：定时查库导出。

## 6. Worker 指标

- `feedora_worker_process_total{worker,result}`：每个 worker（search/notification/growth/rank）成功/失败数。
- `feedora_worker_retry_total{worker}`：重试数。
- `feedora_worker_process_duration_seconds{worker}`：处理耗时。
- `feedora_worker_idempotent_skip_total{worker}`：幂等命中跳过数。

## 7. ES 指标

- `feedora_es_query_duration_seconds{index}`：查询耗时。
- `feedora_es_query_error_total{index}`：查询错误。
- `feedora_es_index_duration_seconds{index}`：写入耗时。
- `feedora_es_index_error_total{index}`：写入错误。

## 8. 业务指标

统一用一个多维 Counter `feedora_business_action_total{action}`，`action` 取：

```
user_register / user_login / post_create / comment_create /
post_like / post_favorite / circle_create / search / file_upload / notification_create
```

在对应 service 成功路径 `+1`。

## 9. Prometheus 暴露接口

Gin 挂载：

```go
r.GET("/metrics", gin.WrapH(observability.MetricsHandler()))
```

`/metrics` 不走鉴权中间件，仅对内网/抓取端开放（生产用网络策略限制）。

Prometheus 抓取（本地，API 在宿主机 8090）：

```yaml
scrape_configs:
  - job_name: feedora-api
    static_configs:
      - targets: ['host.docker.internal:8090']
```

## 10. Grafana 面板设计

- **API 总览**：QPS（`rate(feedora_api_request_total[1m])`）、错误率（`status=~"5.."` 占比）、P95（`histogram_quantile(0.95, rate(feedora_api_request_duration_seconds_bucket[5m]))`）、Top 慢路由。
- **MySQL 总览**：查询 P95、慢查询数、连接池使用。
- **Redis 总览**：命令 P95、命中率（hit/(hit+miss)）、错误数。
- **Kafka / Worker 总览**：生产/消费成功失败、消费 P95、Lag、Outbox pending/failed、各 worker 处理量与重试。
- **ES 搜索总览**：查询 P95、错误率、写入耗时。
- **业务增长总览**：注册/登录/发帖/评论/点赞/搜索/上传的按天增长。
