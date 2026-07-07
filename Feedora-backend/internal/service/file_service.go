package service

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/ossx"
	"github.com/feedora/backend/pkg/utils"
)

const maxUploadSize = 5 << 20 // 5MB

// FileService 文件上传业务逻辑。
type FileService struct {
	files   *repository.FileRepository
	storage ossx.Storage
}

func NewFileService(files *repository.FileRepository, storage ossx.Storage) *FileService {
	return &FileService{files: files, storage: storage}
}

// UploadInput 上传参数。
type UploadInput struct {
	UserID   int64
	BizType  string
	Filename string
	MimeType string
	Size     int64
	Reader   io.Reader
}

// Upload 校验并上传图片文件，写入文件记录。
func (s *FileService) Upload(in UploadInput) (*dto.UploadResult, error) {
	if in.Size > maxUploadSize {
		return nil, errs.New(422, "文件大小不能超过 5MB")
	}
	ext, ok := utils.ImageExtByMime(in.MimeType)
	if !ok {
		ext = strings.ToLower(filepath.Ext(in.Filename))
		if !utils.IsImageExt(ext) {
			return nil, errs.New(422, "仅支持上传图片文件")
		}
	}
	if in.BizType == "" {
		in.BizType = "post_image"
	}
	objectKey := fmt.Sprintf("%s/%s/%s%s", bizDir(in.BizType), time.Now().Format("20060102"), utils.UUID(), ext)
	url, err := s.storage.Upload(objectKey, in.Reader, in.MimeType)
	if err != nil {
		return nil, errs.ErrUploadFailed
	}
	s.files.Create(&model.File{
		UserID: in.UserID, BizType: in.BizType, Filename: in.Filename,
		ObjectKey: objectKey, URL: url, Size: in.Size, MimeType: in.MimeType,
		CreatedAt: time.Now(),
	})
	return &dto.UploadResult{URL: url, Filename: in.Filename, Size: in.Size, MimeType: in.MimeType}, nil
}

// bizDir 根据业务类型返回顶层目录。
func bizDir(bizType string) string {
	switch bizType {
	case "avatar":
		return "avatar"
	case "circle_avatar":
		return "circle"
	case "topic_cover":
		return "topic"
	case "activity_cover":
		return "activity"
	default:
		return "post"
	}
}
