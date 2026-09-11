package service

import (
	"context"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// revokeTTL 用户级吊销时间戳的缓存时长，与 JWT 默认有效期（168h）对齐。
const revokeTTL = 7 * 24 * time.Hour

// AdminService 后台管理业务逻辑。
type AdminService struct {
	admin    *repository.AdminRepository
	users    *repository.UserRepository
	tags     *repository.TagRepository
	topics   *repository.TopicRepository
	circles  *repository.CircleRepository
	comments *repository.CommentRepository
	producer event.Producer
	cache    *cache.Cache
}

func NewAdminService(
	admin *repository.AdminRepository,
	users *repository.UserRepository,
	tags *repository.TagRepository,
	topics *repository.TopicRepository,
	circles *repository.CircleRepository,
	comments *repository.CommentRepository,
	producer event.Producer,
	cch *cache.Cache,
) *AdminService {
	return &AdminService{admin: admin, users: users, tags: tags, topics: topics, circles: circles, comments: comments, producer: producer, cache: cch}
}

// Users 用户列表。
func (s *AdminService) Users(keyword string, page, size int) ([]dto.User, int64, error) {
	page, size = normPage(page, size)
	rows, total, err := s.users.List(keyword, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.User, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToUser(&rows[i], false, 0))
	}
	return list, total, nil
}

// Posts 帖子列表（可按状态过滤）。
func (s *AdminService) Posts(status string, page, size int) ([]model.Post, int64) {
	page, size = normPage(page, size)
	return s.admin.ListPosts(status, offset(page, size), size)
}

// Comments 评论列表。
func (s *AdminService) Comments(page, size int) ([]model.Comment, int64) {
	page, size = normPage(page, size)
	return s.comments.ListAll(offset(page, size), size), s.comments.Count()
}

// Tags 标签列表。
func (s *AdminService) Tags() ([]dto.ContentTag, error) {
	rows, err := s.tags.ListAll()
	if err != nil {
		return nil, errs.ErrInternal
	}
	list := make([]dto.ContentTag, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToContentTag(&rows[i]))
	}
	return list, nil
}

// CreateTag 创建标签。
func (s *AdminService) CreateTag(in dto.CreateTagRequest) (*dto.ContentTag, error) {
	if in.Name == "" {
		return nil, errs.ErrParams
	}
	now := time.Now()
	t := &model.Tag{Name: in.Name, Description: in.Description, Status: model.StatusEnabled, CreatedAt: now, UpdatedAt: now}
	if err := s.tags.Create(t); err != nil {
		return nil, errs.New(409, "标签已存在")
	}
	res := dto.ToContentTag(t)
	return &res, nil
}

// UpdateTag 更新标签。
func (s *AdminService) UpdateTag(id int64, in dto.UpdateTagRequest) (*dto.ContentTag, error) {
	updates := map[string]any{"updated_at": time.Now()}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}
	if err := s.tags.Update(id, updates); err != nil {
		return nil, errs.ErrInternal
	}
	t, _ := s.tags.FindByID(id)
	if t == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToContentTag(t)
	return &res, nil
}

// Topics 话题列表。
func (s *AdminService) Topics(page, size int) ([]dto.Topic, int64, error) {
	page, size = normPage(page, size)
	rows, total, err := s.topics.ListAll(offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.Topic, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToTopic(&rows[i]))
	}
	return list, total, nil
}

// CreateTopic 创建话题。
func (s *AdminService) CreateTopic(in dto.CreateTopicRequest) (*dto.Topic, error) {
	if in.Name == "" {
		return nil, errs.ErrParams
	}
	now := time.Now()
	t := &model.Topic{
		Name: in.Name, Description: in.Description, CoverURL: in.CoverImage,
		IsOfficial: in.IsOfficial, Status: model.StatusEnabled, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.topics.Create(t); err != nil {
		return nil, errs.New(409, "话题已存在")
	}
	s.producer.Publish(event.TopicTopic, event.TopicCreated, t.ID, 0, nil)
	res := dto.ToTopic(t)
	return &res, nil
}

// UpdateTopic 更新话题（仅允许白名单字段）。
func (s *AdminService) UpdateTopic(id int64, in dto.UpdateTopicRequest) (*dto.Topic, error) {
	updates := map[string]any{"updated_at": time.Now()}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.IsOfficial != nil {
		updates["is_official"] = *in.IsOfficial
	}
	if in.IsRecommended != nil {
		updates["is_recommended"] = *in.IsRecommended
	}
	if in.CoverImage != nil {
		updates["cover_url"] = *in.CoverImage
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}
	if err := s.topics.Update(id, updates); err != nil {
		return nil, errs.ErrInternal
	}
	s.producer.Publish(event.TopicTopic, event.TopicUpdated, id, 0, nil)
	t, _ := s.topics.FindByID(id)
	if t == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToTopic(t)
	return &res, nil
}

// SetUserStatus 封禁 / 解禁用户：更新状态、吊销或恢复 token、失效状态缓存，并记录操作日志。
func (s *AdminService) SetUserStatus(adminID, userID int64, status string) error {
	switch status {
	case model.UserNormal, model.UserBanned:
	default:
		return errs.ErrParams
	}
	u, err := s.users.FindByID(userID)
	if err != nil || u == nil {
		return errs.ErrNotFound
	}
	if err := s.admin.UpdateUserStatus(userID, status); err != nil {
		return errs.ErrInternal
	}
	// 失效状态缓存；封禁时吊销该用户全部已签发 token，解禁时恢复。
	ctx := context.Background()
	s.cache.Del(ctx, cache.UserStatusKey(userID))
	if status == model.UserBanned {
		s.cache.SetInt(ctx, cache.UserRevokedBeforeKey(userID), time.Now().Unix(), revokeTTL)
	} else {
		s.cache.Del(ctx, cache.UserRevokedBeforeKey(userID))
	}
	s.log(adminID, "update_user_status", "user", userID, "用户 "+u.Nickname+" 状态变更为 "+status)
	return nil
}

// SetPostStatus 帖子上下架（published / hidden / takedown），变更后发事件供 ES 索引同步。
func (s *AdminService) SetPostStatus(adminID, postID int64, status string) error {
	switch status {
	case model.PostPublished, model.PostHidden, model.PostTakedown:
	default:
		return errs.ErrParams
	}
	p := s.admin.GetPostByID(postID)
	if p == nil {
		return errs.ErrPostNotFound
	}
	if err := s.admin.UpdatePostStatus(postID, status); err != nil {
		return errs.ErrInternal
	}
	s.log(adminID, "update_post_status", "post", postID, "帖子《"+p.Title+"》状态变更为 "+status)
	switch status {
	case model.PostPublished:
		s.producer.Publish(event.TopicPost, event.PostUnhidden, postID, adminID, nil)
	case model.PostHidden:
		s.producer.Publish(event.TopicPost, event.PostHidden, postID, adminID, nil)
	default:
		s.producer.Publish(event.TopicPost, event.PostUpdated, postID, adminID, nil)
	}
	return nil
}

// log 记录后台操作日志。
func (s *AdminService) log(adminID int64, action, targetType string, targetID int64, detail string) {
	name := ""
	if u, _ := s.users.FindByID(adminID); u != nil {
		name = u.Nickname
	}
	s.admin.AddLog(&model.OperationLog{
		AdminID: adminID, AdminName: name, Action: action,
		TargetType: targetType, TargetID: targetID, Detail: detail,
	})
}

// Circles 圈子列表。
func (s *AdminService) Circles(page, size int) ([]model.Circle, int64) {
	page, size = normPage(page, size)
	return s.circles.ListAll(offset(page, size), size)
}

// Stats 后台概览统计。
func (s *AdminService) Stats() dto.AdminStats {
	u, p, c, ci := s.admin.Stats()
	return dto.AdminStats{UserCount: u, PostCount: p, CommentCount: c, CircleCount: ci}
}

// Logs 操作日志列表。
func (s *AdminService) Logs(page, size int) ([]model.OperationLog, int64) {
	page, size = normPage(page, size)
	return s.admin.ListLogs(offset(page, size), size)
}
