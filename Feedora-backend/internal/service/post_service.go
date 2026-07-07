package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// postDetailTTL 帖子详情缓存时长。
const postDetailTTL = 10 * time.Minute

// PostService 帖子业务逻辑，同时承担帖子聚合装配，供其它服务复用。
type PostService struct {
	posts    *repository.PostRepository
	users    *repository.UserRepository
	tags     *repository.TagRepository
	topics   *repository.TopicRepository
	circles  *repository.CircleRepository
	inters   *repository.InteractionRepository
	producer event.Producer
	cache    *cache.Cache
}

func NewPostService(
	posts *repository.PostRepository,
	users *repository.UserRepository,
	tags *repository.TagRepository,
	topics *repository.TopicRepository,
	circles *repository.CircleRepository,
	inters *repository.InteractionRepository,
	producer event.Producer,
	cch *cache.Cache,
) *PostService {
	return &PostService{posts: posts, users: users, tags: tags, topics: topics, circles: circles, inters: inters, producer: producer, cache: cch}
}

// InvalidateDetail 失效帖子详情缓存，供互动 / 评论等模块在计数变化后调用。
func (s *PostService) InvalidateDetail(postID int64) {
	s.cache.Del(context.Background(), cache.PostDetailKey(postID))
}

// ListFilter 帖子列表过滤条件（业务层）。
type ListFilter struct {
	FeedType      string
	Sort          string
	Status        string
	Keyword       string
	TagID         int64
	CircleID      int64
	TopicID       int64
	AuthorID      int64
	IncludeHidden bool
	ViewerID      int64
	Page          int
	PageSize      int
}

// List 分页查询帖子并装配为 DTO。
func (s *PostService) List(f ListFilter) ([]dto.Post, int64, error) {
	page, size := normPage(f.Page, f.PageSize)
	rows, total, err := s.posts.List(repository.PostFilter{
		FeedType:      f.FeedType,
		Sort:          f.Sort,
		Status:        f.Status,
		Keyword:       f.Keyword,
		TagID:         f.TagID,
		CircleID:      f.CircleID,
		TopicID:       f.TopicID,
		AuthorID:      f.AuthorID,
		IncludeHidden: f.IncludeHidden,
		ViewerID:      f.ViewerID,
		Offset:        (page - 1) * size,
		Limit:         size,
	})
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list, err := s.Assemble(rows, f.ViewerID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Get 查询帖子详情。采用 Cache Aside：详情主体（不含 per-viewer 状态）缓存到 Redis，
// 每次按当前用户叠加 liked/favorited，并递增浏览量。
func (s *PostService) Get(id, viewerID int64) (*dto.Post, error) {
	ctx := context.Background()
	key := cache.PostDetailKey(id)

	var detail dto.Post
	if s.cache.GetJSON(ctx, key, &detail) {
		s.posts.IncViewCount(id)
		detail.ViewCount++
		s.overlayViewer(&detail, viewerID)
		return &detail, nil
	}

	p, err := s.posts.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if p == nil || p.Status == model.PostDeleted {
		return nil, errs.ErrPostNotFound
	}
	if p.Status != model.PostPublished && p.AuthorID != viewerID {
		return nil, errs.ErrPostInvisible
	}
	s.posts.IncViewCount(id)
	p.ViewCount++
	// 以 viewer=0 装配可共享的详情主体并缓存。
	list, err := s.Assemble([]model.Post{*p}, 0)
	if err != nil {
		return nil, err
	}
	detail = list[0]
	// 仅对已发布的公开帖子缓存，避免缓存不可见内容。
	if p.Status == model.PostPublished {
		s.cache.SetJSON(ctx, key, detail, postDetailTTL)
	}
	s.overlayViewer(&detail, viewerID)
	return &detail, nil
}

// overlayViewer 叠加当前用户对该帖子的点赞 / 收藏状态。
func (s *PostService) overlayViewer(p *dto.Post, viewerID int64) {
	if viewerID <= 0 {
		return
	}
	p.Liked = s.inters.HasLiked(p.PostID, viewerID)
	p.Favorited = s.inters.HasFavorited(p.PostID, viewerID)
}

// AssembleByIDs 按给定顺序装配帖子（用于点赞 / 收藏列表）。
func (s *PostService) AssembleByIDs(ids []int64, viewerID int64) ([]dto.Post, error) {
	rows, err := s.posts.FindByIDs(ids, true)
	if err != nil {
		return nil, errs.ErrInternal
	}
	byID := map[int64]model.Post{}
	for i := range rows {
		byID[rows[i].ID] = rows[i]
	}
	ordered := make([]model.Post, 0, len(ids))
	for _, id := range ids {
		if p, ok := byID[id]; ok {
			ordered = append(ordered, p)
		}
	}
	return s.Assemble(ordered, viewerID)
}

// Assemble 将帖子模型批量装配为 DTO，一次性加载作者、图片、标签、话题、圈子及互动状态。
func (s *PostService) Assemble(rows []model.Post, viewerID int64) ([]dto.Post, error) {
	result := make([]dto.Post, 0, len(rows))
	if len(rows) == 0 {
		return result, nil
	}
	postIDs := make([]int64, 0, len(rows))
	authorIDs := make([]int64, 0, len(rows))
	circleIDs := make([]int64, 0)
	for i := range rows {
		postIDs = append(postIDs, rows[i].ID)
		authorIDs = append(authorIDs, rows[i].AuthorID)
		if rows[i].CircleID != nil {
			circleIDs = append(circleIDs, *rows[i].CircleID)
		}
	}

	authors, err := s.users.FindByIDs(authorIDs)
	if err != nil {
		return nil, errs.ErrInternal
	}
	circles := s.circles.FindByIDs(circleIDs)
	imagesMap := s.posts.ImagesByPostIDs(postIDs)
	tagsMap := s.tags.FindByPostIDs(postIDs)
	topicsMap := s.topics.FindByPostIDs(postIDs)
	likedSet := s.inters.LikedSet(viewerID, postIDs)
	favSet := s.inters.FavoritedSet(viewerID, postIDs)

	for i := range rows {
		p := &rows[i]
		rel := dto.PostRelations{
			Author:    authors[p.AuthorID],
			Images:    imagesMap[p.ID],
			Tags:      toTagSummaries(tagsMap[p.ID]),
			Topics:    toTopicSummaries(topicsMap[p.ID]),
			Liked:     likedSet[p.ID],
			Favorited: favSet[p.ID],
		}
		if p.CircleID != nil {
			rel.Circle = circles[*p.CircleID]
		}
		result = append(result, dto.ToPost(p, rel))
	}
	return result, nil
}

// Create 发布帖子。
func (s *PostService) Create(authorID int64, in dto.CreatePostRequest) (*dto.Post, error) {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" || strings.TrimSpace(in.Content) == "" {
		return nil, errs.ErrParams
	}
	if in.Visibility == "" {
		in.Visibility = model.VisibilityPublic
	}
	status := model.PostPublished
	switch in.PublishMode {
	case "draft":
		status = model.PostDraft
	case "schedule":
		status = model.PostScheduled
	}
	if in.CircleID != nil {
		if err := s.checkCirclePostPermission(*in.CircleID, authorID); err != nil {
			return nil, err
		}
	}
	now := time.Now()
	p := &model.Post{
		AuthorID:   authorID,
		Title:      in.Title,
		ContentMD:  in.Content,
		Summary:    summarize(in.Content),
		PostType:   model.PostTypeOriginal,
		CircleID:   in.CircleID,
		Visibility: in.Visibility,
		Status:     status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if len(in.Images) > 0 {
		p.CoverURL = in.Images[0]
	}
	if status == model.PostPublished {
		p.PublishedAt = &now
	}
	if err := s.posts.CreateWithRelations(p, in.Images, dedup(in.TagIDs), dedup(in.TopicIDs)); err != nil {
		return nil, errs.ErrInternal
	}
	s.producer.Publish(event.TopicPost, event.PostCreated, p.ID, authorID, map[string]any{"title": p.Title})
	return s.Get(p.ID, authorID)
}

// Update 编辑帖子，仅作者本人可操作。
func (s *PostService) Update(id, userID int64, in dto.UpdatePostRequest) (*dto.Post, error) {
	p, err := s.posts.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if p == nil {
		return nil, errs.ErrPostNotFound
	}
	if p.AuthorID != userID {
		return nil, errs.ErrPostNoPermission
	}
	updates := map[string]any{"updated_at": time.Now()}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Content != nil {
		updates["content_md"] = *in.Content
		updates["summary"] = summarize(*in.Content)
	}
	if in.Visibility != nil {
		updates["visibility"] = *in.Visibility
	}
	if err := s.posts.Update(id, updates); err != nil {
		return nil, errs.ErrInternal
	}
	if in.Images != nil {
		s.posts.ReplaceImages(id, in.Images)
	}
	s.InvalidateDetail(id)
	s.producer.Publish(event.TopicPost, event.PostUpdated, id, userID, nil)
	return s.Get(id, userID)
}

// SetHidden 隐藏 / 取消隐藏帖子，仅作者本人可操作。
func (s *PostService) SetHidden(id, userID int64, hidden bool) error {
	p, err := s.posts.FindByID(id)
	if err != nil {
		return errs.ErrInternal
	}
	if p == nil {
		return errs.ErrPostNotFound
	}
	if p.AuthorID != userID {
		return errs.ErrPostNoPermission
	}
	status := model.PostPublished
	if hidden {
		status = model.PostHidden
	}
	s.posts.Update(id, map[string]any{"status": status, "updated_at": time.Now()})
	s.InvalidateDetail(id)
	if hidden {
		s.producer.Publish(event.TopicPost, event.PostHidden, id, userID, nil)
	} else {
		s.producer.Publish(event.TopicPost, "PostUnhidden", id, userID, nil)
	}
	return nil
}

// Delete 软删除帖子，作者本人或管理员可操作。
func (s *PostService) Delete(id, userID int64, isAdmin bool) error {
	p, err := s.posts.FindByID(id)
	if err != nil {
		return errs.ErrInternal
	}
	if p == nil {
		return errs.ErrPostNotFound
	}
	if p.AuthorID != userID && !isAdmin {
		return errs.ErrPostNoPermission
	}
	s.posts.SoftDelete(id)
	s.users.IncColumn(p.AuthorID, "post_count", -1)
	s.InvalidateDetail(id)
	s.producer.Publish(event.TopicPost, event.PostDeleted, id, userID, nil)
	return nil
}

// Share 分享计数 +1。
func (s *PostService) Share(id int64) error {
	if !s.posts.Exists(id) {
		return errs.ErrPostNotFound
	}
	s.posts.IncColumn(id, "share_count", 1)
	return nil
}

// Repost 转发帖子。
func (s *PostService) Repost(id, userID int64, comment string) (*dto.Post, error) {
	src, err := s.posts.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if src == nil {
		return nil, errs.ErrPostNotFound
	}
	now := time.Now()
	p := &model.Post{
		AuthorID:      userID,
		Title:         "转发：" + src.Title,
		ContentMD:     comment,
		Summary:       summarize(comment),
		PostType:      model.PostTypeRepost,
		SourcePostID:  &src.ID,
		RepostComment: comment,
		Visibility:    model.VisibilityPublic,
		Status:        model.PostPublished,
		PublishedAt:   &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.posts.Create(p); err != nil {
		return nil, errs.ErrInternal
	}
	s.posts.IncColumn(id, "repost_count", 1)
	s.users.IncColumn(userID, "post_count", 1)
	return s.Get(p.ID, userID)
}

// checkCirclePostPermission 校验用户是否有在圈子发帖的权限。
func (s *PostService) checkCirclePostPermission(circleID, userID int64) error {
	circle, err := s.circles.FindByID(circleID)
	if err != nil {
		return errs.ErrInternal
	}
	if circle == nil {
		return errs.ErrCircleNotFound
	}
	member, err := s.circles.FindMember(circleID, userID)
	if err != nil {
		return errs.ErrInternal
	}
	if member == nil || member.Status == model.CircleMemberRemoved {
		return errs.ErrCircleNotJoined
	}
	if member.Status == model.CircleMemberMuted {
		return errs.ErrCircleMuted
	}
	if circle.PostPermission == "admin_only" && member.Role != model.CircleRoleOwner && member.Role != model.CircleRoleModerator {
		return errs.ErrPostNoPermission
	}
	return nil
}

// summarize 从正文截取摘要，最多 120 个字符。
func summarize(content string) string {
	const max = 120
	if utf8.RuneCountInString(content) <= max {
		return content
	}
	runes := []rune(content)
	return string(runes[:max]) + "..."
}

// dedup 去重并过滤非正整数 ID，保持顺序。
func dedup(ids []int64) []int64 {
	seen := map[int64]struct{}{}
	res := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		res = append(res, id)
	}
	return res
}

func toTagSummaries(tags []model.Tag) []dto.TagSummary {
	res := make([]dto.TagSummary, 0, len(tags))
	for i := range tags {
		res = append(res, dto.ToTagSummary(&tags[i]))
	}
	return res
}

func toTopicSummaries(topics []model.Topic) []dto.TopicSummary {
	res := make([]dto.TopicSummary, 0, len(topics))
	for i := range topics {
		res = append(res, dto.ToTopicSummary(&topics[i]))
	}
	return res
}
