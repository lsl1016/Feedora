package utils

import "github.com/google/uuid"

// UUID 生成一个随机 UUID 字符串。
func UUID() string {
	return uuid.NewString()
}
