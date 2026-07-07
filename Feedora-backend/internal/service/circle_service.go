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
)

const circleDetailTTL = 15 * time.Minute

// CircleService 圈子业务逻辑。
type CircleService struct {
	circles  *repository.CircleRepository
	users    *repository.UserRepository
	postSvc  *PostService
	producer event.Producer
	cache    *cache.Cache
}

func NewCircleService(
	circles *repository.CircleRepository,
	users *repository.UserRepository,
	postSvc *PostService,
	producer event.Producer,
	cch *cache.Cache,
) *CircleService {
	return &CircleService{circles: circles, users: users, postSvc: postSvc, producer: producer, cache: cch}
}

// List 按范围分页查询圈子。
func (s *CircleService) List(scope, keyword, category, sort string, viewerID int64, page, size int) ([]dto.Circle, int64, error) {
	page, size = normPage(page, size)
	rows, total, err := s.circles.List(repository.CircleFilter{
		Scope: scope, Keyword: keyword, Category: category, Sort: sort,
		ViewerID: viewerID, Offset: offset(page, size), Limit: size,
	})
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list, err := s.assemble(rows, viewerID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// Get 圈子详情（Cache Aside）。详情主体缓存，按当前用户叠加成员身份。
func (s *CircleService) Get(circleID, viewerID int64) (*dto.Circle, error) {
	ctx := context.Background()
	key := cache.CircleDetailKey(circleID)
	var detail dto.Circle
	if s.cache.GetJSON(ctx, key, &detail) {
		s.overlayMembership(&detail, circleID, viewerID)
		return &detail, nil
	}
	c, err := s.circles.FindByID(circleID)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if c == nil {
		return nil, errs.ErrCircleNotFound
	}
	list, err := s.assemble([]model.Circle{*c}, 0)
	if err != nil {
		return nil, err
	}
	detail = list[0]
	s.cache.SetJSON(ctx, key, detail, circleDetailTTL)
	s.overlayMembership(&detail, circleID, viewerID)
	return &detail, nil
}

// overlayMembership 叠加当前用户在圈子中的身份状态。
func (s *CircleService) overlayMembership(c *dto.Circle, circleID, viewerID int64) {
	if viewerID <= 0 {
		return
	}
	m, _ := s.circles.FindMember(circleID, viewerID)
	if m != nil && m.Status != model.CircleMemberRemoved {
		c.IsJoined = true
		c.MyRole = m.Role
		c.MyStatus = m.Status
	}
}

// invalidate 失效圈子详情缓存。
func (s *CircleService) invalidate(circleID int64) {
	s.cache.Del(context.Background(), cache.CircleDetailKey(circleID))
}

// assemble 装配圈子 DTO，加载圈主与当前用户成员身份。
func (s *CircleService) assemble(rows []model.Circle, viewerID int64) ([]dto.Circle, error) {
	result := make([]dto.Circle, 0, len(rows))
	if len(rows) == 0 {
		return result, nil
	}
	ownerIDs := make([]int64, 0, len(rows))
	circleIDs := make([]int64, 0, len(rows))
	for i := range rows {
		ownerIDs = append(ownerIDs, rows[i].OwnerID)
		circleIDs = append(circleIDs, rows[i].ID)
	}
	owners, err := s.users.FindByIDs(ownerIDs)
	if err != nil {
		return nil, errs.ErrInternal
	}
	members := s.circles.FindMembersByViewer(viewerID, circleIDs)
	for i := range rows {
		result = append(result, dto.ToCircle(&rows[i], owners[rows[i].OwnerID], members[rows[i].ID]))
	}
	return result, nil
}

// Create 创建圈子，创建者自动成为 owner。
func (s *CircleService) Create(ownerID int64, in dto.CreateCircleRequest) (*dto.Circle, error) {
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return nil, errs.ErrParams
	}
	if s.circles.CountByName(in.Name) > 0 {
		return nil, errs.New(409, "圈子名称已存在")
	}
	if in.JoinType == "" {
		in.JoinType = "direct"
	}
	if in.PostPermission == "" {
		in.PostPermission = "all"
	}
	now := time.Now()
	c := &model.Circle{
		OwnerID: ownerID, Name: in.Name, Avatar: in.Avatar, Description: in.Description,
		Category: in.Category, JoinType: in.JoinType, PostPermission: in.PostPermission,
		Rules: in.Rules, Status: model.CircleMemberNormal, MemberCount: 1,
		CreatedAt: now, UpdatedAt: now,
	}
	owner := &model.CircleMember{
		UserID: ownerID, Role: model.CircleRoleOwner, Status: model.CircleMemberNormal,
		JoinedAt: now, UpdatedAt: now,
	}
	if err := s.circles.CreateWithOwner(c, owner); err != nil {
		return nil, errs.ErrInternal
	}
	s.producer.Publish(event.TopicCircle, event.CircleCreated, c.ID, ownerID, map[string]any{"name": c.Name})
	return s.Get(c.ID, ownerID)
}

// Join 加入圈子。direct 类型直接加入。
func (s *CircleService) Join(circleID, userID int64) error {
	c, err := s.circles.FindByID(circleID)
	if err != nil {
		return errs.ErrInternal
	}
	if c == nil {
		return errs.ErrCircleNotFound
	}
	existing, err := s.circles.FindMember(circleID, userID)
	if err != nil {
		return errs.ErrInternal
	}
	if existing != nil {
		if existing.Status == model.CircleMemberRemoved {
			s.circles.UpdateMember(circleID, userID, map[string]any{"status": model.CircleMemberNormal, "updated_at": time.Now()})
			s.circles.IncMemberCount(circleID, 1)
			s.invalidate(circleID)
		}
		return nil
	}
	now := time.Now()
	if err := s.circles.CreateMember(&model.CircleMember{
		CircleID: circleID, UserID: userID, Role: model.CircleRoleMember, Status: model.CircleMemberNormal,
		JoinedAt: now, UpdatedAt: now,
	}); err != nil {
		return errs.ErrInternal
	}
	s.circles.IncMemberCount(circleID, 1)
	s.invalidate(circleID)
	s.producer.Publish(event.TopicCircle, event.CircleJoined, circleID, userID, nil)
	return nil
}

// Leave 退出圈子。owner 不能退出。
func (s *CircleService) Leave(circleID, userID int64) error {
	m, err := s.circles.FindMember(circleID, userID)
	if err != nil {
		return errs.ErrInternal
	}
	if m == nil {
		return errs.ErrCircleNotJoined
	}
	if m.Role == model.CircleRoleOwner {
		return errs.New(422, "圈主不能退出圈子")
	}
	s.circles.DeleteMember(m)
	s.circles.IncMemberCount(circleID, -1)
	s.invalidate(circleID)
	return nil
}

// Members 圈子成员列表。
func (s *CircleService) Members(circleID int64, page, size int) ([]dto.CircleMember, int64, error) {
	page, size = normPage(page, size)
	rows, total := s.circles.Members(circleID, offset(page, size), size)
	userIDs := make([]int64, 0, len(rows))
	for i := range rows {
		userIDs = append(userIDs, rows[i].UserID)
	}
	users, err := s.users.FindByIDs(userIDs)
	if err != nil {
		return nil, 0, errs.ErrInternal
	}
	list := make([]dto.CircleMember, 0, len(rows))
	for i := range rows {
		list = append(list, dto.ToCircleMember(&rows[i], users[rows[i].UserID]))
	}
	return list, total, nil
}

// Posts 圈子内帖子。
func (s *CircleService) Posts(circleID, viewerID int64, sort string, page, size int) ([]dto.Post, int64, error) {
	return s.postSvc.List(ListFilter{CircleID: circleID, ViewerID: viewerID, Sort: sort, Page: page, PageSize: size})
}

// SetRole 修改成员角色，仅圈主可操作。
func (s *CircleService) SetRole(circleID, targetUserID, operatorID int64, role string) error {
	op, err := s.circles.FindMember(circleID, operatorID)
	if err != nil {
		return errs.ErrInternal
	}
	if op == nil || op.Role != model.CircleRoleOwner {
		return errs.ErrForbidden
	}
	return s.circles.UpdateMember(circleID, targetUserID, map[string]any{"role": role, "updated_at": time.Now()})
}

// Mute 禁言成员。days<=0 表示永久。
func (s *CircleService) Mute(circleID, targetUserID, operatorID int64, days int, reason string) error {
	if err := s.requireManage(circleID, operatorID); err != nil {
		return err
	}
	updates := map[string]any{"status": model.CircleMemberMuted, "mute_reason": reason, "updated_at": time.Now()}
	if days > 0 {
		until := time.Now().AddDate(0, 0, days)
		updates["muted_until"] = until
	}
	return s.circles.UpdateMember(circleID, targetUserID, updates)
}

// Unmute 解除禁言。
func (s *CircleService) Unmute(circleID, targetUserID, operatorID int64) error {
	if err := s.requireManage(circleID, operatorID); err != nil {
		return err
	}
	return s.circles.UpdateMember(circleID, targetUserID, map[string]any{
		"status": model.CircleMemberNormal, "mute_reason": "", "muted_until": nil, "updated_at": time.Now(),
	})
}

// Remove 移除成员。
func (s *CircleService) Remove(circleID, targetUserID, operatorID int64) error {
	if err := s.requireManage(circleID, operatorID); err != nil {
		return err
	}
	m, err := s.circles.FindMember(circleID, targetUserID)
	if err != nil {
		return errs.ErrInternal
	}
	if m == nil {
		return errs.ErrNotFound
	}
	if m.Role == model.CircleRoleOwner {
		return errs.New(422, "不能移除圈主")
	}
	s.circles.UpdateMember(circleID, targetUserID, map[string]any{"status": model.CircleMemberRemoved, "updated_at": time.Now()})
	s.circles.IncMemberCount(circleID, -1)
	s.invalidate(circleID)
	return nil
}

// requireManage 校验操作者是否具备圈子管理权限（owner / moderator）。
func (s *CircleService) requireManage(circleID, operatorID int64) error {
	op, err := s.circles.FindMember(circleID, operatorID)
	if err != nil {
		return errs.ErrInternal
	}
	if op == nil || (op.Role != model.CircleRoleOwner && op.Role != model.CircleRoleModerator) {
		return errs.ErrForbidden
	}
	return nil
}
