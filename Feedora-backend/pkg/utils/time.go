package utils

import "time"

// FormatTime 统一时间格式为 RFC3339，前端使用 dayjs 解析。
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// FormatTimePtr 格式化时间指针，nil 或零值返回空串。
func FormatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
