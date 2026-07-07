package service

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// TagService 标签业务逻辑。
type TagService struct {
	tags    *repository.TagRepository
	postSvc *PostService
}

func NewTagService(tags *repository.TagRepository, postSvc *PostService) *TagService {
	return &TagService{tags: tags, postSvc: postSvc}
}

// List 标签列表（启用中）。
func (s *TagService) List() ([]dto.ContentTag, error) {
	rows, err := s.tags.ListEnabled()
	if err != nil {
		return nil, errs.ErrInternal
	}
	list := make([]dto.ContentTag, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToContentTag(&rows[i]))
	}
	return list, nil
}

// Get 标签详情。
func (s *TagService) Get(id int64) (*dto.ContentTag, error) {
	t, err := s.tags.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if t == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToContentTag(t)
	return &res, nil
}

// Posts 标签下的帖子。
func (s *TagService) Posts(tagID, viewerID int64, sort string, page, size int) ([]dto.Post, int64, error) {
	return s.postSvc.List(ListFilter{TagID: tagID, ViewerID: viewerID, Sort: sort, Page: page, PageSize: size})
}
