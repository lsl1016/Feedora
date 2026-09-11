package service

import (
	"strconv"
	"time"
)

// itoa 将 int64 转为字符串。
func itoa(n int64) string { return strconv.FormatInt(n, 10) }

// parseInt 将字符串转为 int64，失败返回 0。
func parseInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

// normPage 规整分页参数，带默认值与上限保护。
func normPage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

// offset 计算 SQL OFFSET。
func offset(page, size int) int {
	return (page - 1) * size
}

// idemWindow 写操作防重复提交的判定时间窗（同作者同内容）。
const idemWindow = 10 * time.Second
