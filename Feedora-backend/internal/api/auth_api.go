package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthAPI 认证接口。
type AuthAPI struct {
	svc *service.AuthService
}

func NewAuthAPI(svc *service.AuthService) *AuthAPI {
	return &AuthAPI{svc: svc}
}

// Login 用户登录
// @Summary  用户登录
// @Tags     认证
// @Accept   json
// @Produce  json
// @Param    body  body  dto.LoginRequest  true  "登录请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /auth/login [post]
func (h *AuthAPI) Login(c *gin.Context) {
	var in dto.LoginRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.Login(in.Account, in.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Register 用户注册
// @Summary  用户注册
// @Tags     认证
// @Accept   json
// @Produce  json
// @Param    body  body  dto.RegisterRequest  true  "注册请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /auth/register [post]
func (h *AuthAPI) Register(c *gin.Context) {
	var in dto.RegisterRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.Register(in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Me 获取当前登录用户信息
// @Summary  获取当前登录用户信息
// @Tags     认证
// @Produce  json
// @Security BearerAuth
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /auth/me [get]
func (h *AuthAPI) Me(c *gin.Context) {
	res, err := h.svc.Me(middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Logout 退出登录（无状态 JWT，登出由前端清除 token）
// @Summary  退出登录
// @Tags     认证
// @Produce  json
// @Security BearerAuth
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /auth/logout [post]
func (h *AuthAPI) Logout(c *gin.Context) {
	response.OK(c, gin.H{})
}
