package worker

import (
	"context"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/logger"
)

// handle 将一条事件分发给各消费者（每个消费者独立幂等）。
func (r *Runner) handle(ctx context.Context, m *event.Message) {
	start := time.Now()
	r.handleSearch(ctx, m)
	r.handleNotification(ctx, m)
	r.handleGrowth(ctx, m)
	r.handleRank(ctx, m)
	r.handleStat(ctx, m)
	logger.Infof("consume event success, eventType:%s, eventId:%s, traceId:%s, aggregateId:%d, durationMs:%d",
		m.EventType, m.EventID, m.TraceID, m.AggregateID, time.Since(start).Milliseconds())
}

// handleSearch 同步 ES 索引。
func (r *Runner) handleSearch(ctx context.Context, m *event.Message) {
	if r.sc == nil {
		return
	}
	const w = "search"
	switch m.EventType {
	case event.PostCreated, event.PostUpdated, event.PostHidden, event.PostDeleted, event.PostUnhidden:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.indexPostByID(ctx, m.AggregateID)
	case event.UserRegistered:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		if u, _ := r.users.FindByID(m.AggregateID); u != nil {
			r.indexUser(ctx, u)
		}
	case event.CircleCreated:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		if c, _ := r.circles.FindByID(m.AggregateID); c != nil {
			r.indexCircle(ctx, c)
		}
	case event.TopicCreated, event.TopicUpdated:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		if t, _ := r.topics.FindByID(m.AggregateID); t != nil {
			r.indexTopic(ctx, t)
		}
	}
}

// handleNotification 生成站内通知并累加未读数。
func (r *Runner) handleNotification(ctx context.Context, m *event.Message) {
	const w = "notification"
	switch m.EventType {
	case event.PostLiked, event.PostFavorited:
		p, _ := r.posts.FindByID(m.AggregateID)
		if p == nil || p.AuthorID == m.UserID {
			return
		}
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		title, content := "收到新的点赞", "有人点赞了你的帖子《"+p.Title+"》"
		if m.EventType == event.PostFavorited {
			title, content = "收到新的收藏", "有人收藏了你的帖子《"+p.Title+"》"
		}
		r.notify(ctx, p.AuthorID, m.UserID, "interaction", title, content, "post", p.ID)
	case event.CommentCreated:
		cm, _ := r.comments.FindByID(m.AggregateID)
		if cm == nil {
			return
		}
		p, _ := r.posts.FindByID(cm.PostID)
		if p == nil || p.AuthorID == cm.UserID {
			return
		}
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.notify(ctx, p.AuthorID, cm.UserID, "interaction", "收到新的评论", "有人评论了你的帖子《"+p.Title+"》", "post", p.ID)
	case event.CircleJoined:
		c, _ := r.circles.FindByID(m.AggregateID)
		if c == nil || c.OwnerID == m.UserID {
			return
		}
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.notify(ctx, c.OwnerID, m.UserID, "circle", "圈子有新成员", "有新成员加入了你的圈子「"+c.Name+"」", "circle", c.ID)
	case event.UserFollowed:
		u, _ := r.users.FindByID(m.UserID)
		if u == nil || m.AggregateID == m.UserID {
			return
		}
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.notify(ctx, m.AggregateID, m.UserID, "follow", "收到新的关注", u.Nickname+" 关注了你", "user", m.UserID)
	}
}

// notify 写入通知并累加未读数。
func (r *Runner) notify(ctx context.Context, userID, actorID int64, category, title, content, targetType string, targetID int64) {
	now := time.Now()
	targetURL := ""
	switch targetType {
	case "post":
		targetURL = "/posts/" + sid(targetID)
	case "circle":
		targetURL = "/circles/" + sid(targetID)
	case "user":
		targetURL = "/users/" + sid(targetID)
	}
	if err := r.notifs.Insert(&model.Notification{
		UserID: userID, ActorID: actorID, Type: category, Title: title, Content: content,
		TargetType: targetType, TargetID: targetID, TargetURL: targetURL,
		ReadStatus: "unread", CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		logger.Errorf("生成通知失败: %v", err)
		return
	}
	// 未读数走 cache-aside：落库后失效缓存，下次查询回源 DB COUNT。
	r.cch.Del(ctx, cache.NotifyUnreadKey(userID))
}

// handleGrowth 异步发放积分。
func (r *Runner) handleGrowth(ctx context.Context, m *event.Message) {
	const w = "growth"
	switch m.EventType {
	case event.PostCreated:
		p, _ := r.posts.FindByID(m.AggregateID)
		if p == nil || !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(p.AuthorID, "create_post", 10, "post", p.ID, "发帖奖励")
	case event.CommentCreated:
		cm, _ := r.comments.FindByID(m.AggregateID)
		if cm == nil || !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(cm.UserID, "create_comment", 3, "comment", cm.ID, "评论奖励")
	case event.PostLiked:
		p, _ := r.posts.FindByID(m.AggregateID)
		if p == nil || !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(p.AuthorID, "post_liked", 2, "post", p.ID, "帖子被点赞")
	case event.PostFavorited:
		p, _ := r.posts.FindByID(m.AggregateID)
		if p == nil || !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(p.AuthorID, "post_favorited", 5, "post", p.ID, "帖子被收藏")
	case event.CircleCreated:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(m.UserID, "create_circle", 20, "circle", m.AggregateID, "创建圈子")
	case event.CircleJoined:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.growth.AddPointLog(m.UserID, "join_circle", 1, "circle", m.AggregateID, "加入圈子")
	}
}

// handleRank 更新热门榜单 ZSet。
func (r *Runner) handleRank(ctx context.Context, m *event.Message) {
	if !r.cch.Enabled() {
		return
	}
	const w = "rank"
	var postID int64
	var delta float64
	switch m.EventType {
	case event.PostCreated:
		postID, delta = m.AggregateID, 1
	case event.PostLiked:
		postID, delta = m.AggregateID, 3
	case event.PostUnliked:
		postID, delta = m.AggregateID, -3
	case event.PostFavorited:
		postID, delta = m.AggregateID, 4
	case event.PostUnfavorited:
		postID, delta = m.AggregateID, -4
	case event.CommentCreated:
		if cm, _ := r.comments.FindByID(m.AggregateID); cm != nil {
			postID, delta = cm.PostID, 5
		}
	case event.CircleJoined:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		member := sid(m.AggregateID)
		for _, tr := range []string{"today", "week", "all"} {
			r.cch.ZIncr(ctx, cache.RankKey("circle", tr), member, 2)
		}
		return
	default:
		return
	}
	if postID == 0 {
		return
	}
	if !r.idem.Claim(m.EventID, w) {
		return
	}
	member := sid(postID)
	for _, tr := range []string{"today", "week", "all"} {
		r.cch.ZIncr(ctx, cache.RankKey("post", tr), member, delta)
	}
}

// handleStat 维护派生统计字段（话题参与人数等）。
func (r *Runner) handleStat(ctx context.Context, m *event.Message) {
	const w = "stat"
	switch m.EventType {
	case event.PostCreated, event.PostDeleted:
		if !r.idem.Claim(m.EventID, w) {
			return
		}
		r.topics.RecomputeParticipantCountsByPost(m.AggregateID)
	}
}
