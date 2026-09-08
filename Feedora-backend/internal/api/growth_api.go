package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GrowthAPI 成长积分接口。
type GrowthAPI struct {
	svc *service.GrowthService
}

func NewGrowthAPI(svc *service.GrowthService) *GrowthAPI {
	return &GrowthAPI{svc: svc}
}

// CheckIn 每日签到
// @Summary  每日签到
// @Tags     成长
// @Produce  json
// @Success  200  {object}  dto.CheckInResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /growth/check-in [post]
func (h *GrowthAPI) CheckIn(c *gin.Context) {
	res, err := h.svc.CheckIn(middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Tasks 成长任务列表
// @Summary  成长任务列表
// @Tags     成长
// @Produce  json
// @Param    req  query  dto.GrowthTaskQuery  false  "成长任务查询参数"
// @Success  200  {object}  dto.TaskListResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /growth/tasks [get]
func (h *GrowthAPI) Tasks(c *gin.Context) {
	var req dto.GrowthTaskQuery
	if !bindQuery(c, &req) {
		return
	}
	response.OK(c, h.svc.Tasks(middleware.CurrentUserID(c), req.Type))
}

// ClaimTask 领取任务奖励
// @Summary  领取任务奖励
// @Tags     成长
// @Produce  json
// @Param    req  path  dto.TaskIDURI  true  "任务路径参数"
// @Success  200  {object}  dto.ClaimedResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /growth/tasks/{taskId}/claim [post]
func (h *GrowthAPI) ClaimTask(c *gin.Context) {
	var req dto.TaskIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.ClaimTask(middleware.CurrentUserID(c), req.TaskID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"claimed": true})
}

// Rankings 成长排行榜
// @Summary  成长排行榜
// @Tags     成长
// @Produce  json
// @Param    req  query  dto.GrowthRankingQuery  false  "成长排行榜查询参数"
// @Success  200  {object}  dto.RankingListResponse
// @Failure  400  {object}  response.Body
// @Router   /growth/rankings [get]
func (h *GrowthAPI) Rankings(c *gin.Context) {
	var req dto.GrowthRankingQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list := h.svc.Rankings(req.Type, req.Range, req.Page, req.PageSize, middleware.CurrentUserID(c))
	response.OK(c, list)
}
