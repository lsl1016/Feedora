package ossx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Storage 是可插拔的文件存储抽象。阶段一提供本地磁盘实现，
// 后续 MinIO / 阿里云 OSS 可实现同一接口。
type Storage interface {
	Upload(objectKey string, reader io.Reader, contentType string) (url string, err error)
	Delete(objectKey string) error
	PublicURL(objectKey string) string
}

// LocalStorage 本地磁盘存储实现（模拟 OSS）。
type LocalStorage struct {
	basePath      string
	publicBaseURL string
}

func NewLocalStorage(basePath, publicBaseURL string) (*LocalStorage, error) {
	if basePath == "" {
		basePath = "./uploads"
	}
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		basePath:      basePath,
		publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
	}, nil
}

// BasePath 返回本地根目录，用于对外提供静态文件访问。
func (s *LocalStorage) BasePath() string { return s.basePath }

func (s *LocalStorage) Upload(objectKey string, reader io.Reader, _ string) (string, error) {
	objectKey = strings.TrimLeft(objectKey, "/")
	full := filepath.Join(s.basePath, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", err
	}
	f, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}
	return s.PublicURL(objectKey), nil
}

func (s *LocalStorage) Delete(objectKey string) error {
	objectKey = strings.TrimLeft(objectKey, "/")
	return os.Remove(filepath.Join(s.basePath, filepath.FromSlash(objectKey)))
}

func (s *LocalStorage) PublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s", s.publicBaseURL, strings.TrimLeft(objectKey, "/"))
}
