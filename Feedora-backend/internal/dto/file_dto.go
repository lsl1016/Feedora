package dto

// UploadFileForm 文件上传请求。
type UploadFileForm struct {
	BizType string `form:"bizType" example:"post_image"` // 业务类型。
}

// UploadResult 文件上传响应。
type UploadResult struct {
	URL      string `json:"url"`      // 文件访问地址。
	Filename string `json:"filename"` // 原始文件名。
	Size     int64  `json:"size"`     // 文件大小。
	MimeType string `json:"mimeType"` // 文件 MIME 类型。
}

// UploadResponse 文件上传成功响应。
type UploadResponse struct {
	TraceEnvelope
	Data UploadResult `json:"data"` // 业务数据。
}
