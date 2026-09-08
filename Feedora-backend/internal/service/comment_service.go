package service

import (
	"context"
	"strings"
	"time"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/utils"
)

// CommentService 评论业务逻辑。
type CommentService struct {
	comments *repository.CommentRepository
	posts    *repository.PostRepository
	users    *repository.UserRepository
	inters   *repository.InteractionRepository
	producer event.Producer
	cache    *cache.Cache
}

func NewCommentService(
	comments *repository.CommentRepository,
	posts *repository.PostRepository,
	users *repository.UserRepository,
	inters *repository.InteractionRepository,
	producer event.Producer,
	cch *cache.Cache,
) *CommentService {
	return &CommentService{comments: comments, posts: posts, users: users, inters: inters, producer: producer, cache: cch}
}

// ListByPost 查询帖子的评论树（根评论 + 其下回复）。
func (s *CommentService) ListByPost(postID, viewerID int64) ([]dto.Comment, error) {
	rows, err := s.comments.ListByPost(postID)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if len(rows) == 0 {
		return []dto.Comment{}, nil
	}
	userIDs := make([]int64, 0, len(rows))
	commentIDs := make([]int64, 0, len(rows))
	for i := range rows {
		userIDs = append(userIDs, rows[i].UserID)
		commentIDs = append(commentIDs, rows[i].ID)
	}
	users, err := s.users.FindByIDs(userIDs)
	if err != nil {
		return nil, errs.ErrInternal
	}
	likedSet := s.inters.CommentLikedSet(viewerID, commentIDs)

	toDTO := func(m *model.Comment) dto.Comment {
		return dto.Comment{
			CommentID: m.ID,
			PostID:    m.PostID,
			UserID:    m.UserID,
			User:      dto.ToUserSummary(users[m.UserID]),
			Content:   m.Content,
			LikeCount: m.LikeCount,
			Liked:     likedSet[m.ID],
			Status:    m.Status,
			Replies:   []dto.Comment{},
			CreatedAt: utils.FormatTime(m.CreatedAt),
			UpdatedAt: utils.FormatTime(m.UpdatedAt),
		}
	}

	rootIndex := map[int64]int{}
	roots := make([]dto.Comment, 0)
	for i := range rows {
		if rows[i].ParentID == 0 {
			roots = append(roots, toDTO(&rows[i]))
			rootIndex[rows[i].ID] = len(roots) - 1
		}
	}
	for i := range rows {
		if rows[i].ParentID == 0 {
			continue
		}
		root := rows[i].RootID
		if root == 0 {
			root = rows[i].ParentID
		}
		if idx, ok := rootIndex[root]; ok {
			roots[idx].Replies = append(roots[idx].Replies, toDTO(&rows[i]))
		}
	}
	return roots, nil
}

// Create 发表根评论。
func (s *CommentService) Create(postID, userID int64, content string) (*dto.Comment, error) {
	return s.create(postID, userID, 0, 0, nil, content)
}

// Reply 回复某条评论。
func (s *CommentService) Reply(parentID, userID int64, content string, replyToUserID *int64) (*dto.Comment, error) {
	parent, err := s.comments.FindByID(parentID)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if parent == nil {
		return nil, errs.ErrCommentNotFound
	}
	root := parent.RootID
	if root == 0 {
		root = parent.ID
	}
	return s.create(parent.PostID, userID, parentID, root, replyToUserID, content)
}

func (s *CommentService) create(postID, userID, parentID, rootID int64, replyTo *int64, content string) (*dto.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errs.ErrParams
	}
	p, err := s.posts.FindByID(postID)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if p == nil {
		return nil, errs.ErrPostNotFound
	}
	if p.Status != model.PostPublished {
		return nil, errs.ErrPostInvisible
	}
	now := time.Now()
	m := &model.Comment{
		PostID:        postID,
		UserID:        userID,
		ParentID:      parentID,
		RootID:        rootID,
		ReplyToUserID: replyTo,
		Content:       content,
		Status:        model.CommentNormal,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.comments.CreateWithCounters(m); err != nil {
		return nil, errs.ErrInternal
	}
	s.cache.Del(context.Background(), cache.PostDetailKey(postID))
	s.producer.Publish(event.TopicComment, event.CommentCreated, m.ID, userID, map[string]any{"postId": postID, "postAuthorId": p.AuthorID})

	u, _ := s.users.FindByID(userID)
	res := dto.Comment{
		CommentID: m.ID, PostID: postID, UserID: userID, User: dto.ToUserSummary(u),
		Content: content, LikeCount: 0, Liked: false, Status: model.CommentNormal, Replies: []dto.Comment{},
		CreatedAt: utils.FormatTime(now), UpdatedAt: utils.FormatTime(now),
	}
	return &res, nil
}

// Delete 删除评论（软删除），作者本人或管理员可操作。
func (s *CommentService) Delete(commentID, userID int64, isAdmin bool) error {
	m, err := s.comments.FindByID(commentID)
	if err != nil {
		return errs.ErrInternal
	}
	if m == nil {
		return errs.ErrCommentNotFound
	}
	if m.UserID != userID && !isAdmin {
		return errs.ErrForbidden
	}
	s.comments.DeleteWithCounters(m)
	s.cache.Del(context.Background(), cache.PostDetailKey(m.PostID))
	return nil
}

// SetLike 评论点赞 / 取消点赞，唯一索引保证幂等。
func (s *CommentService) SetLike(commentID, userID int64, like bool) error {
	m, err := s.comments.FindByID(commentID)
	if err != nil {
		return errs.ErrInternal
	}
	if m == nil {
		return errs.ErrCommentNotFound
	}
	if like {
		if s.inters.AddCommentLike(commentID, userID) {
			s.comments.IncLikeCount(commentID, 1)
		}
	} else {
		if s.inters.RemoveCommentLike(commentID, userID) {
			s.comments.IncLikeCount(commentID, -1)
		}
	}
	return nil
}

// ListMine 查询我的评论列表。
func (s *CommentService) ListMine(userID int64, page, size int) ([]dto.MyCommentItem, int64, error) {
	page, size = normPage(page, size)
	rows, total, err := s.comments.ListMine(userID, offset(page, size), size)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	postIDs := make([]int64, 0, len(rows))
	for i := range rows {
		postIDs = append(postIDs, rows[i].PostID)
	}
	titles := map[int64]string{}
	if len(postIDs) > 0 {
		posts, _ := s.posts.FindByIDs(postIDs, false)
		for i := range posts {
			titles[posts[i].ID] = posts[i].Title
		}
	}
	list := make([]dto.MyCommentItem, 0, len(rows))
	for i := range rows {
		title := titles[rows[i].PostID]
		if title == "" {
			title = "原帖已删除"
		}
		list = append(list, dto.MyCommentItem{
			CommentID: rows[i].ID, PostID: rows[i].PostID, PostTitle: title,
			Content: rows[i].Content, LikeCount: rows[i].LikeCount, Status: rows[i].Status,
			CreatedAt: utils.FormatTime(rows[i].CreatedAt),
		})
	}
	return list, total, nil
}
