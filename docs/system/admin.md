---
title: 后台管理模块功能文档
date: 2026-09-09
version: v1.1
type: system
module: admin
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/admin_router.go
  - internal/api/admin_api.go
  - internal/dto/admin_dto.go
  - internal/service/admin_service.go
  - internal/repository/admin_repository.go
  - internal/model/admin.go
  - pkg/middleware/admin.go
  - pkg/middleware/auth.go
summary: 面向平台管理员的后台查询与治理接口，含仪表盘统计、标签与话题管理、用户封禁与帖子上下架（事件驱动 ES 同步）及操作日志读写。
---

## 1. 模块概述

后台管理模块（admin）挂载在 `/api/v1/admin` 路由组下，全部接口要求 JWT 登录且角色为 `admin`。模块提供用户、帖子、评论、圈子、操作日志的分页查询，仪表盘四类实体总量统计，标签与话题的创建和更新能力，以及用户封禁/解禁与帖子上下架两类治理写操作。帖子状态变更与话题创建/更新会发布 Kafka 事件，由 Worker 消费后同步 ES 索引。评论删除与用户角色变更未在本模块实现。

## 2. 接口清单

以下路径省略 `/api/v1` 前缀。

| 路径 | 方法 | 功能 | 控制器 |
| --- | --- | --- | --- |
| /admin/users | GET | 用户分页列表，按 `keyword` 模糊匹配昵称或账号 | `AdminAPI.Users` |
| /admin/users/:userId/status | PUT | 封禁 / 解禁用户（`normal` / `banned`） | `AdminAPI.UserStatus` |
| /admin/posts | GET | 帖子分页列表，可按 `status` 过滤 | `AdminAPI.Posts` |
| /admin/posts/:postId/status | PUT | 帖子上下架（`published` / `hidden` / `takedown`） | `AdminAPI.PostStatus` |
| /admin/comments | GET | 评论分页列表 | `AdminAPI.Comments` |
| /admin/tags | GET | 标签全量列表 | `AdminAPI.Tags` |
| /admin/tags | POST | 创建标签 | `AdminAPI.CreateTag` |
| /admin/tags/:tagId | PUT | 更新标签（部分字段） | `AdminAPI.UpdateTag` |
| /admin/topics | GET | 话题分页列表 | `AdminAPI.Topics` |
| /admin/topics | POST | 创建话题 | `AdminAPI.CreateTopic` |
| /admin/topics/:topicId | PUT | 更新话题（部分字段） | `AdminAPI.UpdateTopic` |
| /admin/circles | GET | 圈子分页列表 | `AdminAPI.Circles` |
| /admin/dashboard/stats | GET | 用户、帖子、评论、圈子四项总量统计 | `AdminAPI.Stats` |
| /admin/operation-logs | GET | 操作日志分页列表 | `AdminAPI.Logs` |

## 3. 核心逻辑

### 管理员鉴权

- 路由组以 `g := x.v1.Group("/admin", x.authMW, x.adminMW)` 注册，两个中间件按顺序执行，先认证后鉴权。
- `middleware.Auth`（`pkg/middleware/auth.go`）解析 `Authorization: Bearer <token>` 中的 JWT，解析失败返回 401（`errs.ErrUnauth`）并中止；解析成功将 `claims.UserID` 与 `claims.Role` 写入 gin 上下文。
- `middleware.Admin`（`pkg/middleware/admin.go`）读取上下文中的角色，不等于字符串 `"admin"` 即返回 403（`errs.ErrForbidden`）并中止。
- 角色来源于 JWT claims，签发时取自 `users.role` 字段（默认 `user`）。管理员账号由 `cmd/seed` 种子程序创建（账号 `admin`，`role = "admin"`）。角色校验不回查数据库，修改 `users.role` 对已签发的 token 不生效，需重新登录后新 token 才携带新角色。
- 角色只有 `admin` 一级，无更细粒度的权限划分；所有 admin 接口权限等价。

### 标签与话题管理

- 创建：`name` 为空返回参数错误；构造实体时 `Status` 固定为 `enabled`；落库依赖 `tags.name` / `topics.name` 上的唯一索引（`uk_tag_name`、`uk_topic_name`）防重，重名触发唯一冲突后由服务层统一映射为 409「标签已存在」/「话题已存在」。该映射不区分错误类型，创建时的其他数据库错误同样以 409 返回。
- 更新：请求体各字段为指针类型，仅非 nil 字段进入更新白名单——标签为 `name`、`description`、`status`；话题为 `name`、`description`、`is_official`、`is_recommended`、`cover_url`、`status`。更新成功后回读实体，回读不到则返回 404。`status` 入参未做枚举校验，可写入任意字符串。
- 话题创建/更新落库成功后向 `TopicTopic`（`feedora.topic.events`）发布 `TopicCreated` / `TopicUpdated` 事件（`user_id` 为 0），由 Worker 的 `handleSearch` 回查话题并维护 `feedora_topics` 增量索引。
- 标签列表接口不分页，一次性返回全量；话题列表分页。

### 用户封禁与帖子上下架

- `PUT /admin/users/:userId/status`：请求体 `AdminUserStatusRequest`（`status` 必填）。`SetUserStatus` 校验状态枚举（仅 `normal` / `banned`，其余返回 400）与用户存在性（不存在返回 404），更新 `users.status` 后写操作日志。被封禁用户登录被 auth 模块拦截（返回 10003）；对已签发的 JWT 不生效。
- `PUT /admin/posts/:postId/status`：请求体 `AdminPostStatusRequest`（`status` 必填）。`SetPostStatus` 校验状态枚举（仅 `published` / `hidden` / `takedown`，其余返回 400）与帖子存在性（不存在返回 20001），更新 `posts.status` 后写操作日志并按目标状态向 `TopicPost`（`feedora.post.events`）发布事件供 ES 索引同步：`published` → `PostUnhidden`、`hidden` → `PostHidden`、`takedown` → `PostUpdated`；事件 `user_id` 为管理员 ID。
- 两接口成功均返回 `{"updated": true}`。操作日志经私有 `log()` 写入 `operation_logs`，`AdminName` 回查管理员昵称，`action` 为 `update_user_status` / `update_post_status`，`detail` 记录目标昵称/标题与目标状态。
- 事件发布依赖 `kafka.enabled`：关闭时生产者为 Noop 实现，仅打印日志。

### 列表查询行为

- 用户列表复用 `UserRepository.List`：固定过滤 `status <> 'banned'`（被封禁用户在后台列表不可见，解禁需已知 userId 直接调用状态接口），`keyword` 对 `nickname`、`account` 做 LIKE 匹配，按 `id ASC` 排序；返回的 `dto.User` 不携带关注状态。
- 帖子列表由 `AdminRepository.ListPosts` 实现：`status` 为空或 `"all"` 时不过滤，否则按 `status = ?` 精确匹配，按 `id DESC` 排序，直接返回 `model.Post` 原始字段（DTO `AdminPostItem` 仅作 Swagger 声明，序列化字段名与模型导出字段一致）。
- 评论总数与列表分两次查询（`CommentRepository.ListAll` + `Count`）。
- 统计接口对 `users`、`posts`、`comments`、`circles` 四张表各执行一次全表 `Count`，无时间维度。

### 未实现的能力（否定事实）

- 无用户角色变更、用户详情接口；`SetUserStatus` 不限制目标账号角色，可将管理员账号置为 banned。
- 无评论的下架、恢复、删除接口；帖子仅支持状态枚举内的上下架，无物理删除。
- 无圈子管理写接口（创建、解散、推荐位设置均无）。
- `operation_logs` 的写入仅来自 `SetUserStatus` / `SetPostStatus` 两类操作，标签与话题管理不记操作日志；`AddLog` 写入失败被静默忽略。

## 4. 数据模型

模块不设独立 `admins` 表，管理员身份即 `users` 表的 `role` 字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `role` | varchar(32)，默认 `user` | 取值 `admin` 时通过 `Admin` 中间件校验 |

模块自有表 `operation_logs`（`model.OperationLog`）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int64 主键 | |
| `admin_id` | int64，带索引 | 操作管理员 ID |
| `admin_name` | varchar(64) | 管理员名称 |
| `action` | varchar(64) | 操作动作 |
| `target_type` | varchar(64) | 目标类型 |
| `target_id` | int64 | 目标 ID |
| `detail` | varchar(1000) | 详情描述 |
| `created_at` | time | 创建时间 |

跨模块操作：`AdminService` 除自有 `AdminRepository`（帖子列表、日志列表、统计、帖子/用户状态更新、日志写入）外，直接注入并调用 `UserRepository`、`TagRepository`、`TopicRepository`、`CircleRepository`、`CommentRepository` 完成查询与标签、话题写入，并注入 `event.Producer` 发布话题与帖子状态事件；不经过各业务模块的 Service 层。`operation_logs` 已注册在 `model/all.go` 的自动迁移列表中，`update_user_status` / `update_post_status` 两类操作产生写入。

## 5. 配置项与限制

无新增配置。鉴权依赖全局 JWT 配置（`middleware.Auth(opts.JWT)`），JWT 密钥与过期时间沿用应用级配置。

分页约束（两层归一化，`api/context.go` 与 `service/common.go` 逻辑一致）：`page` 缺省或非正取 1，`pageSize` 缺省或非正取 10，上限 100。

其他限制：

- 标签、话题创建的错误统一映射为 409，无法区分重名与其他数据库故障。
- 标签/话题更新的回读错误被忽略（`FindByID` 的 error 不检查），仅以 nil 判断 404。
- 统计接口为四次全表 Count，无缓存。
- 封禁对已签发的 JWT 不生效，token 有效期内被封禁用户仍可访问需认证接口。
- `takedown` 状态复用 `PostUpdated` 事件，ES 侧按帖子当前状态重建索引文档。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.1 | 2026-09-09 | Feedora 项目组 | 新增用户封禁与帖子上下架接口，操作日志开始写入 |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
