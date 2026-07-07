package worker

import (
	"context"
	"strconv"

	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/logger"
)

func sid(id int64) string { return strconv.FormatInt(id, 10) }

// reindexAll 启动时全量重建各索引，保证已有数据可被搜索。
func (r *Runner) reindexAll(ctx context.Context) {
	var posts []model.Post
	r.db.Where("status <> ?", model.PostDeleted).Find(&posts)
	for i := range posts {
		r.indexPost(ctx, &posts[i])
	}

	var users []model.User
	r.db.Where("status <> ?", model.UserBanned).Find(&users)
	for i := range users {
		r.indexUser(ctx, &users[i])
	}

	var circles []model.Circle
	r.db.Where("status = ?", model.CircleMemberNormal).Find(&circles)
	for i := range circles {
		r.indexCircle(ctx, &circles[i])
	}

	var topics []model.Topic
	r.db.Where("status = ?", model.StatusEnabled).Find(&topics)
	for i := range topics {
		r.indexTopic(ctx, &topics[i])
	}
	logger.Infof("ES 全量重建完成: posts=%d users=%d circles=%d topics=%d", len(posts), len(users), len(circles), len(topics))
}

// indexPostByID 按 ID 重新索引帖子（不存在或已删除则从索引移除）。
func (r *Runner) indexPostByID(ctx context.Context, id int64) {
	if r.sc == nil {
		return
	}
	p, _ := r.posts.FindByID(id)
	if p == nil || p.Status == model.PostDeleted {
		_ = r.sc.DeleteDoc(ctx, r.sc.PostIndex(), sid(id))
		return
	}
	r.indexPost(ctx, p)
}

func (r *Runner) indexPost(ctx context.Context, p *model.Post) {
	if r.sc == nil {
		return
	}
	authorName := ""
	if a, _ := r.users.FindByID(p.AuthorID); a != nil {
		authorName = a.Nickname
	}
	tagNames := []string{}
	for _, t := range r.tags.FindByPostIDs([]int64{p.ID})[p.ID] {
		tagNames = append(tagNames, t.Name)
	}
	topicNames := []string{}
	for _, t := range r.topics.FindByPostIDs([]int64{p.ID})[p.ID] {
		topicNames = append(topicNames, t.Name)
	}
	var circleID int64
	if p.CircleID != nil {
		circleID = *p.CircleID
	}
	doc := map[string]any{
		"postId": p.ID, "title": p.Title, "content": p.ContentMD, "summary": p.Summary,
		"authorId": p.AuthorID, "authorName": authorName,
		"tagNames": tagNames, "topicNames": topicNames,
		"circleId": circleID, "visibility": p.Visibility, "status": p.Status,
		"likeCount": p.LikeCount, "commentCount": p.CommentCount, "favoriteCount": p.FavoriteCount,
		"hotScore": p.HotScore, "createdAt": p.CreatedAt,
	}
	if err := r.sc.IndexDoc(ctx, r.sc.PostIndex(), sid(p.ID), doc); err != nil {
		logger.Errorf("索引帖子失败 postId=%d: %v", p.ID, err)
	}
}

func (r *Runner) indexUser(ctx context.Context, u *model.User) {
	if r.sc == nil {
		return
	}
	doc := map[string]any{
		"userId": u.ID, "nickname": u.Nickname, "bio": u.Bio, "avatar": u.Avatar,
		"status": u.Status, "followerCount": u.FollowerCount, "postCount": u.PostCount, "createdAt": u.CreatedAt,
	}
	if err := r.sc.IndexDoc(ctx, r.sc.UserIndex(), sid(u.ID), doc); err != nil {
		logger.Errorf("索引用户失败 userId=%d: %v", u.ID, err)
	}
}

func (r *Runner) indexCircle(ctx context.Context, c *model.Circle) {
	if r.sc == nil {
		return
	}
	doc := map[string]any{
		"circleId": c.ID, "name": c.Name, "description": c.Description, "avatar": c.Avatar,
		"category": c.Category, "status": c.Status, "memberCount": c.MemberCount, "postCount": c.PostCount,
		"featuredCount": c.FeaturedCount, "isRecommended": c.IsRecommended, "createdAt": c.CreatedAt,
	}
	if err := r.sc.IndexDoc(ctx, r.sc.CircleIndex(), sid(c.ID), doc); err != nil {
		logger.Errorf("索引圈子失败 circleId=%d: %v", c.ID, err)
	}
}

func (r *Runner) indexTopic(ctx context.Context, t *model.Topic) {
	if r.sc == nil {
		return
	}
	doc := map[string]any{
		"topicId": t.ID, "name": t.Name, "description": t.Description, "coverUrl": t.CoverURL,
		"isOfficial": t.IsOfficial, "isRecommended": t.IsRecommended,
		"participantCount": t.ParticipantCount, "postCount": t.PostCount, "status": t.Status, "createdAt": t.CreatedAt,
	}
	if err := r.sc.IndexDoc(ctx, r.sc.TopicIndex(), sid(t.ID), doc); err != nil {
		logger.Errorf("索引话题失败 topicId=%d: %v", t.ID, err)
	}
}
