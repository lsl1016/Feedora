package logger

import (
	"log"
	"os"
	"strings"
)

// 轻量日志封装。阶段一保持依赖精简，后续可无侵入替换为 zap。
// 支持通过环境变量 LOG_LEVEL 控制输出级别：debug < info < warn < error。
var std = log.New(os.Stdout, "", log.LstdFlags)

// 日志级别。数值越小越详细。
const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
)

// 当前生效的最低输出级别，默认 info。
var minLevel = parseLevel(os.Getenv("LOG_LEVEL"))

// parseLevel 解析级别字符串，非法或空值回退到 info。
func parseLevel(s string) int {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return levelDebug
	case "warn", "warning":
		return levelWarn
	case "error":
		return levelError
	default:
		return levelInfo
	}
}

// SetLevel 运行时调整最低输出级别（如 "debug"）。
func SetLevel(s string) { minLevel = parseLevel(s) }

func Debugf(format string, args ...interface{}) {
	if minLevel <= levelDebug {
		std.Printf("[DEBUG] "+format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if minLevel <= levelInfo {
		std.Printf("[INFO] "+format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if minLevel <= levelWarn {
		std.Printf("[WARN] "+format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if minLevel <= levelError {
		std.Printf("[ERROR] "+format, args...)
	}
}
