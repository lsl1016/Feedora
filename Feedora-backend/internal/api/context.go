package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// pageParams 从 query 解析分页参数，带默认值与上限保护。
func pageParams(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.Query("page"))
	size, _ = strconv.Atoi(c.Query("pageSize"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return
}

// paramID 从路径参数解析 int64 ID。
func paramID(c *gin.Context, name string) int64 {
	id, _ := strconv.ParseInt(c.Param(name), 10, 64)
	return id
}

// queryID 从 query 解析 int64 ID。
func queryID(c *gin.Context, name string) int64 {
	id, _ := strconv.ParseInt(c.Query(name), 10, 64)
	return id
}
