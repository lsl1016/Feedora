---
title: 关注模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: follow
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/follow_router.go
  - internal/api/follow_api.go
  - internal/dto/follow_dto.go
  - internal/service/follow_service.go
  - internal/repository/follow_repository.go
  - internal/model/user_follow.go
  - internal/cache/keys.go
summary: 用户关注关系的建立、取消、状态查询与列表，以及关注中心页面的动态、圈子、话题、标签聚合。
---

## 1. 模块概述

关注模块提供用户对用户的关注关系管理，包括关注、取消关注、关注状态查询、关注列表与粉丝列表。模块同时承担关注中心的聚合展示：`following-feed` 按 `all`/`user`/`circle`/`topic` 四个页签聚合帖子动态，`following-circles`/`following-topics`/`following-tags` 分别列出用户加入的圈子、参与过的话题和使用过的标签。圈子、话题、标签均不存在独立的关注关系表：圈子取自圈子成员关系，话题与标签由用户发帖行为推导。

## 2. 接口清单

路由注册于 `registerFollow`，实际前缀为 `/api/v1`，下表省略该前缀。

| 路径 | 方法 | 功能 | 控制器 |
| --- | --- | --- | --- |
| `/users/:userId/follow` | POST | 关注用户（需登录） | `FollowAPI.Follow` |
| `/users/:userId/follow` | DELETE | 取消关注用户（需登录） | `FollowAPI.Unfollow` |
| `/users/:userId/follow/state` | GET | 查询当前用户是否关注目标用户（公开） | `FollowAPI.State` |
| `/users/:userId/following` | GET | 获取指定用户关注的人列表（公开） | `FollowAPI.Following` |
| `/users/:userId/followers` | GET | 获取指定用户的粉丝列表（公开） | `FollowAPI.Followers` |
| `/users/me/following` | GET | 获取当前用户关注的人列表（需登录） | `FollowAPI.MyFollowing` |
| `/users/me/following-feed` | GET | 获取当前用户的关注动态（需登录） | `FollowAPI.MyFollowingFeed` |
| `/users/me/following-circles` | GET | 获取当前用户加入的圈子（需登录） | `FollowAPI.MyFollowingCircles` |
| `/users/me/following-topics` | GET | 获取当前用户参与的话题（需登录） | `FollowAPI.MyFollowingTopics` |
| `/users/me/following-tags` | GET | 获取当前用户使用过的标签（需登录） | `FollowAPI.MyFollowingTags` |

`/users/me/*` 为静态路由，注册顺序上先于 `:userId` 通配匹配，避免 `me` 命中 `UserIDURI` 的数字校验。`/users/:userId/follow/state` 未挂认证中间件，未登录访问不报错，返回 `following=false`。

## 3. 核心逻辑

### 关注与取关的幂等处理

`FollowService.Follow` 依次校验：发起方已登录（`followerID > 0`）、不能关注自己（返回 `ErrParams`）、目标用户存在（`UserRepository.FindByID` 为 `nil` 时返回 `ErrNotFound`）。落库由 `FollowRepository.Follow` 执行 `Create`，命中 `user_follows` 唯一键冲突时返回 `created=false`；冲突判定为双保险：`gorm.ErrDuplicatedKey` 或错误信息包含 `Duplicate entry` / `Error 1062`。

仅当 `created=true` 时执行后续动作，重复关注不重复触发：

- 双方计数增减：`IncColumn(followerID, "following_count", 1)` 与 `IncColumn(followeeID, "follower_count", 1)`，取关时方向取 `-1`，为原子增减，不回读重算。
- 失效双方资料缓存：删除 `user:profile:{followerID}` 与 `user:profile:{followeeID}`（`cache.UserProfileKey`）。
- 发布事件：`producer.Publish(event.TopicUser, event.UserFollowed, followeeID, followerID, nil)`，取关对应 `event.UserUnfollowed`。

`Unfollow` 按 `follower_id + followee_id` 物理删除，以 `RowsAffected > 0` 判断是否实际删除；关系不存在时返回 `removed=false`，同样幂等。`Unfollow` 不校验目标用户是否存在（`Follow` 校验，两者不对称）。

### 关注状态查询

`State` 在未登录（`viewerID <= 0`，该路由无认证中间件）或查看自己时直接返回 `false`，否则执行 `IsFollowing` 的 `COUNT` 查询。

### following-feed 聚合

采用拉模式：每次请求实时查询数据库，无预计算时间线，无推送扇出，无 Redis 缓存。四个页签统一走 `PostService.List`，构造 `ListFilter{ViewerID, Page, PageSize, Sort: "latest"}`，按页签区分：

| 页签 | 过滤来源 | 实现 | 条目 `feedType` |
| --- | --- | --- | --- |
| `all` / `user` | 关注的作者 | `filter.FeedType = "following"`，`PostRepository.List` 中以子查询 `author_id IN (SELECT followee_id FROM user_follows WHERE follower_id = ?)` 过滤 | `post` |
| `circle` | 加入的圈子 | 先 `JoinedCircles(viewerID, 0, 100)` 取前 100 个圈子，ID 集合作为 `CircleIDs` 传入帖子查询（圈内帖不受公开可见性限制） | `circle` |
| `topic` | 参与的话题 | 先 `ParticipatedTopics(viewerID, 0, 100)` 取前 100 个话题，ID 集合作为 `TopicIDs` 传入帖子查询 | `topic` |

查询结果映射为 `FollowingFeedItem`：`FeedID` 与 `TargetID` 均取帖子 ID，`TargetURL` 固定为 `/posts/{postID}`，`SourceName`/`SourceAvatar` 取作者昵称与头像。`feedType` 取当前页签的值，同一页返回的所有条目取值相同，不区分单条动态的实际来源。

### 粉丝/关注列表分页

`ListFollowing`/`ListFollowers` 以 `IN` 子查询取对端用户 ID，再查 `users` 表并过滤 `status = UserNormal`（非正常状态用户不出现在列表和计数中），按 `id DESC` 排序，`LIMIT/OFFSET` 偏移分页，`total` 为独立的 `COUNT` 查询。返回用户对象经 `dto.ToUser(&u, false, 0)` 转换，不含关注状态字段。

「关注的圈子」「关注的话题」「关注的标签」同为偏移分页：圈子来自 `circle_members`（`status = CircleMemberNormal`）；话题来自用户未删除帖子（`author_id = ? AND deleted_at IS NULL`）关联的 `post_topics`；标签来自同一帖子集合关联的 `post_tags`，排序为 `use_count DESC, id DESC`，其余列表均按 `id DESC`。

## 4. 数据模型

`user_follows` 表（`model.UserFollow`）：

| 字段 | 列 | 说明 |
| --- | --- | --- |
| `ID` | `id` | 主键 |
| `FollowerID` | `follower_id` | 关注发起方，唯一索引第 1 列 |
| `FolloweeID` | `followee_id` | 被关注方，唯一索引第 2 列 |
| `CreatedAt` | `created_at` | 关注时间 |

唯一索引 `uk_user_follow(follower_id, followee_id)` 保证同一方向的关注关系唯一。表无状态字段，语义为"存在即关注"：关注是插入一行，取关是物理删除该行，不存在软删除、待确认、屏蔽等中间状态。`(A, B)` 与 `(B, A)` 是两行独立记录，互不影响。

## 5. 配置项与限制

无新增配置项，模块不读取独立配置。代码内置的边界值：

| 项 | 值 | 说明 |
| --- | --- | --- |
| 分页默认页码 / 页大小 | `1` / `10` | `normalizePageRequest` 与 `normPage` 内置 |
| 分页页大小上限 | `100` | 超出时截断为 `100` |
| feed 聚合圈子/话题上限 | 各 `100` 个 | `circle`/`topic` 页签固定取前 100 个，超出部分不参与聚合 |

当前实现的其他边界：

- 圈子、话题、标签的"关注"不产生任何关注记录，无关注时间，取值为成员关系或发帖行为的实时推导。
- 关注关系不写入 Redis，Redis 中不存在关注集合或时间线缓存；关注变化仅删除双方 `user:profile:{id}` 资料缓存。
- `following-feed` 无推模式与预计算，关注数或帖子量大时每次请求均执行实时查询。
- `Unfollow` 不校验目标用户是否存在。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
