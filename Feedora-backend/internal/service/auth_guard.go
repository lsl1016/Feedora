package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/jwtx"
)

// userStatusCacheTTL 用户状态缓存有效期；封禁/解禁时由 admin 侧主动失效。
const userStatusCacheTTL = 10 * time.Minute

// AuthGuard 鉴权链补充校验：登出 token 黑名单、用户级吊销（封禁）、用户状态。
// Redis 未启用时退化为每次请求查库校验用户状态，黑名单与吊销不生效。
type AuthGuard struct {
	users *repository.UserRepository
	cache *cache.Cache
}

func NewAuthGuard(users *repository.UserRepository, cch *cache.Cache) *AuthGuard {
	return &AuthGuard{users: users, cache: cch}
}

// Check 校验 token 与用户状态，返回非 nil error 表示拒绝。
func (g *AuthGuard) Check(claims *jwtx.Claims, token string) error {
	ctx := context.Background()

	// 1. 登出黑名单（按 token）。
	if n, ok := g.cache.GetIntOK(ctx, cache.TokenBLKey(hashToken(token))); ok && n > 0 {
		return errs.ErrUnauth
	}

	// 2. 用户级吊销（封禁）：签发时间早于吊销时间戳的 token 全部失效。
	if n, ok := g.cache.GetIntOK(ctx, cache.UserRevokedBeforeKey(claims.UserID)); ok &&
		claims.IssuedAt != nil && claims.IssuedAt.Unix() < n {
		return errs.ErrUnauth
	}

	// 3. 用户状态：cache-aside，未命中回源 DB。
	status := ""
	if !g.cache.GetJSON(ctx, cache.UserStatusKey(claims.UserID), &status) {
		u, err := g.users.FindByID(claims.UserID)
		if err != nil || u == nil {
			return errs.ErrUnauth
		}
		status = u.Status
		g.cache.SetJSON(ctx, cache.UserStatusKey(claims.UserID), status, userStatusCacheTTL)
	}
	if status == model.UserBanned {
		return errs.ErrUserBanned
	}
	return nil
}

// hashToken 计算 token 的 SHA-256 十六进制摘要，用作黑名单 Key 的一部分。
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
