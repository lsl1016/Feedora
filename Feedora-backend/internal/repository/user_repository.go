package repository

import (
	"errors"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问。
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID 按 ID 查询用户，不存在返回 (nil, nil)。
func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByAccount 按账号查询用户，不存在返回 (nil, nil)。
func (r *UserRepository) FindByAccount(account string) (*model.User, error) {
	var u model.User
	err := r.db.Where("account = ?", account).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByIDs 批量查询用户，返回 id->用户 映射。
func (r *UserRepository) FindByIDs(ids []int64) (map[int64]*model.User, error) {
	res := map[int64]*model.User{}
	if len(ids) == 0 {
		return res, nil
	}
	var us []model.User
	if err := r.db.Where("id IN ?", ids).Find(&us).Error; err != nil {
		return nil, err
	}
	for i := range us {
		res[us[i].ID] = &us[i]
	}
	return res, nil
}

// Create 创建用户。
func (r *UserRepository) Create(u *model.User) error {
	return r.db.Create(u).Error
}

// CountByAccount 统计账号数量，用于判断是否已存在。
func (r *UserRepository) CountByAccount(account string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("account = ?", account).Count(&count).Error
	return count, err
}

// Update 按 ID 更新字段。
func (r *UserRepository) Update(id int64, updates map[string]any) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// List 按关键词分页查询用户。
func (r *UserRepository) List(keyword string, offset, limit int) ([]model.User, int64, error) {
	q := r.db.Model(&model.User{}).Where("status <> ?", model.UserBanned)
	if keyword != "" {
		q = q.Where("nickname LIKE ? OR account LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.User
	if err := q.Order("id ASC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// IncColumn 对指定计数列做增量（可为负），并做非负保护。
func (r *UserRepository) IncColumn(id int64, column string, delta int) {
	if delta >= 0 {
		r.db.Model(&model.User{}).Where("id = ?", id).
			UpdateColumn(column, gorm.Expr(column+" + ?", delta))
	} else {
		r.db.Model(&model.User{}).Where("id = ?", id).
			UpdateColumn(column, gorm.Expr("GREATEST("+column+" - ?, 0)", -delta))
	}
}
