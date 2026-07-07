package utils

import "strings"

// 允许上传的图片 MIME 类型及对应扩展名。
var imageMimeExt = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// ImageExtByMime 根据 MIME 返回扩展名，第二个返回值表示是否为允许的图片类型。
func ImageExtByMime(mime string) (string, bool) {
	ext, ok := imageMimeExt[mime]
	return ext, ok
}

// IsImageExt 判断扩展名是否为图片。
func IsImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return true
	}
	return false
}
