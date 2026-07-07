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
