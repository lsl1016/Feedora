---
title: 认证模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: auth
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/auth_router.go
  - internal/api/auth_api.go
  - internal/dto/auth_dto.go
  - internal/service/auth_service.go
  - internal/model/user.go
  - pkg/jwtx/jwt.go
  - pkg/middleware/auth.go
  - pkg/middleware/context.go
summary: auth 模块提供账号密码注册、登录、当前用户信息获取与登出，基于 bcrypt 密码哈希与 HS256 无状态 JWT 完成身份认证。
---

## 1. 模块概述

auth 模块负责用户身份认证，覆盖注册、账号密码登录、获取当前登录用户信息与退出登录四项能力，由 `AuthAPI` / `AuthService` / `jwtx.Manager` / `middleware.Auth` 组成的竖切片实现。认证采用无状态 JWT（HS256 签名），密码使用 bcrypt 哈希存储。本模块不涉及内容类资源（帖子、评论、圈子等）的管理，也不负责用户资料编辑（属于 user 模块）。

## 2. 接口清单

路由统一注册在 `/api/v1` 前缀下，由 `registerAuth`（`internal/router/auth_router.go`）完成。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/auth/register` | POST | 注册新用户并直接返回登录态 | `AuthAPI.Register` |
| `/auth/login` | POST | 账号密码登录 | `AuthAPI.Login` |
| `/auth/me` | GET | 获取当前登录用户信息（需登录） | `AuthAPI.Me` |
| `/auth/logout` | POST | 退出登录（需登录，服务端无实际动作） | `AuthAPI.Logout` |

## 3. 核心逻辑

### 3.1 密码哈希

- 注册时调用 `utils.HashPassword`（`pkg/utils/password.go`）使用 `bcrypt.GenerateFromPassword` 生成哈希，cost 取 `bcrypt.DefaultCost`（10），存入 `users.password_hash`。
- 登录时调用 `utils.CheckPassword`，以 `bcrypt.CompareHashAndPassword` 比对哈希与明文，返回 `bool`。

### 3.2 JWT 签发与校验

- `jwtx.Manager`（`pkg/jwtx/jwt.go`）持有 `secret []byte` 与 `expireHours int`，由 `jwtx.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)` 在 `internal/app/app.go` 装配；当 `expireHours <= 0` 时回退为 `168`（7 天）。
- 签发 `Generate(userID, role)`：构造 `jwtx.Claims`，以 `jwt.SigningMethodHS256`（HMAC-SHA256）签名，输出字符串 token。
- 解析 `Parse(tokenStr)`：`jwt.ParseWithClaims` 校验签名与过期时间；keyfunc 中显式拒绝非 HMAC 族签名算法（返回 `unexpected signing method`）；校验失败或 claims 断言失败返回错误。

### 3.3 Claims 结构与 token 生命周期

`jwtx.Claims` 自定义字段与标准字段如下：

| 字段 | 类型 | 来源 | 说明 |
|---|---|---|---|
| `userId` | int64 | `u.ID` | 用户 ID |
| `role` | string | `u.Role` | 用户角色（`user` / `admin`） |
| `exp` | NumericDate | 签发时间 + `expireHours` 小时 | 过期时间 |
| `iat` | NumericDate | 签发时刻 | 签发时间 |

未设置 `iss`、`sub`、`aud` 等其余注册声明。token 生命周期为单一访问令牌：注册与登录成功即签发；不设刷新令牌（refresh token）；服务端不存 token、无黑名单，token 在过期前始终有效；登出由前端清除 token，服务端不做失效处理。

### 3.4 鉴权中间件与上下文取当前用户

- `middleware.Auth(jm)`（`pkg/middleware/auth.go`）：从 `Authorization` 请求头解析 `Bearer ` 前缀提取 token，调用 `jm.Parse`；失败返回 401（`errs.ErrUnauth`）并 `Abort`；成功后将 `claims.UserID`、`claims.Role` 以键 `userId`、`role`（`CtxUserID` / `CtxRole`，定义于 `pkg/middleware/context.go`）写入 `gin.Context`。
- `middleware.OptionalAuth(jm)`：在 `router.New` 中全局挂载于 `/api/v1` 分组，尝试解析 token，成功写入上下文，失败放行，用于匿名可访问但需识别登录态的接口。
- 控制器与服务层通过 `middleware.CurrentUserID(c)` / `middleware.CurrentRole(c)` 读取当前用户，未登录时分别返回 `0` 与空串。`AuthAPI.Me` 即以 `middleware.CurrentUserID(c)` 传入 `AuthService.Me`。

### 3.5 注册流程（`AuthService.Register`）

1. 对 `Account`、`Nickname` 做 `TrimSpace`；账号或密码为空返回 `errs.ErrParams`。
2. `ConfirmPassword` 非空时参与校验，与 `Password` 不一致返回 422「两次输入的密码不一致」；为空则跳过该检查。
3. `Nickname` 为空时默认取 `Account`。
4. `CountByAccount` 查重，已存在返回 `errs.ErrAccountExists`。
5. bcrypt 哈希后创建用户：`Role = model.RoleUser`（`user`）、`Status = model.UserNormal`（`normal`）、`Level = 1`。
6. 通过 `event.Producer` 向主题 `event.TopicUser`（`user.events`）发布 `event.UserRegistered` 事件（operatorID 与 targetID 均为新用户 ID，payload 携带 `account`）。
7. 调用 `issue` 签发 token，返回 `dto.LoginResult`（`token` + `user`），即注册成功即视为已登录。

### 3.6 登录流程（`AuthService.Login`）

1. 对账号 `TrimSpace`，账号或密码为空返回 `errs.ErrParams`。
2. `FindByAccount` 查询用户；用户不存在与密码校验失败统一返回 `errs.ErrAccountOrPassword`，不区分具体原因。
3. `Status == model.UserBanned`（`banned`）返回 `errs.ErrUserBanned`。
4. 通过 `issue` 签发 token 并返回登录结果。

### 3.7 登出（`AuthAPI.Logout`）

无状态 JWT 方案下登出接口直接返回空对象 `gin.H{}`，路由仍挂 `middleware.Auth` 要求携带有效 token；服务端不登记、不吊销 token，token 失效完全依赖前端删除与自然过期。

## 4. 数据模型

本模块不新增数据表，全部读写复用 user 模块的 `users` 表（`model.User`，`internal/model/user.go`，表名 `users`）。认证流程涉及的关键字段：

| 字段 | 列 | 说明 |
|---|---|---|
| `ID` | `id` | 主键，int64，签入 claims 的 `userId` |
| `Account` | `account` | 登录账号，`size:64`，唯一索引 `uk_account`，注册查重依据 |
| `PasswordHash` | `password_hash` | bcrypt 哈希，`size:255`，不出现在任何响应 DTO |
| `Role` | `role` | 角色，默认 `user`，签入 claims 的 `role` |
| `Status` | `status` | 账号状态，默认 `normal`，登录时拦截 `banned` |
| `DeletedAt` | `deleted_at` | gorm 软删除标记，带索引 |

`/auth/me` 与登录响应中的用户对象由 `dto.ToUser(u, false, 0)` 组装（`internal/dto/user_dto.go`），包含昵称、头像、角色、状态、积分、等级及各计数字段，不含密码哈希。

## 5. 配置项与限制

配置定义于 `pkg/config/config.go` 的 `JWTConfig`，取值来自 `configs/config.yaml`：

| 配置键 | 类型 | 默认值（config.yaml） | 说明 |
|---|---|---|---|
| `jwt.secret` | string | `feedora-secret-change-me` | HS256 签名密钥 |
| `jwt.expireHours` | int | `168` | token 有效期（小时）；代码侧 `expireHours <= 0` 时回退 `168` |

限制与未实现项：

- 无刷新令牌机制，token 过期后需重新调用 `/auth/login`。
- 无服务端 token 吊销 / 黑名单，登出仅由前端删除 token 完成。
- 注册接口无验证码、无邮箱或手机号校验，账号为任意非空字符串。
- 密码无最小长度等强度校验，仅要求非空。
- 未实现第三方登录（OAuth）与密码找回。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
|---|---|---|---|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
