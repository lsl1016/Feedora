---
title: 搜索模块功能文档
date: 2026-09-09
version: v1.1
type: system
module: search
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/search_router.go
  - internal/api/search_api.go
  - internal/dto/search_dto.go
  - internal/service/search_service.go
  - internal/search/client.go
  - pkg/esx/elasticsearch.go
  - internal/event/event_type.go
  - internal/event/topic.go
  - internal/worker/handlers.go
  - internal/worker/index.go
summary: 搜索模块基于 Elasticsearch 提供帖子、用户、圈子、话题四类对象的综合搜索、帖子标题联想与 Redis 热词榜，索引数据由业务事件（含话题创建/更新）经 Outbox 与 Kafka 异步同步，并在 Worker 启动时全量重建。
---

## 1. 模块概述

搜索模块基于 Elasticsearch 提供四类对象的全文检索：帖子、用户、圈子、话题。对外提供三个只读接口：综合搜索（`/search`，支持按类型筛选）、搜索联想（`/search/suggest`，仅帖子标题/摘要）与热门搜索词（`/search/hot-keywords`，Redis ZSet 统计）。索引数据不在业务请求链路中写入，而是由业务事件经 Outbox 表、Kafka 送达 Worker 异步写入 ES，Worker 每次启动时还会全量重建索引。评论不参与搜索。`internal/worker/search/` 目录仅为阶段占位，实际的索引同步逻辑位于 `internal/worker/handlers.go` 与 `internal/worker/index.go`。

## 2. 接口清单

路由注册于 `internal/router/search_router.go`，挂载于 `/api/v1` 前缀下，未挂认证中间件，均为匿名可访问。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/search` | GET | 综合搜索（支持 `keyword`、`type`、分页参数） | `SearchAPI.Search` |
| `/search/suggest` | GET | 搜索联想建议（帖子标题，最多 8 条） | `SearchAPI.Suggest` |
| `/search/hot-keywords` | GET | 热门搜索词 Top10 | `SearchAPI.HotKeywords` |

## 3. 核心逻辑

### 3.1 支持的搜索对象

| 对象 | ES 索引 | 检索字段（`^2` 为权重加成） | 过滤条件 |
|---|---|---|---|
| 帖子 | `feedora_posts` | `title^2`、`content`、`summary`、`tagNames` | `status.keyword=published`、`visibility.keyword=public` |
| 用户 | `feedora_users` | `nickname^2`、`bio` | `status.keyword=normal` |
| 圈子 | `feedora_circles` | `name^2`、`description`、`category` | `status.keyword=normal` |
| 话题 | `feedora_topics` | `name^2`、`description` | `status.keyword=enabled` |

评论没有对应的 ES 索引，也没有评论相关的事件处理分支，评论内容不可搜索。用户与圈子索引仅保存检索和展示所需的冗余字段。

### 3.2 索引 mapping 与分词器

索引由 `search.Client.EnsureIndices` 创建，创建时不带请求体，使用 ES 默认动态映射，未显式配置任何分词器或字段 mapping。服务层查询中使用的 `status.keyword`、`visibility.keyword` 依赖动态映射自动生成的 `keyword` 子字段做精确过滤。

### 3.3 索引同步链路

同步链路为：业务 Service 写库后发布事件 → `event.OutboxProducer` 写入 MySQL `event_outbox` 表（topic 为 `kafka.topicPrefix` + `.` + topic 名，默认前缀 `feedora`）→ Worker 的 `dispatchLoop` 按 `worker.outboxIntervalSeconds` 间隔批量取出 pending 记录投递 Kafka（失败按 10/30/60/300/600 秒退避重试，超过 `worker.maxRetry` 标记 failed）→ Worker 消费组（`kafka.consumerGroup`，默认 `feedora-worker`）消费 `feedora.post.events`、`feedora.user.events`、`feedora.circle.events`、`feedora.topic.events` 等 topic → `Runner.handleSearch` 以 `search` 为 worker 名做幂等 Claim 后从 MySQL 重新加载实体 → 调用 `IndexDoc`/`DeleteDoc` 写入 ES（每次写入均带 `refresh=true`）。

`handleSearch` 处理的事件类型：

| 事件类型 | 来源 topic | 处理 |
|---|---|---|
| `PostCreated`、`PostUpdated`、`PostHidden`、`PostDeleted`、`PostUnhidden` | `feedora.post.events` | `indexPostByID`：帖子不存在或状态为 `deleted` 时从索引删除文档，否则重建帖子文档（含作者昵称、标签名、话题名） |
| `UserRegistered` | `feedora.user.events` | `indexUser`：写入用户文档 |
| `CircleCreated` | `feedora.circle.events` | `indexCircle`：写入圈子文档 |
| `TopicCreated`、`TopicUpdated` | `feedora.topic.events` | `indexTopic`：按 `AggregateID` 回查话题，存在时重建话题文档，不存在时跳过（事件由后台 `CreateTopic` / `UpdateTopic` 成功后发布） |

`feedora.topic.events` 由 Worker 启动时 `EnsureTopics` 预建，与其他业务 topic 一致。话题删除不发布事件。

Worker 启动时（`Runner.Run`）先调用 `EnsureIndices` 确保四个索引存在，随后执行 `reindexAll` 全量重建：帖子取 `status <> deleted`，用户取 `status <> banned`，圈子取 `status = normal`，话题取 `status = enabled`，逐条重新写入索引。

### 3.4 综合搜索实现

`SearchService.Search(keyword, typ, page, size, viewerID)` 的行为：

- ES 客户端未启用（`sc == nil`）时直接返回 500 错误"搜索服务未启用"。
- `keyword` 去掉首尾空白后非空时，向 Redis ZSet `search:hot_keywords` 执行 `ZIncrBy +1` 记录热词。
- `type` 为空、`all`、`post`、`user`、`circle`、`topic` 时命中对应分支；`type` 为空或 `all` 时依次查询四类索引，否则只查指定类型。
- 每类索引构造 `multi_match` 查询（`search.MatchQuery`），叠加表 3-1 的 term 过滤；`keyword` 为空时退化为 `match_all`（过滤条件仍然生效），即空关键词请求返回的是过滤后的存量文档分页。
- `SearchIDs` 以 `_source: false` 只取命中文档的 `_id`，解析为 int64 后按 ID 顺序回源装配：帖子走 `PostService.AssembleByIDs`（带 viewerID，匿名访问时为 0），用户走 `UserRepository.FindByIDs` 后转 `dto.ToUser`，圈子走 `CircleRepository.FindByIDs` 后连作者信息转 `dto.ToCircle`，话题逐 ID 走 `TopicRepository.FindByID` 转成 `dto.ToTopic`。
- 单类索引查询出错时静默跳过（`err == nil` 才装配），该类型不贡献结果，错误不向上抛出。
- `type=all` 时四类各自独立应用相同的 `from/size` 分页，返回扁平混合列表，`total` 为四类命中数之和，不存在跨类型统一排序或归并；页面里每类最多出现 `size` 条。

### 3.5 联想实现

`SearchService.Suggest(keyword)` 仅查询 `feedora_posts` 索引（字段 `title^2`、`summary`，过滤 `published` + `public`），取前 8 条，按 ES 默认相关性排序，经 `PostRepository.FindByIDs` 回源后返回 `SuggestItem`（`type` 固定为 `post`，`targetUrl` 为 `/posts/{postId}`）。ES 客户端未启用、关键词为空或查询出错时均返回空列表。

### 3.6 热词实现

`SearchService.HotKeywords()` 对 Redis ZSet `search:hot_keywords` 执行 `ZRevRangeWithScores` 取分数最高的前 10 个关键词返回字符串数组。热词数据完全来源于综合搜索接口的 `ZIncrBy` 计数，无过期时间、无清理任务；Redis 未启用时返回空数组。

### 3.7 高亮与分页

当前实现不做搜索结果高亮：`SearchIDs` 关闭 `_source` 且只解析 `_id`，响应中不返回命中的字段片段，也无 `highlight` DSL。排序未显式指定，采用 ES 默认的 `_score` 相关性排序。分页参数经 `normPage` 规整：`page` 默认 1，`pageSize` 默认 10、上限 100，`from = (page-1) * size`。

### 3.8 当前实现的局限

- 话题创建与更新（后台接口）经 `TopicCreated` / `TopicUpdated` 事件增量同步到索引；话题删除不发布事件，已删话题的文档保留在索引中，直到下次 Worker 启动全量重建。
- 用户文档仅在 `UserRegistered` 事件时写入，昵称、简介等资料修改不产生事件，索引中的用户信息会过期。
- 圈子文档仅在 `CircleCreated` 事件时写入，圈子资料修改、成员数与帖子数变化不同步。
- 帖子文档中的 `likeCount`、`commentCount`、`favoriteCount`、`hotScore` 仅在帖子自身的增删改事件触发重建时刷新，点赞、评论、收藏事件不触发帖子重建，计数在索引中长期滞后。
- 索引使用默认动态映射与默认分词器，未配置中文分词器，中文检索效果依赖 ES 默认行为。
- 每次 Worker 重启都对全部数据执行全量重建，且每条文档写入均带 `refresh=true`。
- `type=all` 的分页按类型各自独立执行，`total` 为各类型之和，翻页时各类型条数不均衡。
- 单类型 ES 查询失败被静默吞掉，接口仍返回 200，缺少降级提示。

## 4. 索引结构

索引名由 `elasticsearch.indexPrefix`（默认 `feedora`）拼接生成，文档 `_id` 为对应业务对象 ID 的十进制字符串。

### feedora_posts

| 字段 | 来源 | 说明 |
|---|---|---|
| `postId` | `posts.id` | 帖子 ID |
| `title` / `content` / `summary` | 帖子字段 | 标题、正文（`content_md`）、摘要 |
| `authorId` / `authorName` | 作者 | 作者 ID 与昵称（写入时实时查库） |
| `tagNames` / `topicNames` | 关联查询 | 帖子标签名与话题名数组 |
| `circleId` | `posts.circle_id` | 所属圈子 ID，无则为 0 |
| `visibility` / `status` | 帖子字段 | 可见性与状态 |
| `likeCount` / `commentCount` / `favoriteCount` / `hotScore` | 帖子字段 | 计数与热度 |
| `createdAt` | `posts.created_at` | 创建时间 |

### feedora_users

| 字段 | 说明 |
|---|---|
| `userId` / `nickname` / `bio` / `avatar` | 用户基础信息 |
| `status` | 用户状态 |
| `followerCount` / `postCount` | 粉丝数与帖子数 |
| `createdAt` | 创建时间 |

### feedora_circles

| 字段 | 说明 |
|---|---|
| `circleId` / `name` / `description` / `avatar` | 圈子基础信息 |
| `category` / `status` | 分类与状态 |
| `memberCount` / `postCount` / `featuredCount` | 成员数、帖子数、精华数 |
| `isRecommended` | 是否推荐 |
| `createdAt` | 创建时间 |

### feedora_topics

| 字段 | 说明 |
|---|---|
| `topicId` / `name` / `description` / `coverUrl` | 话题基础信息 |
| `isOfficial` / `isRecommended` | 官方与推荐标记 |
| `participantCount` / `postCount` / `status` | 参与数、帖子数、状态 |
| `createdAt` | 创建时间 |

## 5. 配置项

`configs/config.yaml` 中与本模块相关的配置：

| 配置键 | 默认值 | 说明 |
|---|---|---|
| `elasticsearch.enabled` | `true` | 是否启用 ES；为 `false` 时 API 进程不创建 `search.Client` |
| `elasticsearch.addresses` | `http://localhost:9200` | ES 地址列表 |
| `elasticsearch.username` / `elasticsearch.password` | 空 | ES 认证信息 |
| `elasticsearch.indexPrefix` | `feedora` | 索引名前缀 |
| `kafka.enabled` / `kafka.brokers` / `kafka.topicPrefix` / `kafka.consumerGroup` | `true` / `localhost:9092` / `feedora` / `feedora-worker` | 事件投递与消费，决定索引同步是否运转 |
| `worker.batchSize` / `worker.outboxIntervalSeconds` / `worker.maxRetry` | 100 / 2 / 5 | Outbox 投递批量、间隔与最大重试 |
| `redis.enabled` | `true` | 热词统计与查询依赖 Redis |

降级行为：

- API 进程启动时若 `elasticsearch.enabled=true` 但连接失败，进程启动失败退出；`enabled=false` 时 `search.Client` 为 `nil`，`/search` 返回 500"搜索服务未启用"，`/search/suggest` 返回空数组，`/search/hot-keywords` 仅依赖 Redis，不受影响。
- Worker 进程中 ES 连接失败仅记录日志，`sc` 保持 `nil` 继续运行，跳过全部索引写入。
- `kafka.enabled=false` 时 API 侧使用 `NoopProducer`（事件仅打印日志不投递），Worker 入口直接退出，索引不存在任何增量同步；同时搜索接口仍可查询既有索引内容。
- Redis 未启用时热词不记录、`/search/hot-keywords` 返回空数组，搜索与联想功能不受影响。

## 6. 历史版本

| 版本 | 日期 | 维护者 | 说明 |
|---|---|---|---|
| v1.1 | 2026-09-09 | Feedora 项目组 | 新增 TopicCreated/TopicUpdated 事件与 topic.events 订阅，话题索引增量同步 |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
