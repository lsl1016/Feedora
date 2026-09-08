package service

import (
	"context"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// InteractionService 点赞 / 收藏业务逻辑。
type InteractionService struct {
	posts    *repository.PostRepository
	users    *repository.UserRepository
	inters   *repository.InteractionRepository
	producer event.Producer
	cache    *cache.Cache
}

func NewInteractionService(
	posts *repository.PostRepository,
	users *repository.UserRepository,
	inters *repository.InteractionRepository,
	producer event.Producer,
	cch *cache.Cache,
) *InteractionService {
	return &InteractionService{posts: posts, users: users, inters: inters, producer: producer, cache: cch}
}

// invalidate 失效帖子详情缓存（计数变化后）。
func (s *InteractionService) invalidate(postID int64) {
	s.cache.Del(context.Background(), cache.PostDetailKey(postID))
}

// SetLike 点赞 / 取消点赞帖子，唯一索引保证幂等。
func (s *InteractionService) SetLike(postID, userID int64, like bool) (*dto.InteractionResult, error) {
	p, err := s.posts.FindByID(postID)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if p == nil {
		return nil, errs.ErrPostNotFound
	}
	if like {
		if s.inters.AddPostLike(postID, userID) {
			s.posts.IncColumn(postID, "like_count", 1)
			s.users.IncColumn(p.AuthorID, "like_count", 1)
			s.producer.Publish(event.TopicInteraction, event.PostLiked, postID, userID, nil)
		}
	} else {
		if s.inters.RemovePostLike(postID, userID) {
			s.posts.IncColumn(postID, "like_count", -1)
			s.users.IncColumn(p.AuthorID, "like_count", -1)
			s.producer.Publish(event.TopicInteraction, event.PostUnliked, postID, userID, nil)
		}
	}
	s.invalidate(postID)
	return s.status(postID, userID)
}

// SetFavorite 收藏 / 取消收藏帖子。
func (s *InteractionService) SetFavorite(postID, userID int64, fav bool) (*dto.InteractionResult, error) {
	if !s.posts.Exists(postID) {
		return nil, errs.ErrPostNotFound
	}
	if fav {
		if s.inters.AddPostFavorite(postID, userID) {
			s.posts.IncColumn(postID, "favorite_count", 1)
			s.producer.Publish(event.TopicInteraction, event.PostFavorited, postID, userID, nil)
		}
	} else {
		if s.inters.RemovePostFavorite(postID, userID) {
			s.posts.IncColumn(postID, "favorite_count", -1)
			s.producer.Publish(event.TopicInteraction, event.PostUnfavorited, postID, userID, nil)
		}
	}
	s.invalidate(postID)
	return s.status(postID, userID)
}

// status 查询帖子最新计数与当前用户互动状态。
func (s *InteractionService) status(postID, userID int64) (*dto.InteractionResult, error) {
	p, err := s.posts.FindByID(postID)
	if err != nil || p == nil {
		return nil, errs.ErrPostNotFound
	}
	return &dto.InteractionResult{
		Liked:         s.inters.HasLiked(postID, userID),
		Favorited:     s.inters.HasFavorited(postID, userID),
		LikeCount:     p.LikeCount,
		FavoriteCount: p.FavoriteCount,
	}, nil
}
