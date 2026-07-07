package repository

import (
	"errors"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// CircleRepository 圈子数据访问。
type CircleRepository struct {
	db *gorm.DB
}

func NewCircleRepository(db *gorm.DB) *CircleRepository {
	return &CircleRepository{db: db}
}

// CircleFilter 圈子列表过滤条件。
type CircleFilter struct {
	Scope    string
	Keyword  string
	Category string
	Sort     string
	ViewerID int64
	Offset   int
	Limit    int
}

// List 按范围分页查询圈子。
func (r *CircleRepository) List(f CircleFilter) ([]model.Circle, int64, error) {
	q := r.db.Model(&model.Circle{}).Where("status = ?", model.CircleMemberNormal)
	switch f.Scope {
	case "recommended":
		q = q.Where("is_recommended = ?", true)
	case "joined":
		if f.ViewerID <= 0 {
			return []model.Circle{}, 0, nil
		}
		q = q.Where("id IN (?)", r.db.Model(&model.CircleMember{}).
			Select("circle_id").Where("user_id = ? AND status <> ?", f.ViewerID, model.CircleMemberRemoved))
	case "created":
		if f.ViewerID <= 0 {
			return []model.Circle{}, 0, nil
		}
		q = q.Where("owner_id = ?", f.ViewerID)
	}
	if f.Keyword != "" {
		q = q.Where("name LIKE ? OR description LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	order := "created_at DESC"
	if f.Sort == "hot" {
		order = "member_count DESC, post_count DESC"
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Circle
	if err := q.Order(order).Offset(f.Offset).Limit(f.Limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListAll 分页查询全部圈子（后台）。
func (r *CircleRepository) ListAll(offset, limit int) ([]model.Circle, int64) {
	var total int64
	r.db.Model(&model.Circle{}).Count(&total)
	var rows []model.Circle
	r.db.Order("id ASC").Offset(offset).Limit(limit).Find(&rows)
	return rows, total
}

// FindByID 按 ID 查询圈子，不存在返回 (nil, nil)。
func (r *CircleRepository) FindByID(id int64) (*model.Circle, error) {
	var c model.Circle
	err := r.db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// FindByIDs 批量查询圈子，返回 id->圈子 映射。
func (r *CircleRepository) FindByIDs(ids []int64) map[int64]*model.Circle {
	res := map[int64]*model.Circle{}
	if len(ids) == 0 {
		return res
	}
	var cs []model.Circle
	r.db.Where("id IN ?", ids).Find(&cs)
	for i := range cs {
		res[cs[i].ID] = &cs[i]
	}
	return res
}

// CountByName 统计同名圈子数量。
func (r *CircleRepository) CountByName(name string) int64 {
	var count int64
	r.db.Model(&model.Circle{}).Where("name = ?", name).Count(&count)
	return count
}

// CreateWithOwner 在事务中创建圈子并写入圈主成员记录。
func (r *CircleRepository) CreateWithOwner(c *model.Circle, owner *model.CircleMember) error {
	return tx(r.db, func(t *gorm.DB) error {
		if err := t.Create(c).Error; err != nil {
			return err
		}
		owner.CircleID = c.ID
		return t.Create(owner).Error
	})
}

// IncMemberCount 对成员数做增量（可为负）。
func (r *CircleRepository) IncMemberCount(id int64, delta int) {
	if delta >= 0 {
		r.db.Model(&model.Circle{}).Where("id = ?", id).UpdateColumn("member_count", gorm.Expr("member_count + ?", delta))
	} else {
		r.db.Model(&model.Circle{}).Where("id = ?", id).UpdateColumn("member_count", gorm.Expr("GREATEST(member_count - ?, 0)", -delta))
	}
}

// FindMember 查询用户在圈子的成员记录，不存在返回 (nil, nil)。
func (r *CircleRepository) FindMember(circleID, userID int64) (*model.CircleMember, error) {
	var m model.CircleMember
	err := r.db.Where("circle_id = ? AND user_id = ?", circleID, userID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// FindMembersByViewer 批量查询查看者在多个圈子的成员记录，返回 circleID->成员。
func (r *CircleRepository) FindMembersByViewer(viewerID int64, circleIDs []int64) map[int64]*model.CircleMember {
	res := map[int64]*model.CircleMember{}
	if viewerID <= 0 || len(circleIDs) == 0 {
		return res
	}
	var ms []model.CircleMember
	r.db.Where("user_id = ? AND circle_id IN ?", viewerID, circleIDs).Find(&ms)
	for i := range ms {
		res[ms[i].CircleID] = &ms[i]
	}
	return res
}

// CreateMember 新增成员。
func (r *CircleRepository) CreateMember(m *model.CircleMember) error {
	return r.db.Create(m).Error
}

// UpdateMember 更新成员记录字段。
func (r *CircleRepository) UpdateMember(circleID, userID int64, updates map[string]any) error {
	return r.db.Model(&model.CircleMember{}).
		Where("circle_id = ? AND user_id = ?", circleID, userID).Updates(updates).Error
}

// DeleteMember 物理删除成员记录（退出圈子）。
func (r *CircleRepository) DeleteMember(m *model.CircleMember) {
	r.db.Delete(m)
}

// Members 分页查询圈子成员（排除已移除）。
func (r *CircleRepository) Members(circleID int64, offset, limit int) ([]model.CircleMember, int64) {
	q := r.db.Model(&model.CircleMember{}).Where("circle_id = ? AND status <> ?", circleID, model.CircleMemberRemoved)
	var total int64
	q.Count(&total)
	var rows []model.CircleMember
	q.Order("FIELD(role,'owner','moderator','reviewer','member'), joined_at ASC").
		Offset(offset).Limit(limit).Find(&rows)
	return rows, total
}
