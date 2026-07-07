package api

import (
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// RankAPI 热门榜单接口。
type RankAPI struct {
	svc *service.RankService
}

func NewRankAPI(svc *service.RankService) *RankAPI {
	return &RankAPI{svc: svc}
}

// HotRanks 热门榜单
// @Summary  热门榜单
// @Tags     榜单
// @Produce  json
// @Param    rankType   query  string  false  "榜单类型"
// @Param    timeRange  query  string  false  "时间范围"
// @Param    page       query  int     false  "页码"
// @Param    pageSize   query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /hot/ranks [get]
func (h *RankAPI) HotRanks(c *gin.Context) {
	page, size := pageParams(c)
	list := h.svc.HotRanks(c.Query("rankType"), c.Query("timeRange"), page, size)
	response.OK(c, list)
}
