package service

import (
	"context"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	"gorm.io/gorm"
)

// RankService 热门榜单与排行榜（读 Redis ZSet，为空时回退 DB 排序）。
type RankService struct {
	db      *gorm.DB
	posts   *repository.PostRepository
	users   *repository.UserRepository
	topics  *repository.TopicRepository
	circles *repository.CircleRepository
	cache   *cache.Cache
}

func NewRankService(
	db *gorm.DB,
	posts *repository.PostRepository,
	users *repository.UserRepository,
	topics *repository.TopicRepository,
	circles *repository.CircleRepository,
	cch *cache.Cache,
) *RankService {
	return &RankService{db: db, posts: posts, users: users, topics: topics, circles: circles, cache: cch}
}

// HotRanks 热门榜单：post / circle / topic。
func (s *RankService) HotRanks(rankType, timeRange string, page, size int) []dto.HotRankItem {
	if rankType == "" {
		rankType = "post"
	}
	if timeRange == "" {
		timeRange = "today"
	}
	page, size = normPage(page, size)
	ctx := context.Background()

	// 优先从 Redis ZSet 取 ID+分数。
	ids := []int64{}
	scores := map[int64]float64{}
	for _, z := range s.cache.ZTop(ctx, cache.RankKey(rankType, timeRange), offset(page, size), size) {
		if idStr, ok := z.Member.(string); ok {
			if id := parseInt(idStr); id > 0 {
				ids = append(ids, id)
				scores[id] = z.Score
			}
		}
	}
	// ZSet 为空时回退 DB 排序，保证页面有数据。
	if len(ids) == 0 {
		return s.hotFromDB(rankType, offset(page, size), size)
	}
	return s.buildHot(rankType, ids, scores, offset(page, size))
}

func (s *RankService) buildHot(rankType string, ids []int64, scores map[int64]float64, baseRank int) []dto.HotRankItem {
	items := make([]dto.HotRankItem, 0, len(ids))
	switch rankType {
	case "circle":
		cs := s.circles.FindByIDs(ids)
		for i, id := range ids {
			if c := cs[id]; c != nil {
				items = append(items, dto.HotRankItem{
					Rank: baseRank + i + 1, Score: scores[id], TargetID: id, CircleID: id,
					Name: c.Name, MemberCount: c.MemberCount, PostCount: c.PostCount, FeaturedPostCount: c.FeaturedCount,
				})
			}
		}
	case "topic":
		for i, id := range ids {
			if t, _ := s.topics.FindByID(id); t != nil {
				items = append(items, dto.HotRankItem{
					Rank: baseRank + i + 1, Score: scores[id], TargetID: id, TopicID: id,
					Name: t.Name, PostCount: t.PostCount, ParticipantCount: t.ParticipantCount,
				})
			}
		}
	default: // post
		posts, _ := s.posts.FindByIDs(ids, false)
		byID := map[int64]model.Post{}
		authorIDs := []int64{}
		for i := range posts {
			byID[posts[i].ID] = posts[i]
			authorIDs = append(authorIDs, posts[i].AuthorID)
		}
		authors, _ := s.users.FindByIDs(authorIDs)
		for i, id := range ids {
			if p, ok := byID[id]; ok {
				name := ""
				if a := authors[p.AuthorID]; a != nil {
					name = a.Nickname
				}
				items = append(items, dto.HotRankItem{
					Rank: baseRank + i + 1, Score: scores[id], TargetID: id, PostID: id,
					Title: p.Title, AuthorName: name, HotScore: p.HotScore,
					LikeCount: p.LikeCount, CommentCount: p.CommentCount,
				})
			}
		}
	}
	return items
}

func (s *RankService) hotFromDB(rankType string, off, size int) []dto.HotRankItem {
	items := []dto.HotRankItem{}
	switch rankType {
	case "circle":
		var cs []model.Circle
		s.db.Where("status = ?", model.CircleMemberNormal).Order("member_count DESC, post_count DESC").Offset(off).Limit(size).Find(&cs)
		for i := range cs {
			items = append(items, dto.HotRankItem{
				Rank: off + i + 1, Score: float64(cs[i].MemberCount), TargetID: cs[i].ID, CircleID: cs[i].ID,
				Name: cs[i].Name, MemberCount: cs[i].MemberCount, PostCount: cs[i].PostCount, FeaturedPostCount: cs[i].FeaturedCount,
			})
		}
	case "topic":
		var ts []model.Topic
		s.db.Where("status = ?", model.StatusEnabled).Order("participant_count DESC, post_count DESC").Offset(off).Limit(size).Find(&ts)
		for i := range ts {
			items = append(items, dto.HotRankItem{
				Rank: off + i + 1, Score: float64(ts[i].ParticipantCount), TargetID: ts[i].ID, TopicID: ts[i].ID,
				Name: ts[i].Name, PostCount: ts[i].PostCount, ParticipantCount: ts[i].ParticipantCount,
			})
		}
	default:
		var ps []model.Post
		s.db.Where("status = ?", model.PostPublished).Order("hot_score DESC, like_count DESC, created_at DESC").Offset(off).Limit(size).Find(&ps)
		authorIDs := make([]int64, 0, len(ps))
		for i := range ps {
			authorIDs = append(authorIDs, ps[i].AuthorID)
		}
		authors, _ := s.users.FindByIDs(authorIDs)
		for i := range ps {
			name := ""
			if a := authors[ps[i].AuthorID]; a != nil {
				name = a.Nickname
			}
			items = append(items, dto.HotRankItem{
				Rank: off + i + 1, Score: float64(ps[i].HotScore), TargetID: ps[i].ID, PostID: ps[i].ID,
				Title: ps[i].Title, AuthorName: name, HotScore: ps[i].HotScore,
				LikeCount: ps[i].LikeCount, CommentCount: ps[i].CommentCount,
			})
		}
	}
	return items
}

// Rankings 用户 / 圈子排行榜。type: active/contribution/creator/circle。
func (s *RankService) Rankings(rankType, timeRange string, page, size int, currentUserID int64) []dto.RankingItem {
	if rankType == "" {
		rankType = "active"
	}
	if timeRange == "" {
		timeRange = "all"
	}
	page, size = normPage(page, size)
	off := offset(page, size)

	// 阶段一/二初期 ZSet 可能为空，直接以 DB 排序给出稳定结果。
	items := []dto.RankingItem{}
	if rankType == "circle" {
		var cs []model.Circle
		s.db.Where("status = ?", model.CircleMemberNormal).Order("member_count DESC").Offset(off).Limit(size).Find(&cs)
		for i := range cs {
			items = append(items, dto.RankingItem{
				Rank: off + i + 1, TargetID: cs[i].ID, TargetType: "circle", Name: cs[i].Name,
				Avatar: cs[i].Avatar, MemberCount: cs[i].MemberCount, Score: float64(cs[i].MemberCount),
			})
		}
		return items
	}
	order := "point_count DESC"
	switch rankType {
	case "creator", "contribution":
		order = "post_count DESC, like_count DESC"
	case "active":
		order = "point_count DESC, post_count DESC"
	}
	var us []model.User
	s.db.Where("status = ?", model.UserNormal).Order(order).Offset(off).Limit(size).Find(&us)
	for i := range us {
		score := us[i].PointCount
		if rankType == "creator" || rankType == "contribution" {
			score = us[i].PostCount
		}
		items = append(items, dto.RankingItem{
			Rank: off + i + 1, TargetID: us[i].ID, TargetType: "user", Name: us[i].Nickname,
			Avatar: us[i].Avatar, Level: us[i].Level, LevelName: dto.LevelNameOf(us[i].Level), Points: us[i].PointCount,
			PostCount: us[i].PostCount, LikeReceivedCount: us[i].LikeCount, Score: float64(score),
			IsCurrentUser: us[i].ID == currentUserID,
		})
	}
	return items
}
