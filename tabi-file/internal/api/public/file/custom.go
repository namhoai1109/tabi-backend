package file

import "tabi-file/internal/model"

// ListFileResponse contains list of file
// swagger:model ListFileResponse
type ListFileResponse struct {
	Data []*model.FileResponse `json:"data"`
}

// UploadPresignedRequest model
// swagger:model UploadPresignedRequest
type UploadPresignedRequest struct {
	Data []*FilePreSignedRequest `json:"data" validate:"required"`
}

// FilePreSignedRequest model
// swagger:model FilePreSignedRequest
type FilePreSignedRequest struct {
	FileName       string `json:"file_name"`
	AttachmentID   int    `json:"attachment_id"`
	AttachmentType string `json:"attachment_type"`
	ContentType    string `json:"content_type"`
	FileSize       int    `json:"file_size"`
	Field          string `json:"field"`
	SortOrder      int    `json:"sort_order,omitempty"`
}
