# Feedora 后端（阶段一 MVP）

开发者知识社区 V2.1 阶段一后端，基于 Gin + GORM + MySQL + JWT，跑通「注册登录 → 浏览 → 发帖 → 评论 → 点赞收藏 → 圈子 → 文件上传 → 后台基础管理」核心闭环。

采用 router / api / service / repository / dto / model 分层架构。Kafka 与 OSS 采用可插拔接口：阶段一分别使用 noop（打印日志）与本地磁盘实现，后续可无侵入替换为真实 Kafka / MinIO / 阿里云 OSS。

## 技术栈

- 语言 / 框架：Go 1.24 + Gin
- ORM / 数据库：GORM + MySQL 8.0
- 鉴权：JWT（HS256）+ bcrypt
- 配置：YAML
- 文件存储：本地磁盘（模拟 OSS）
- 事件：noop 事件生产者（预留 Kafka）

## 分层架构

```text
请求 → router（路由注册）→ api（HTTP 处理）→ service（业务逻辑）→ repository（数据访问）→ MySQL
                                              ↓
                                          event（领域事件，阶段二接 Kafka）
```

- `internal/router`：总路由入口 + 各模块 `*_router.go`
- `internal/api`：HTTP 处理器，解析参数、返回统一响应
- `internal/service`：业务逻辑、权限校验、事务编排、事件发送
- `internal/repository`：GORM 数据访问（唯一操作数据库的层）
- `internal/dto`：与前端 `types.ts` 对应的请求 / 响应结构与转换器
- `internal/model`：GORM 模型（按域拆分，`all.go` 统一注册 AutoMigrate）
- `internal/event`：统一领域事件结构与生产者
- `internal/app`：应用装配（组装 DB、存储、事件、各层依赖）

## 目录结构

```text
Feedora-backend
├── cmd
│   ├── api/main.go        # API 服务入口
│   ├── worker/main.go     # Kafka Worker 入口（阶段二）
│   ├── migrate/main.go    # 数据库迁移入口
│   └── seed/main.go       # 种子数据
├── configs/config.yaml
├── internal
│   ├── app router api service repository dto model event
│   ├── cache search domain          # 阶段二占位
│   └── worker/{dispatcher,consumer,search,notification,growth,rank,stat,counter}
├── pkg
│   ├── config database jwtx ossx errors logger middleware response utils
│   └── redisx esx kafkax validator singleflightx   # 阶段二占位
├── migrations                       # 版本化表结构变更 SQL（scripts/apply-migrations.sh 执行）
├── deployments scripts               # 部署与常用脚本（建新表由 GORM AutoMigrate 完成）
└── uploads                          # 本地上传目录（运行时生成）
```

## 快速开始

### 1. 准备数据库

需要一个可用的 MySQL 8.x，并创建数据库（默认 `feedora`）：

```sql
CREATE DATABASE feedora CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

在 `configs/config.yaml` 中确认 `mysql.dsn`。默认：

```text
root:root@tcp(localhost:3306)/feedora?charset=utf8mb4&parseTime=True&loc=Local
```

### 2. 写入种子数据（可选）

```bash
go run ./cmd/seed
```

默认账号：`admin` / `zhangsan` / `lisi` / `wangwu`，密码均为 `123456`（admin 为管理员）。

### 3. 启动服务

```bash
go run ./cmd/api
```

服务默认监听 `:8090`，接口前缀 `/api/v1`，启动时自动执行 AutoMigrate 建表。

### 4. 生成 Swagger 接口文档

如未安装 `swag`，先执行：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

在 `Feedora-backend` 目录下执行：

```bash
swag init -g cmd/api/main.go -o docs
```

生成后的文档输出到 `docs/` 目录，启动服务后可通过 `http://localhost:8090/swagger/index.html` 访问。

### 5. 前端联调

`Feedora-frontend/.env.development`：

```text
VITE_API_MODE=real
VITE_API_BASE_URL=http://localhost:8090/api/v1
```

## 统一响应

成功：`{ "code": 0, "message": "success", "data": {}, "traceId": "trace-xxx", "timestamp": "..." }`

分页 `data`：`{ "list": [], "total": 0, "page": 1, "pageSize": 10 }`

业务错误码返回 HTTP 200 + 非 0 code；401/403/404/422 等映射到对应 HTTP 状态码。

## 已实现接口

认证：`POST /auth/register|login`、`GET /auth/me`、`POST /auth/logout`

用户：`GET /users`、`GET /users/:userId`、`PUT /users/me/profile`、`GET /users/me/posts|comments|liked-posts|favorite-posts`

帖子：`GET /posts`、`GET /posts/:id`、`POST /posts`、`PUT /posts/:id`、`DELETE /posts/:id`、`PUT /posts/:id/hide|unhide`、`POST /posts/:id/share|repost`

互动：`POST|DELETE /posts/:id/like`、`POST|DELETE /posts/:id/favorite`

评论：`GET /posts/:id/comments`、`POST /comments`、`POST /comments/:id/replies`、`DELETE /comments/:id`、`POST|DELETE /comments/:id/like`

标签 / 话题：`GET /tags`、`GET /tags/:id`、`GET /tags/:id/posts`、`GET /topics`、`GET /topics/:id`、`GET /topics/:id/posts`

圈子：`GET /circles`、`POST /circles`、`GET /circles/:id`、`POST /circles/:id/join|leave`、`GET /circles/:id/members|posts`、成员角色/禁言/移除管理

文件：`POST /files/upload`（multipart，字段 `file`、`bizType`）

后台（需 admin）：`GET /admin/users|posts|comments|tags|topics|circles`、`POST|PUT /admin/tags`、`POST|PUT /admin/topics`、`GET /admin/dashboard/stats`、`GET /admin/operation-logs`

## 阶段二：Redis + Kafka Worker + Elasticsearch

阶段二在阶段一基础上增加缓存加速、事件驱动与全文搜索。

### 基础设施（docker-compose）

`deployments/docker-compose.yml` 基于已有镜像启动中间件（API/Worker 在本地跑）：

| 服务 | 镜像 | 端口 | 说明 |
|---|---|---|---|
| mysql | mysql:8.0 | 3307→3306 | 宿主机 3306 已被占用，故映射 3307 |
| redis | redis:7-alpine | 6379 | 缓存 / 计数 / 榜单 ZSet / 未读数 |
| kafka | apache/kafka:3.7.2 | 9092 | KRaft 单节点事件总线 |
| elasticsearch | groupflow-elasticsearch:latest | 9200 | 帖子/用户/圈子/话题全文搜索 |
| ~~minio~~ | minio/minio | 9000/9001 | 当前环境无法拉取，OSS 暂用本地磁盘（已在 compose 中注释，接口兼容可随时切换） |

```bash
cd deployments && docker compose up -d      # 启动中间件
cd .. && go run ./cmd/seed                   # 首次写入种子数据
go run ./cmd/api                             # 启动 API（:8090）
go run ./cmd/worker                          # 启动 Worker（Outbox 投递 + 事件消费）
```

### 数据流

```text
API 写业务表 + event_outbox（同一入口）
  → Outbox Dispatcher 扫描 pending 投递 Kafka（失败退避重试）
    → Worker 消费（消费组 + worker_event_records/唯一索引 幂等）
      ├─ Search Worker：写 Elasticsearch 索引（隐藏/删除自动剔除）
      ├─ Notification Worker：生成站内通知 + Redis 未读数
      ├─ Growth Worker：发放积分流水 + 累加用户积分
      └─ Rank Worker：更新 Redis 榜单 ZSet
```

Worker 启动时会 `EnsureTopics` 预建 Kafka Topic，并对现有数据做一次 ES 全量重建。

### 新增接口

搜索：`GET /search`、`GET /search/suggest`、`GET /search/hot-keywords`

榜单：`GET /hot/ranks`（post/circle/topic，Redis ZSet 为空时回退 DB 排序）

通知：`GET /notifications`、`GET /notifications/unread-count`、`PUT /notifications/read-all`、`PUT /notifications/:id/read`

成长：`POST /growth/check-in`、`GET /growth/tasks`、`POST /growth/tasks/:id/claim`、`GET /growth/rankings`

### 缓存策略（Cache Aside）

帖子详情、用户资料、圈子详情走 Redis 缓存；帖子/圈子详情主体缓存后按当前用户叠加 liked/favorited、成员身份等 per-viewer 状态；写操作（编辑/隐藏/删除/点赞/收藏/评论/加入退出）失效对应缓存。未启用 Redis 时自动降级为直查 DB。

### 配置开关

`configs/config.yaml` 中 `redis.enabled` / `elasticsearch.enabled` / `kafka.enabled` 均可独立开关；关闭 Kafka 时事件生产者退回 Noop（仅打印日志），系统仍可运行。


