// Package errs 定义统一的业务错误结构与错误码。
// 包名使用 errs 以避免与标准库 errors 冲突。
package errs

// Error 表示携带稳定错误码与提示信息的业务错误。
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string { return e.Message }

// New 创建业务错误。
func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}
