package service

import (
	"context"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

const userProfileTTL = 30 * time.Minute

// UserService 用户业务逻辑。
type UserService struct {
	users      *repository.UserRepository
	inters     *repository.InteractionRepository
	postSvc    *PostService
	commentSvc *CommentService
	cache      *cache.Cache
}

func NewUserService(
	users *repository.UserRepository,
	inters *repository.InteractionRepository,
	postSvc *PostService,
	commentSvc *CommentService,
	cch *cache.Cache,
) *UserService {
	return &UserService{users: users, inters: inters, postSvc: postSvc, commentSvc: commentSvc, cache: cch}
}

// List 用户列表，支持 keyword 关键词过滤。
func (s *UserService) List(keyword string, page, size int) ([]dto.User, int64, error) {
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

// Get 用户主页（Cache Aside）。
func (s *UserService) Get(id int64) (*dto.User, error) {
	ctx := context.Background()
	key := cache.UserProfileKey(id)
	var cached dto.User
	if s.cache.GetJSON(ctx, key, &cached) {
		return &cached, nil
	}
	u, err := s.users.FindByID(id)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if u == nil {
		return nil, errs.ErrNotFound
	}
	res := dto.ToUser(u, false, 0)
	s.cache.SetJSON(ctx, key, res, userProfileTTL)
	return &res, nil
}

// UpdateProfile 修改我的资料。
func (s *UserService) UpdateProfile(userID int64, in dto.UpdateProfileRequest) (*dto.User, error) {
	updates := map[string]any{"updated_at": time.Now()}
	if in.Nickname != nil {
		updates["nickname"] = *in.Nickname
	}
	if in.Avatar != nil {
		updates["avatar"] = *in.Avatar
	}
	if in.Bio != nil {
		updates["bio"] = *in.Bio
	}
	if err := s.users.Update(userID, updates); err != nil {
		return nil, errs.ErrInternal
	}
	s.cache.Del(context.Background(), cache.UserProfileKey(userID))
	return s.Get(userID)
}

// MyPosts 我的帖子（含隐藏、草稿等状态）。
func (s *UserService) MyPosts(userID int64, page, size int) ([]dto.Post, int64, error) {
	return s.postSvc.List(ListFilter{
		AuthorID:      userID,
		ViewerID:      userID,
		Status:        "all",
		IncludeHidden: true,
		Page:          page,
		PageSize:      size,
	})
}

// MyComments 我的评论。
func (s *UserService) MyComments(userID int64, page, size int) ([]dto.MyCommentItem, int64, error) {
	return s.commentSvc.ListMine(userID, page, size)
}

// MyLikedPosts 我点赞的帖子。
func (s *UserService) MyLikedPosts(userID int64, page, size int) ([]dto.Post, int64, error) {
	return s.postsByRelation("post_likes", userID, page, size)
}

// MyFavoritePosts 我收藏的帖子。
func (s *UserService) MyFavoritePosts(userID int64, page, size int) ([]dto.Post, int64, error) {
	return s.postsByRelation("post_favorites", userID, page, size)
}

// postsByRelation 根据点赞 / 收藏关系表查询帖子，保持关系时间倒序。
func (s *UserService) postsByRelation(table string, userID int64, page, size int) ([]dto.Post, int64, error) {
	page, size = normPage(page, size)
	ids := s.inters.PostIDsByUser(table, userID)
	total := int64(len(ids))
	if total == 0 {
		return []dto.Post{}, 0, nil
	}
	start := offset(page, size)
	if start > len(ids) {
		start = len(ids)
	}
	end := start + size
	if end > len(ids) {
		end = len(ids)
	}
	list, err := s.postSvc.AssembleByIDs(ids[start:end], userID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
