package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// QueryHash 计算查询字符串的哈希，用于阶段二缓存 Key（如 feed:home:{queryHash}）。
func QueryHash(raw string) string {
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}
