package errs

// 与 HTTP 状态码对齐的通用错误码。
var (
	Success       = New(0, "success")
	ErrParams     = New(400, "请求参数错误")
	ErrUnauth     = New(401, "未登录或 token 过期")
	ErrForbidden  = New(403, "无权限")
	ErrNotFound   = New(404, "资源不存在")
	ErrConflict   = New(409, "资源冲突")
	ErrValidation = New(422, "业务校验失败")
	ErrRateLimit  = New(429, "请求过于频繁")
	ErrInternal   = New(500, "服务内部错误")
)

// 业务错误码。
var (
	ErrAccountOrPassword = New(10001, "账号或密码错误")
	ErrAccountExists     = New(10002, "账号已存在")
	ErrUserBanned        = New(10003, "用户已被封禁")
	ErrPostNotFound      = New(20001, "帖子不存在")
	ErrPostInvisible     = New(20002, "帖子不可见")
	ErrPostNoPermission  = New(20003, "无权操作该帖子")
	ErrCommentNotFound   = New(30001, "评论不存在")
	ErrCircleNotFound    = New(40001, "圈子不存在")
	ErrCircleMuted       = New(40002, "你已被该圈子禁言")
	ErrCircleNotJoined   = New(40003, "未加入圈子")
	ErrUploadFailed      = New(50001, "文件上传失败")
	ErrTaskNotDone       = New(60001, "任务进度未达标，暂不能领取")
	ErrTaskClaimed       = New(60002, "任务奖励已领取过")
	ErrTaskClaimViaCheckIn = New(60003, "该任务通过每日签到完成，无需领取")
)
