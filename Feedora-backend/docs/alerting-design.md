# 告警设计

基于 Prometheus 指标（见 `metrics-design.md`）配置告警规则，经 Alertmanager 分级通知。

## 1. 告警目标

及时发现：接口异常、服务不可用、Kafka 堆积、Worker 消费失败、MySQL/Redis/ES 异常、错误率升高、关键业务量异常。

## 2. 告警级别

| 级别 | 含义 | 响应 |
|---|---|---|
| P0 | 严重故障，核心不可用 | 立即电话/IM，7×24 |
| P1 | 重要故障，部分不可用或错误率高 | 15 分钟内响应 |
| P2 | 一般异常，性能下降 | 工作时间处理 |
| P3 | 观察类，趋势预警 | 记录跟进 |

## 3. API 告警规则

| 规则 | 表达式（示意） | 级别 |
|---|---|---|
| 5xx 错误率高 | `sum(rate(feedora_api_request_total{status=~"5.."}[5m])) / sum(rate(feedora_api_request_total[5m])) > 0.05` | P1 |
| P95 延迟高 | `histogram_quantile(0.95, sum(rate(feedora_api_request_duration_seconds_bucket[5m])) by (le)) > 1` | P2 |
| 服务无请求（疑似挂） | `sum(rate(feedora_api_request_total[5m])) == 0` 持续 5m | P0 |
| 登录接口错误率异常 | `route="/auth/login"` 且错误率 > 30% | P1 |
| 发帖接口错误率异常 | `route="/posts"` `method=POST` 错误率 > 10% | P1 |

## 4. Worker 告警规则

| 规则 | 示意 | 级别 |
|---|---|---|
| 消费失败率高 | `sum(rate(feedora_worker_process_total{result="fail"}[5m])) / sum(rate(feedora_worker_process_total[5m])) > 0.05` | P1 |
| Kafka Lag 超阈值 | `max(kafka_consumergroup_lag{group="feedora-worker"}) > 1000` 持续 10m | P1 |
| Outbox pending 持续增长 | `delta(feedora_outbox_pending[15m]) > 100` | P2 |
| Outbox failed > 0 | `feedora_outbox_failed > 0` | P1 |

## 5. MySQL 告警规则

| 规则 | 示意 | 级别 |
|---|---|---|
| 慢查询激增 | `increase(feedora_mysql_slow_total[5m]) > 50` | P2 |
| 连接池耗尽 | `feedora_mysql_conns_inuse / feedora_mysql_conns_open > 0.9` 持续 5m | P1 |
| SQL 错误率升高 | `rate(feedora_mysql_error_total[5m]) > 1` | P1 |

## 6. Redis 告警规则

| 规则 | 示意 | 级别 |
|---|---|---|
| Redis 不可用 | `redis_up == 0`（redis-exporter） | P0 |
| 命令延迟高 | P95 `feedora_redis_command_duration_seconds` > 0.05 持续 10m | P2 |
| 命中率过低 | `hit/(hit+miss) < 0.7` 持续 30m | P3 |
| 内存使用率高 | `redis_memory_used_bytes / redis_memory_max_bytes > 0.85` | P2 |

## 7. ES 告警规则

| 规则 | 示意 | 级别 |
|---|---|---|
| 查询错误率升高 | `rate(feedora_es_query_error_total[5m]) > 1` | P1 |
| 查询 P95 过高 | `histogram_quantile(0.95, ...es_query_duration...) > 0.5` | P2 |
| 集群 yellow/red | `elasticsearch_cluster_health_status != green` | P1(yellow)/P0(red) |
| 索引写入失败 | `rate(feedora_es_index_error_total[5m]) > 0` | P1 |

## 8. 业务告警

| 规则 | 示意 | 级别 |
|---|---|---|
| 发帖量归零 | `sum(increase(feedora_business_action_total{action="post_create"}[30m])) == 0`（业务时段） | P1 |
| 登录量骤降 | 环比昨日同时段下降 > 70% | P2 |
| 搜索失败率高 | 搜索接口错误率 > 20% | P2 |
| 上传失败率高 | `file_upload` 失败率 > 20% | P2 |

## 9. 告警处理建议（排查步骤）

- **API 5xx 高**：Grafana 看 Top 慢/错路由 → Loki 按 traceId 查错误日志 → 定位 Service/Repository/依赖。
- **Worker 失败率高**：看 worker 日志 `errorMessage` → 是否某类事件失败 → 检查 ES/DB/Redis 依赖 → 必要时暂停投递。
- **Kafka Lag 高**：确认 Worker 是否存活、处理耗时是否升高 → 扩 Worker 或排查慢依赖。
- **Outbox pending 增长**：Dispatcher 是否存活、Kafka 是否可达 → 看 `outbox dispatch` 日志。
- **MySQL 连接池耗尽**：排查慢 SQL、是否有未释放连接、是否需扩连接数。
- **Redis 不可用**：立即确认实例状态；缓存降级逻辑（本项目 cache 层 nil-safe）应保证核心可用。
- **ES red**：检查磁盘/分片；搜索可临时降级为 DB（榜单已具备 DB 回退）。
