package service

import (
	"strings"
	"time"

	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/jwtx"
	"github.com/feedora/backend/pkg/utils"
)

// AuthService 认证业务逻辑。
type AuthService struct {
	users    *repository.UserRepository
	jwt      *jwtx.Manager
	producer event.Producer
}

func NewAuthService(users *repository.UserRepository, jm *jwtx.Manager, producer event.Producer) *AuthService {
	return &AuthService{users: users, jwt: jm, producer: producer}
}

// Register 注册新用户。
func (s *AuthService) Register(in dto.RegisterRequest) (*dto.LoginResult, error) {
	in.Account = strings.TrimSpace(in.Account)
	in.Nickname = strings.TrimSpace(in.Nickname)
	if in.Account == "" || in.Password == "" {
		return nil, errs.ErrParams
	}
	if in.ConfirmPassword != "" && in.ConfirmPassword != in.Password {
		return nil, errs.New(422, "两次输入的密码不一致")
	}
	if in.Nickname == "" {
		in.Nickname = in.Account
	}
	count, err := s.users.CountByAccount(in.Account)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if count > 0 {
		return nil, errs.ErrAccountExists
	}
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, errs.ErrInternal
	}
	now := time.Now()
	u := &model.User{
		Account:      in.Account,
		PasswordHash: hash,
		Nickname:     in.Nickname,
		Role:         model.RoleUser,
		Status:       model.UserNormal,
		Level:        1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(u); err != nil {
		return nil, errs.ErrInternal
	}
	s.producer.Publish(event.TopicUser, event.UserRegistered, u.ID, u.ID, map[string]any{"account": u.Account})
	return s.issue(u)
}

// Login 账号密码登录。
func (s *AuthService) Login(account, password string) (*dto.LoginResult, error) {
	account = strings.TrimSpace(account)
	if account == "" || password == "" {
		return nil, errs.ErrParams
	}
	u, err := s.users.FindByAccount(account)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if u == nil || !utils.CheckPassword(u.PasswordHash, password) {
		return nil, errs.ErrAccountOrPassword
	}
	if u.Status == model.UserBanned {
		return nil, errs.ErrUserBanned
	}
	return s.issue(u)
}

// Me 获取当前登录用户信息。
func (s *AuthService) Me(userID int64) (*dto.User, error) {
	u, err := s.users.FindByID(userID)
	if err != nil || u == nil {
		return nil, errs.ErrUnauth
	}
	res := dto.ToUser(u, false, 0)
	return &res, nil
}

// issue 签发 token 并组装登录结果。
func (s *AuthService) issue(u *model.User) (*dto.LoginResult, error) {
	token, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return &dto.LoginResult{Token: token, User: dto.ToUser(u, false, 0)}, nil
}
