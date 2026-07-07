package dto

// LoginRequest 登录请求。
type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Account         string `json:"account"`
	Nickname        string `json:"nickname"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// LoginResult 登录响应。
type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
