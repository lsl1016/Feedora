// Package cache 封装 Redis 业务缓存、计数、榜单 ZSet、通知未读数与幂等。
package cache

import "fmt"

// Redis Key 统一生成，格式：业务域:对象:标识[:附加维度]。
func PostDetailKey(postID int64) string   { return fmt.Sprintf("post:detail:%d", postID) }
func UserProfileKey(userID int64) string  { return fmt.Sprintf("user:profile:%d", userID) }
func CircleDetailKey(cid int64) string    { return fmt.Sprintf("circle:detail:%d", cid) }
func NotifyUnreadKey(userID int64) string { return fmt.Sprintf("notify:unread:%d", userID) }

// RankKey 榜单 ZSet：rank:{rankType}:{timeRange}，如 rank:post:today。
func RankKey(rankType, timeRange string) string {
	return fmt.Sprintf("rank:%s:%s", rankType, timeRange)
}

// UserRankKey 用户排行榜：rank:user:{type}:{range}，如 rank:user:creator:all。
func UserRankKey(rankType, timeRange string) string {
	return fmt.Sprintf("rank:user:%s:%s", rankType, timeRange)
}

// HotKeywordsKey 热门搜索词 ZSet。
const HotKeywordsKey = "search:hot_keywords"

// 用户状态缓存（鉴权链消费，10 分钟 TTL，封禁/解禁时主动失效）。
func UserStatusKey(userID int64) string { return fmt.Sprintf("user:status:%d", userID) }

// TokenBLKey 登出黑名单：值为 1，TTL 为 token 剩余有效期。
func TokenBLKey(tokenHash string) string { return fmt.Sprintf("auth:bl:%s", tokenHash) }

// UserRevokedBeforeKey 用户级吊销时间戳（Unix 秒）：签发时间早于该值的 token 全部失效，用于封禁等场景。
func UserRevokedBeforeKey(userID int64) string { return fmt.Sprintf("auth:revoked_before:%d", userID) }

// IdemPostKey / IdemCommentKey 写操作防重复提交（同作者同内容时间窗）。
func IdemPostKey(userID int64, hash string) string {
	return fmt.Sprintf("idem:post:%d:%s", userID, hash)
}
func IdemCommentKey(userID int64, hash string) string {
	return fmt.Sprintf("idem:comment:%d:%s", userID, hash)
}
