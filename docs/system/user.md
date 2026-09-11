---
title: 用户模块功能文档
date: 2026-09-09
version: v1.2
type: system
module: user
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/user_router.go
  - internal/api/user_api.go
  - internal/dto/user_dto.go
  - internal/service/user_service.go
  - internal/repository/user_repository.go
  - internal/repository/interaction_repository.go
  - internal/model/user.go
summary: 提供用户列表与详情查询、当前用户资料维护，以及我的帖子、评论、点赞、收藏四类个人内容聚合视图。
---

## 1. 模块概述

用户模块提供用户的公开查询能力（列表与详情）和当前用户的个人资料维护能力，并聚合"我的"视图：我的帖子、我的评论、我点赞的帖子、我收藏的帖子。模块本身不涉及注册、登录、密码等认证流程，用户记录由认证模块创建，本模块只读取和更新资料字段。用户详情查询采用 Cache Aside 缓存策略。

## 2. 接口清单

路由注册于 `internal/router/user_router.go`，全部挂载在 `/api/v1` 分组下；带 `authMW` 中间件的接口要求登录。下表路径省略 `/api/v1` 前缀。

| 路径 | 方法 | 功能 | 控制器 |
|------|------|------|--------|
| `/users` | GET | 分页获取用户列表，支持 `keyword` 关键词过滤昵称或账号 | `UserAPI.List` |
| `/users/:userId` | GET | 获取用户详情 | `UserAPI.Get` |
| `/users/me/profile` | PUT | 更新当前用户资料（昵称、头像、签名） | `UserAPI.UpdateProfile` |
| `/users/me/posts` | GET | 分页获取当前用户自己的帖子（含隐藏与草稿） | `UserAPI.MyPosts` |
| `/users/me/comments` | GET | 分页获取当前用户的评论 | `UserAPI.MyComments` |
| `/users/me/liked-posts` | GET | 分页获取当前用户点赞的帖子 | `UserAPI.MyLikedPosts` |
| `/users/me/favorite-posts` | GET | 分页获取当前用户收藏的帖子 | `UserAPI.MyFavoritePosts |

共 7 个接口。`/users` 与 `/users/:userId` 为公开接口，不加认证中间件；其余 5 个 `me` 系列接口经 `authMW` 认证，用户 ID 从 `middleware.CurrentUserID(c)` 获取，不信任请求参数。

## 3. 核心逻辑

### 3.1 资料更新规则

`UserService.UpdateProfile` 采用指针字段的部分更新语义：请求体 `dto.UpdateProfileRequest` 中 `Nickname`、`Avatar`、`Bio` 均为指针，仅非 `nil` 字段写入更新 map，未传字段保持原值；`updated_at` 总是被刷新。更新成功后先删除该用户的资料缓存，再通过 `Get` 重新查询返回，重查过程会重建缓存。service 层对三个字段的长度与格式不做校验，约束仅来自数据库列定义（`nickname` 64、`avatar` 512、`bio` 500）。

密码字段 `password_hash` 不在可更新范围内，本模块不提供任何密码读取或修改路径。

### 3.2 头像处理

头像仅作为字符串地址存储与更新（如 `/static/avatar.png`），本模块不含头像文件上传、裁剪或存储处理逻辑，前端直接提交头像地址字符串。

### 3.3 用户详情缓存

`UserService.Get` 使用 Cache Aside 模式：缓存键为 `user:profile:{id}`，命中则直接反序列化返回 `dto.User`；未命中时经 `UserRepository.FindByID` 查库，记录不存在（返回 `nil, nil`）映射为 `ErrNotFound`，查到后以 `userProfileTTL`（30 分钟）写入缓存。

### 3.4 用户列表过滤

`UserRepository.List` 固定排除 `status = banned` 的用户（`status <> ?`），`keyword` 非空时按 `nickname LIKE` 或 `account LIKE` 模糊匹配，按 `id ASC` 排序。响应的 `dto.User` 包含 `account` 登录账号字段，且该接口为公开接口。

### 3.5 我的帖子

`UserService.MyPosts` 委托 `PostService.List`，过滤条件为 `AuthorID = ViewerID = 当前用户`、`Status = "all"`、`IncludeHidden = true`，因此返回结果包含该用户处于隐藏状态和草稿状态的帖子。

### 3.6 点赞与收藏帖子

`MyLikedPosts` 与 `MyFavoritePosts` 共用 `postsByRelation` 私有方法：经 `InteractionRepository.PagePostIDsByUser(table, userId, offset, limit)` 在关系表（`post_likes` 或 `post_favorites`）上执行 DB 级分页——先 `COUNT` 取该用户关系总数，再按关系 `id` 倒序以 `LIMIT/OFFSET` 取当前页帖子 ID，最后经 `PostService.AssembleByIDs` 组装帖子 DTO。每次请求只加载当前页的关系记录，不再全量取 ID 后内存切页（原 `PostIDsByUser` 全量方法已删除）。

### 3.7 统计字段维护

users 表的计数字段（`post_count`、`comment_count`、`follower_count`、`following_count`、`like_count`、`point_count`）不由本模块接口修改，由其它模块通过 `UserRepository.IncColumn` 增量维护：`delta >= 0` 时执行 `column + ?` 自增，`delta < 0` 时执行 `GREATEST(column - ?, 0)` 自减并做非负保护；该方法使用 `UpdateColumn`，不触发 GORM 钩子、不刷新 `updated_at`。

### 3.8 DTO 组装规则

`dto.ToUser` 将模型转为完整信息时存在以下阶段一近似实现：

- `experience` 直接取 `point_count`，未单独累积经验值；
- `nextLevelExperience` 按公式 `level * 100` 计算，无等级配置表；
- `badgeCount` 恒为 `0`，勋章功能未落地；
- 本模块内 `Get` 与 `List` 调用 `ToUser` 时固定传入 `checkedInToday = false`、`continuousCheckInDays = 0`，因此这两个接口返回的签到字段恒为 `false` / `0`，实际签到状态不在此查询。

等级名称由 `levelName` 函数按等级阈值映射：`>= 10` 为"资深专家"，`>= 6` 为"活跃达人"，`>= 3` 为"进阶用户"，其余为"新手上路"；`LevelNameOf` 导出该映射供排行榜等其它模块复用。

## 4. 数据模型

users 表（`model.User`，表名 `users`）关键字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `account` | string(64) | 登录账号，唯一索引 `uk_account` |
| `password_hash` | string(255) | 密码哈希，本模块不读写 |
| `nickname` / `avatar` / `bio` | string | 资料字段，长度上限分别为 64 / 512 / 500 |
| `role` | string(32) | 用户角色，默认 `user` |
| `status` | string(32) | 用户状态，默认 `normal`，存在 `banned` 值（`model.UserBanned`），建索引 |
| `follower_count` / `following_count` | int64 | 粉丝数 / 关注数，默认 0 |
| `post_count` / `comment_count` / `like_count` | int64 | 发帖数 / 评论数 / 获赞数，默认 0 |
| `point_count` | int64 | 积分，默认 0 |
| `level` | int | 用户等级，默认 1 |
| `created_at` / `updated_at` | time | 创建 / 更新时间 |
| `deleted_at` | gorm.DeletedAt | 软删除标记，建索引 |

软删语义：`deleted_at` 使用 GORM 标准软删除，本模块所有查询与更新自动附加 `deleted_at IS NULL` 条件，软删用户对全部接口不可见。本模块不提供恢复软删用户或硬删除用户的接口。

## 5. 配置项与限制

无新增配置。用户资料缓存 TTL 由代码常量 `userProfileTTL`（`internal/service/user_service.go`，30 分钟）固定，不可通过配置文件调整。

当前实现的其他限制：

- 用户创建（注册）不在本模块，`UserRepository.Create`、`FindByAccount`、`CountByAccount` 供认证模块等外部调用，本模块路由不触达；
- 用户列表与详情为公开接口且响应包含 `account` 字段；
- `experience`、`nextLevelExperience`、`badgeCount`、签到字段为占位或近似值，见 3.8 节。

## 6. 历史版本

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.2 | 2026-09-09 | Feedora 项目组 | 我的点赞/收藏列表改为 DB 级分页（PagePostIDsByUser：COUNT + LIMIT/OFFSET 按关系 id 倒序），删除全量 ID 内存分页 |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
