package dto

// LoginRequest 登录请求。
type LoginRequest struct {
	Account  string `json:"account" binding:"required" example:"feedora"` // 登录账号。
	Password string `json:"password" binding:"required" example:"123456"` // 登录密码。
}

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Account         string `json:"account" binding:"required" example:"feedora"` // 注册账号。
	Nickname        string `json:"nickname" example:"Feedora 用户"`                // 用户昵称。
	Password        string `json:"password" binding:"required" example:"123456"` // 登录密码。
	ConfirmPassword string `json:"confirmPassword" example:"123456"`             // 确认密码。
}

// LoginResult 登录响应数据。
type LoginResult struct {
	Token string `json:"token"` // 访问令牌。
	User  User   `json:"user"`  // 当前登录用户。
}

// LoginResponse 登录成功响应。
type LoginResponse struct {
	TraceEnvelope
	Data LoginResult `json:"data"` // 业务数据。
}

// CurrentUserResponse 当前用户信息响应。
type CurrentUserResponse struct {
	TraceEnvelope
	Data User `json:"data"` // 业务数据。
}
