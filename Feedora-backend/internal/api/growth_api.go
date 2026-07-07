package api

import (
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
// @Success  200  {object}  response.Body
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
// @Param    type  query  string  false  "任务类型"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /growth/tasks [get]
func (h *GrowthAPI) Tasks(c *gin.Context) {
	response.OK(c, h.svc.Tasks(middleware.CurrentUserID(c), c.Query("type")))
}

// ClaimTask 领取任务奖励
// @Summary  领取任务奖励
// @Tags     成长
// @Produce  json
// @Param    taskId  path  int  true  "任务ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /growth/tasks/{taskId}/claim [post]
func (h *GrowthAPI) ClaimTask(c *gin.Context) {
	response.OK(c, gin.H{"claimed": true})
}

// Rankings 成长排行榜
// @Summary  成长排行榜
// @Tags     成长
// @Produce  json
// @Param    type      query  string  false  "排行类型"
// @Param    range     query  string  false  "时间范围"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /growth/rankings [get]
func (h *GrowthAPI) Rankings(c *gin.Context) {
	page, size := pageParams(c)
	list := h.svc.Rankings(c.Query("type"), c.Query("range"), page, size, middleware.CurrentUserID(c))
	response.OK(c, list)
}
