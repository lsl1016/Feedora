package response

import (
	stderrors "errors"
	"net/http"
	"time"

	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Body struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	TraceID   string      `json:"traceId"`
	Timestamp string      `json:"timestamp"`
}

// PageData 是前端约定的统一分页数据结构。
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func traceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

// OK 写出成功响应。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:      0,
		Message:   "success",
		Data:      data,
		TraceID:   traceID(c),
		Timestamp: now(),
	})
}

// Page 写出分页成功响应。
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	OK(c, PageData{List: list, Total: total, Page: page, PageSize: pageSize})
}

// Fail 写出错误响应。通用错误码映射到对应 HTTP 状态码，业务错误码返回 HTTP 200。
// 所有错误在此统一落盘，便于凭 traceId 快速排查：5xx 记 ERROR，其余记 WARN。
func Fail(c *gin.Context, err error) {
	var be *errs.Error
	if !stderrors.As(err, &be) {
		be = &errs.Error{Code: errs.ErrInternal.Code, Message: err.Error()}
	}
	status := httpStatus(be.Code)
	// 按严重程度分级记录：服务端错误 ERROR，客户端/业务错误 WARN。
	logf := logger.Warnf
	if status >= http.StatusInternalServerError {
		logf = logger.Errorf
	}
	logf("request failed, traceId:%s method:%s path:%s status:%d code:%d message:%s userId:%d",
		traceID(c), c.Request.Method, c.Request.URL.Path, status, be.Code, be.Message, currentUserID(c))
	c.JSON(status, Body{
		Code:      be.Code,
		Message:   be.Message,
		Data:      nil,
		TraceID:   traceID(c),
		Timestamp: now(),
	})
}

// currentUserID 从上下文读取登录用户 ID，未登录返回 0。
// 就地读取避免引入 middleware 包造成循环依赖。
func currentUserID(c *gin.Context) int64 {
	if v, ok := c.Get("userId"); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

func httpStatus(code int) int {
	switch code {
	case 400, 401, 403, 404, 409, 422, 429, 500:
		return code
	default:
		return http.StatusOK
	}
}
