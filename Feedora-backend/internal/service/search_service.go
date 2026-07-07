package service

import (
	"context"
	"strings"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	"github.com/feedora/backend/internal/search"
	errs "github.com/feedora/backend/pkg/errors"
)

// SearchService 搜索业务逻辑（读 ES + Redis 热词）。
type SearchService struct {
	sc      *search.Client
	posts   *repository.PostRepository
	users   *repository.UserRepository
	topics  *repository.TopicRepository
	circles *repository.CircleRepository
	postSvc *PostService
	cache   *cache.Cache
}

func NewSearchService(
	sc *search.Client,
	posts *repository.PostRepository,
	users *repository.UserRepository,
	topics *repository.TopicRepository,
	circles *repository.CircleRepository,
	postSvc *PostService,
	cch *cache.Cache,
) *SearchService {
	return &SearchService{sc: sc, posts: posts, users: users, topics: topics, circles: circles, postSvc: postSvc, cache: cch}
}

// Search 全局 / 分类搜索。返回混合类型的扁平列表与总数，
// 由前端按 postId/userId/topicId/circleId 键区分渲染。
func (s *SearchService) Search(keyword, typ string, page, size int, viewerID int64) ([]any, int64, error) {
	list := []any{}
	if s.sc == nil {
		return list, 0, errs.New(500, "搜索服务未启用")
	}
	page, size = normPage(page, size)
	from := offset(page, size)
	ctx := context.Background()
	var total int64

	// 记录热门搜索词。
	if kw := strings.TrimSpace(keyword); kw != "" {
		s.cache.ZIncr(ctx, cache.HotKeywordsKey, kw, 1)
	}

	if typ == "" || typ == "all" || typ == "post" {
		q := search.MatchQuery(keyword, []string{"title^2", "content", "summary", "tagNames"},
			map[string]any{"status.keyword": "published", "visibility.keyword": "public"})
		ids, t, err := s.sc.SearchIDs(ctx, s.sc.PostIndex(), q, from, size)
		if err == nil {
			if posts, e := s.postSvc.AssembleByIDs(ids, viewerID); e == nil {
				for i := range posts {
					list = append(list, posts[i])
				}
			}
			total += t
		}
	}
	if typ == "" || typ == "all" || typ == "user" {
		q := search.MatchQuery(keyword, []string{"nickname^2", "bio"}, map[string]any{"status.keyword": "normal"})
		ids, t, err := s.sc.SearchIDs(ctx, s.sc.UserIndex(), q, from, size)
		if err == nil {
			users, _ := s.users.FindByIDs(ids)
			for _, id := range ids {
				if u := users[id]; u != nil {
					list = append(list, dto.ToUser(u, false, 0))
				}
			}
			total += t
		}
	}
	if typ == "" || typ == "all" || typ == "circle" {
		q := search.MatchQuery(keyword, []string{"name^2", "description", "category"}, map[string]any{"status.keyword": "normal"})
		ids, t, err := s.sc.SearchIDs(ctx, s.sc.CircleIndex(), q, from, size)
		if err == nil {
			cs := s.circles.FindByIDs(ids)
			for _, id := range ids {
				if c := cs[id]; c != nil {
					owner, _ := s.users.FindByID(c.OwnerID)
					list = append(list, dto.ToCircle(c, owner, nil))
				}
			}
			total += t
		}
	}
	if typ == "" || typ == "all" || typ == "topic" {
		q := search.MatchQuery(keyword, []string{"name^2", "description"}, map[string]any{"status.keyword": "enabled"})
		ids, t, err := s.sc.SearchIDs(ctx, s.sc.TopicIndex(), q, from, size)
		if err == nil {
			for _, id := range ids {
				if tp, _ := s.topics.FindByID(id); tp != nil {
					list = append(list, dto.ToTopic(tp))
				}
			}
			total += t
		}
	}
	return list, total, nil
}

// Suggest 搜索联想，返回帖子标题联想。
func (s *SearchService) Suggest(keyword string) ([]dto.SuggestItem, error) {
	items := []dto.SuggestItem{}
	if s.sc == nil || strings.TrimSpace(keyword) == "" {
		return items, nil
	}
	ctx := context.Background()
	q := search.MatchQuery(keyword, []string{"title^2", "summary"},
		map[string]any{"status.keyword": "published", "visibility.keyword": "public"})
	ids, _, err := s.sc.SearchIDs(ctx, s.sc.PostIndex(), q, 0, 8)
	if err != nil {
		return items, nil
	}
	posts, _ := s.posts.FindByIDs(ids, true)
	byID := map[int64]*model.Post{}
	for i := range posts {
		byID[posts[i].ID] = &posts[i]
	}
	for _, id := range ids {
		if p := byID[id]; p != nil {
			items = append(items, dto.SuggestItem{
				Type: "post", Title: p.Title, TargetID: p.ID,
				TargetURL: "/posts/" + itoa(p.ID),
			})
		}
	}
	return items, nil
}

// HotKeywords 热门搜索词 Top10（前端需要字符串数组）。
func (s *SearchService) HotKeywords() []string {
	items := []string{}
	if !s.cache.Enabled() {
		return items
	}
	for _, z := range s.cache.ZTop(context.Background(), cache.HotKeywordsKey, 0, 10) {
		if kw, ok := z.Member.(string); ok && kw != "" {
			items = append(items, kw)
		}
	}
	return items
}
