package repository

import (
	"time"

	"github.com/feedora/backend/internal/model"
	"gorm.io/gorm"
)

// GrowthRepository 成长积分数据访问。
type GrowthRepository struct {
	db *gorm.DB
}

func NewGrowthRepository(db *gorm.DB) *GrowthRepository {
	return &GrowthRepository{db: db}
}

// AddPointLog 写入积分流水并累加用户积分。依赖唯一索引
// uk_user_action_biz 防止重复发放；返回 true 表示本次实际发放。
func (r *GrowthRepository) AddPointLog(userID int64, action string, point int, bizType string, bizID int64, remark string) bool {
	log := &model.UserPointLog{
		UserID: userID, Action: action, Point: point,
		BizType: bizType, BizID: bizID, Remark: remark, CreatedAt: time.Now(),
	}
	res := r.db.Create(log)
	if res.Error != nil || res.RowsAffected == 0 {
		return false // 已发放过（唯一索引冲突）
	}
	r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("point_count", gorm.Expr("point_count + ?", point))
	return true
}

// HasPointLog 判断某业务是否已发放过积分。
func (r *GrowthRepository) HasPointLog(userID int64, action, bizType string, bizID int64) bool {
	var n int64
	r.db.Model(&model.UserPointLog{}).
		Where("user_id = ? AND action = ? AND biz_type = ? AND biz_id = ?", userID, action, bizType, bizID).
		Count(&n)
	return n > 0
}

// CountRecentCheckIns 统计用户最近 n 天内的签到次数（用于连续签到近似）。
func (r *GrowthRepository) CountRecentCheckIns(userID int64, days int) int64 {
	since := time.Now().AddDate(0, 0, -days)
	var n int64
	r.db.Model(&model.UserPointLog{}).
		Where("user_id = ? AND action = ? AND created_at >= ?", userID, "check_in", since).
		Count(&n)
	return n
}
