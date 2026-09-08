package service

import (
	"time"

	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
)

const checkInPoints = 5

// GrowthService 成长积分业务逻辑。
type GrowthService struct {
	growth *repository.GrowthRepository
	users  *repository.UserRepository
	rank   *RankService
}

func NewGrowthService(growth *repository.GrowthRepository, users *repository.UserRepository, rank *RankService) *GrowthService {
	return &GrowthService{growth: growth, users: users, rank: rank}
}

// CheckIn 每日签到，同日重复签到不重复发放积分。
func (s *GrowthService) CheckIn(userID int64) (*dto.CheckInResult, error) {
	today := int64(time.Now().Year()*10000 + int(time.Now().Month())*100 + time.Now().Day())
	awarded := 0
	if s.growth.AddPointLog(userID, "check_in", checkInPoints, "checkin", today, "每日签到") {
		awarded = checkInPoints
	}
	u, err := s.users.FindByID(userID)
	if err != nil || u == nil {
		return nil, errs.ErrInternal
	}
	continuous := int(s.growth.CountRecentCheckIns(userID, 7))
	return &dto.CheckInResult{Points: u.PointCount, ContinuousDays: continuous, AwardedPoints: awarded}, nil
}

// Tasks 返回成长任务列表（阶段二用静态任务 + 用户当前进度近似）。
func (s *GrowthService) Tasks(userID int64, typ string) []dto.Task {
	u, _ := s.users.FindByID(userID)
	posted := int64(0)
	commented := int64(0)
	if u != nil {
		posted = u.PostCount
		commented = u.CommentCount
	}
	tasks := []dto.Task{
		{TaskID: 1, Title: "发布首篇帖子", Description: "发布你的第一篇帖子", Type: "newbie", RewardPoints: 10, TargetValue: 1, CurrentValue: min64(posted, 1), Status: doneIf(posted >= 1), ActionText: "去发帖", ActionURL: "/posts/new"},
		{TaskID: 2, Title: "参与评论", Description: "发表 3 条评论", Type: "daily", RewardPoints: 3, TargetValue: 3, CurrentValue: min64(commented, 3), Status: doneIf(commented >= 3), ActionText: "去评论", ActionURL: "/"},
		{TaskID: 3, Title: "每日签到", Description: "坚持每日签到", Type: "daily", RewardPoints: checkInPoints, TargetValue: 1, CurrentValue: 0, Status: "todo", ActionText: "去签到", ActionURL: "/growth"},
	}
	if typ != "" && typ != "all" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if t.Type == typ {
				filtered = append(filtered, t)
			}
		}
		return filtered
	}
	return tasks
}

// Rankings 排行榜，委托 RankService。
func (s *GrowthService) Rankings(rankType, timeRange string, page, size int, currentUserID int64) []dto.RankingItem {
	return s.rank.Rankings(rankType, timeRange, page, size, currentUserID)
}

// ClaimTask 领取任务奖励：进度达标后发放积分。
// 唯一索引 uk_user_action_biz(action=task_claim, biz_type=task) 保证同一任务只发一次；
// 签到任务（ID=3）不走领取，由签到接口直接发放。
func (s *GrowthService) ClaimTask(userID, taskID int64) error {
	u, err := s.users.FindByID(userID)
	if err != nil || u == nil {
		return errs.ErrInternal
	}
	var progress, target int64
	var points int
	var title string
	switch taskID {
	case 1:
		progress, target, points, title = u.PostCount, 1, 10, "发布首篇帖子"
	case 2:
		progress, target, points, title = u.CommentCount, 3, 3, "参与评论"
	case 3:
		return errs.ErrTaskClaimViaCheckIn
	default:
		return errs.ErrParams
	}
	if progress < target {
		return errs.ErrTaskNotDone
	}
	if !s.growth.AddPointLog(userID, "task_claim", points, "task", taskID, "任务奖励："+title) {
		return errs.ErrTaskClaimed
	}
	return nil
}

func min64(v, max int64) int {
	if v > max {
		return int(max)
	}
	return int(v)
}

func doneIf(cond bool) string {
	if cond {
		return "done"
	}
	return "todo"
}
