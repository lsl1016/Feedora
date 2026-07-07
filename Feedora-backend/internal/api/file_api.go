package api

import (
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// FileAPI 文件上传接口。
type FileAPI struct {
	svc *service.FileService
}

func NewFileAPI(svc *service.FileService) *FileAPI {
	return &FileAPI{svc: svc}
}

// Upload 上传文件
// @Summary  上传文件
// @Tags     文件
// @Accept   multipart/form-data
// @Produce  json
// @Param    file      formData  file    true   "上传的文件"
// @Param    bizType   formData  string  false  "业务类型"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /files/upload [post]
func (h *FileAPI) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	src, err := fileHeader.Open()
	if err != nil {
		response.Fail(c, errs.ErrUploadFailed)
		return
	}
	defer src.Close()

	res, err := h.svc.Upload(service.UploadInput{
		UserID:   middleware.CurrentUserID(c),
		BizType:  c.DefaultPostForm("bizType", "post_image"),
		Filename: fileHeader.Filename,
		MimeType: fileHeader.Header.Get("Content-Type"),
		Size:     fileHeader.Size,
		Reader:   src,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
