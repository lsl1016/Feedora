package api

import (
	"github.com/feedora/backend/internal/dto"
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
// @Param    req  query  dto.HotRankQuery  false  "热门榜单查询参数"
// @Success  200  {object}  dto.HotRankListResponse
// @Failure  400  {object}  response.Body
// @Router   /hot/ranks [get]
func (h *RankAPI) HotRanks(c *gin.Context) {
	var req dto.HotRankQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list := h.svc.HotRanks(req.RankType, req.TimeRange, req.Page, req.PageSize)
	response.OK(c, list)
}
