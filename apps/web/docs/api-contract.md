# 前后端接口契约

接口前缀：`/api/v1`

统一响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "traceId": "trace-id",
  "timestamp": "2026-07-05T10:00:00+08:00"
}
```

分页响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "pageSize": 10
  }
}
```

## Auth

### 登录

`POST /auth/login`

```json
{
  "account": "zhangsan",
  "password": "123456"
}
```

响应：

```json
{
  "token": "jwt-token",
  "user": {}
}
```

### 注册

`POST /auth/register`

```json
{
  "account": "lisi",
  "nickname": "李四",
  "password": "123456",
  "confirmPassword": "123456"
}
```

### 当前用户

`GET /auth/me`

## User

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/users/{userId}` | 用户主页 |
| PUT | `/users/me/profile` | 修改资料 |
| GET | `/users/me/posts` | 我的帖子 |
| GET | `/users/me/comments` | 我的评论 |
| GET | `/users/me/liked-posts` | 我的点赞 |
| GET | `/users/me/favorite-posts` | 我的收藏 |

## Post

### 获取帖子列表

`GET /posts`

Query：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | number | 是 | 页码 |
| pageSize | number | 是 | 每页数量 |
| feedType | string | 否 | recommend/latest/following/hot |
| timeRange | string | 否 | all/today/week/month |
| tagId | number | 否 | 标签 ID |
| circleId | number | 否 | 圈子 ID |
| topicId | number | 否 | 话题 ID |
| keyword | string | 否 | 关键词 |

### 帖子详情

`GET /posts/{postId}`

### 创建帖子

`POST /posts`

```json
{
  "title": "Go 项目部署复盘",
  "content": "正文内容",
  "images": [],
  "tagIds": [1, 4],
  "topicIds": [2],
  "circleId": 1,
  "visibility": "public",
  "publishMode": "now",
  "scheduledAt": null
}
```

### 隐藏 / 取消隐藏

| 方法 | 路径 |
|---|---|
| PUT | `/posts/{postId}/hide` |
| PUT | `/posts/{postId}/unhide` |

### 互动

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/posts/{postId}/like` | 点赞 |
| POST | `/posts/{postId}/favorite` | 收藏 |
| POST | `/posts/{postId}/share` | 分享 |
| POST | `/posts/{postId}/repost` | 转发 |

## Comment

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/posts/{postId}/comments` | 评论列表 |
| POST | `/comments` | 发表评论 |
| DELETE | `/comments/{commentId}` | 删除评论 |

## Tag

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/tags` | 标签列表 |
| GET | `/tags/{tagId}` | 标签详情 |
| GET | `/tags/{tagId}/posts` | 标签内容聚合 |

## Circle

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/circles` | 圈子列表，支持 `scope=all/recommended/joined/created` |
| POST | `/circles` | 创建圈子 |
| GET | `/circles/{circleId}` | 圈子详情 |
| POST | `/circles/{circleId}/join` | 加入或申请加入 |
| GET | `/circles/{circleId}/members` | 圈子成员 |
| PUT | `/circles/{circleId}/members/{userId}/role` | 修改角色 |
| PUT | `/circles/{circleId}/members/{userId}/mute` | 禁言 |
| PUT | `/circles/{circleId}/members/{userId}/unmute` | 解禁 |
| DELETE | `/circles/{circleId}/members/{userId}` | 踢出 |

### 圈子聊天

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/circles/{circleId}/chat/messages` | 聊天消息 |
| POST | `/circles/{circleId}/chat/messages` | 发送消息 |

### 圈子智能摘要

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/circles/{circleId}/ai-summary/subscription` | 获取配置 |
| PUT | `/circles/{circleId}/ai-summary/subscription` | 保存配置 |
| PUT | `/circles/{circleId}/ai-summary/subscription/cancel` | 取消订阅 |
| GET | `/users/me/circle-summary-subscriptions` | 我的摘要订阅 |

## Topic

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/topics` | 话题列表 |
| GET | `/topics/{topicId}` | 话题详情 |
| GET | `/topics/{topicId}/posts` | 话题帖子 |

## Search

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/search` | 搜索结果 |
| GET | `/search/suggest` | 搜索联想 |
| GET | `/search/hot-keywords` | 热门搜索 |

## Hot

`GET /hot/ranks`

Query：

| 参数 | 类型 | 说明 |
|---|---|---|
| rankType | string | post/circle/topic |
| timeRange | string | today/week/all |

## Growth

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/growth/check-in` | 签到 |
| GET | `/growth/tasks` | 任务列表 |
| POST | `/growth/tasks/{taskId}/claim` | 领取奖励 |
| GET | `/growth/rankings` | 排行榜 |

## Activity

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/activities` | 活动列表 |
| GET | `/activities/{activityId}` | 活动详情 |
| POST | `/activities/{activityId}/vote` | 投票 |
| POST | `/activities/{activityId}/check-in` | 打卡 |

## Workspace

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/workspace/dashboard` | 工作空间首页统计 |
| GET | `/workspace/notes` | 笔记列表 |
| POST | `/workspace/notes` | 新建笔记 |
| GET | `/workspace/notes/{noteId}` | 笔记详情 |
| PUT | `/workspace/notes/{noteId}` | 更新笔记 |
| DELETE | `/workspace/notes/{noteId}` | 删除笔记 |
| POST | `/workspace/notes/{noteId}/publish-as-post` | 发布为帖子 |
| GET | `/workspace/knowledge-bases` | 知识库列表 |
| POST | `/workspace/knowledge-bases` | 新建知识库 |
| GET | `/workspace/knowledge-bases/{id}` | 知识库详情 |
| POST | `/workspace/knowledge-bases/{id}/invite` | 邀请成员 |

## AI

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/ai/model-config` | 用户模型配置 |
| PUT | `/ai/model-config` | 保存配置 |
| POST | `/ai/model-config/test` | 测试连接 |
| GET | `/ai/chat/sessions` | 会话列表 |
| POST | `/ai/chat/sessions` | 创建会话 |
| GET | `/ai/chat/sessions/{sessionId}` | 会话消息 |
| POST | `/ai/chat/sessions/{sessionId}/messages` | 发送消息 |
| PUT | `/ai/chat/messages/{messageId}` | 编辑消息 |
| POST | `/ai/topic-agent` | 话题 Agent |

## Admin AI

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/admin/ai/models` | 基础模型列表 |
| POST | `/admin/ai/models` | 新增模型 |
| PUT | `/admin/ai/models/{modelId}` | 编辑模型 |
| PUT | `/admin/ai/models/{modelId}/default` | 设为默认 |
| PUT | `/admin/ai/models/{modelId}/enable` | 启用 |
| PUT | `/admin/ai/models/{modelId}/disable` | 禁用 |
| GET | `/admin/ai/token-config` | Token 配置 |
| PUT | `/admin/ai/token-config` | 保存 Token 配置 |
| GET | `/admin/ai/token-usages` | 用户 Token 使用情况 |
