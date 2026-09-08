package service

import (
	"context"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

// FollowService 关注体系：关注 / 取关、关注与粉丝列表、关注动态聚合。
type FollowService struct {
	follows  *repository.FollowRepository
	users    *repository.UserRepository
	posts    *PostService
	producer event.Producer
	cache    *cache.Cache
}

func NewFollowService(
	follows *repository.FollowRepository,
	users *repository.UserRepository,
	posts *PostService,
	producer event.Producer,
	cch *cache.Cache,
) *FollowService {
	return &FollowService{follows: follows, users: users, posts: posts, producer: producer, cache: cch}
}

// invalidateProfiles 关注关系变化后失效双方的资料缓存（粉丝 / 关注数）。
func (s *FollowService) invalidateProfiles(a, b int64) {
	ctx := context.Background()
	s.cache.Del(ctx, cache.UserProfileKey(a))
	s.cache.Del(ctx, cache.UserProfileKey(b))
}

// Follow 关注用户（幂等，重复关注不重复计数）。
func (s *FollowService) Follow(followerID, followeeID int64) error {
	if followerID <= 0 {
		return errs.ErrUnauth
	}
	if followerID == followeeID {
		return errs.ErrParams
	}
	u, err := s.users.FindByID(followeeID)
	if err != nil {
		return errs.ErrInternal
	}
	if u == nil {
		return errs.ErrNotFound
	}
	created, err := s.follows.Follow(followerID, followeeID)
	if err != nil {
		return errs.ErrInternal
	}
	if created {
		s.users.IncColumn(followerID, "following_count", 1)
		s.users.IncColumn(followeeID, "follower_count", 1)
		s.invalidateProfiles(followerID, followeeID)
		s.producer.Publish(event.TopicUser, event.UserFollowed, followeeID, followerID, nil)
	}
	return nil
}

// Unfollow 取消关注（幂等）。
func (s *FollowService) Unfollow(followerID, followeeID int64) error {
	if followerID <= 0 {
		return errs.ErrUnauth
	}
	removed, err := s.follows.Unfollow(followerID, followeeID)
	if err != nil {
		return errs.ErrInternal
	}
	if removed {
		s.users.IncColumn(followerID, "following_count", -1)
		s.users.IncColumn(followeeID, "follower_count", -1)
		s.invalidateProfiles(followerID, followeeID)
		s.producer.Publish(event.TopicUser, event.UserUnfollowed, followeeID, followerID, nil)
	}
	return nil
}

// State 查询 viewer 是否关注了 target（未登录或自己返回 false）。
func (s *FollowService) State(viewerID, targetID int64) (bool, error) {
	if viewerID <= 0 || viewerID == targetID {
		return false, nil
	}
	ok, err := s.follows.IsFollowing(viewerID, targetID)
	if err != nil {
		return false, errs.ErrInternal
	}
	return ok, nil
}

// Following 关注的人分页列表。
func (s *FollowService) Following(followerID int64, page, size int) ([]dto.User, int64, error) {
	page, size = normPage(page, size)
	users, total, err := s.follows.ListFollowing(followerID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.User, 0, len(users))
	for i := range users {
		list = append(list, dto.ToUser(&users[i], false, 0))
	}
	return list, total, nil
}

// Followers 粉丝分页列表。
func (s *FollowService) Followers(followeeID int64, page, size int) ([]dto.User, int64, error) {
	page, size = normPage(page, size)
	users, total, err := s.follows.ListFollowers(followeeID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.User, 0, len(users))
	for i := range users {
		list = append(list, dto.ToUser(&users[i], false, 0))
	}
	return list, total, nil
}

// FollowingFeed 关注动态：all/user 看关注作者的帖子，circle 看加入圈子的帖子，
// topic 看自己参与话题下的帖子。
func (s *FollowService) FollowingFeed(viewerID int64, tab string, page, size int) ([]dto.FollowingFeedItem, int64, error) {
	page, size = normPage(page, size)
	filter := ListFilter{ViewerID: viewerID, Page: page, PageSize: size, Sort: "latest"}
	feedType := "post"
	switch tab {
	case "circle":
		circles, _, err := s.follows.JoinedCircles(viewerID, 0, 100)
		if err != nil {
			return nil, 0, errs.ErrInternal
		}
		ids := make([]int64, 0, len(circles))
		for i := range circles {
			ids = append(ids, circles[i].ID)
		}
		filter.CircleIDs = ids
		feedType = "circle"
	case "topic":
		topics, _, err := s.follows.ParticipatedTopics(viewerID, 0, 100)
		if err != nil {
			return nil, 0, errs.ErrInternal
		}
		ids := make([]int64, 0, len(topics))
		for i := range topics {
			ids = append(ids, topics[i].ID)
		}
		filter.TopicIDs = ids
		feedType = "topic"
	default: // all / user
		filter.FeedType = "following"
	}
	posts, total, err := s.posts.List(filter)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.FollowingFeedItem, 0, len(posts))
	for i := range posts {
		p := posts[i]
		list = append(list, dto.FollowingFeedItem{
			FeedID:       p.PostID,
			FeedType:     feedType,
			Title:        p.Title,
			Summary:      p.Summary,
			TargetID:     p.PostID,
			TargetURL:    "/posts/" + itoa(p.PostID),
			SourceName:   p.Author.Nickname,
			SourceAvatar: p.Author.Avatar,
			CreatedAt:    p.CreatedAt,
		})
	}
	return list, total, nil
}

// FollowingCircles 关注中心「关注的圈子」：用户加入的圈子。
func (s *FollowService) FollowingCircles(userID int64, page, size int) ([]dto.Circle, int64, error) {
	page, size = normPage(page, size)
	circles, total, err := s.follows.JoinedCircles(userID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.Circle, 0, len(circles))
	for i := range circles {
		list = append(list, dto.ToCircle(&circles[i], nil, nil))
	}
	return list, total, nil
}

// FollowingTopics 关注中心「关注的话题」：用户发帖参与的话题。
func (s *FollowService) FollowingTopics(userID int64, page, size int) ([]dto.Topic, int64, error) {
	page, size = normPage(page, size)
	topics, total, err := s.follows.ParticipatedTopics(userID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.Topic, 0, len(topics))
	for i := range topics {
		list = append(list, dto.ToTopic(&topics[i]))
	}
	return list, total, nil
}

// FollowingTags 关注中心「关注的标签」：用户发帖使用的标签。
func (s *FollowService) FollowingTags(userID int64, page, size int) ([]dto.ContentTag, int64, error) {
	page, size = normPage(page, size)
	tags, total, err := s.follows.UsedTags(userID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.ContentTag, 0, len(tags))
	for i := range tags {
		list = append(list, dto.ToContentTag(&tags[i]))
	}
	return list, total, nil
}
