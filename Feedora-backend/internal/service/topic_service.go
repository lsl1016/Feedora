package service

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// TopicService 话题业务逻辑。
type TopicService struct {
	topics  *repository.TopicRepository
	postSvc *PostService
}

func NewTopicService(topics *repository.TopicRepository, postSvc *PostService) *TopicService {
	return &TopicService{topics: topics, postSvc: postSvc}
}

// List 话题广场，tab 支持 all/official/hot/latest。
func (s *TopicService) List(tab string, page, size int) ([]dto.Topic, int64, error) {
	page, size = normPage(page, size)
	rows, total, err := s.topics.List(tab, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.Topic, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToTopic(&rows[i]))
	}
	return list, total, nil
}

// Get 话题详情。
func (s *TopicService) Get(id int64) (*dto.Topic, error) {
	t, err := s.topics.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if t == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToTopic(t)
	return &res, nil
}

// Posts 话题下的帖子。
func (s *TopicService) Posts(topicID, viewerID int64, sort string, page, size int) ([]dto.Post, int64, error) {
	return s.postSvc.List(ListFilter{TopicID: topicID, ViewerID: viewerID, Sort: sort, Page: page, PageSize: size})
}
