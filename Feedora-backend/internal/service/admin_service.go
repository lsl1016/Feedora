package service

import (
	"time"

	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// AdminService 后台管理业务逻辑。
type AdminService struct {
	admin    *repository.AdminRepository
	users    *repository.UserRepository
	tags     *repository.TagRepository
	topics   *repository.TopicRepository
	circles  *repository.CircleRepository
	comments *repository.CommentRepository
}

func NewAdminService(
	admin *repository.AdminRepository,
	users *repository.UserRepository,
	tags *repository.TagRepository,
	topics *repository.TopicRepository,
	circles *repository.CircleRepository,
	comments *repository.CommentRepository,
) *AdminService {
	return &AdminService{admin: admin, users: users, tags: tags, topics: topics, circles: circles, comments: comments}
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
	t, _ := s.topics.FindByID(id)
	if t == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToTopic(t)
	return &res, nil
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
