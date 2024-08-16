package model

import "time"

// File represents the file model
// swagger:model
type File struct {
	ID             int    `json:"id" gorm:"primary_key"`
	PathName       string `json:"path_name" gorm:"type:varchar(255)"`
	FileName       string `json:"file_name" gorm:"type:varchar(255)"`
	FileSize       int64  `json:"file_size"`
	ContentType    string `json:"content_type" gorm:"content_type:varchar(255)"`
	Field          string `json:"field" gorm:"type:text"`
	AttachmentID   int    `json:"attachment_id"`
	AttachmentType string `json:"attachment_type" gorm:"type:varchar(255)"`
	SortOrder      int    `json:"sort_order"`
	FullURL        string `json:"full_url" gorm:"type:text"`
	UploadURL      string `json:"upload_url" gorm:"type:text"`

	Base
}

func (f *File) ToFileResponse() *FileResponse {
	return &FileResponse{
		ID:             f.ID,
		CreatedAt:      f.CreatedAt,
		PathName:       f.PathName,
		FileName:       f.FileName,
		Field:          f.Field,
		AttachmentID:   f.AttachmentID,
		AttachmentType: f.AttachmentType,
		GetURL:         f.FullURL,
		UploadURL:      f.UploadURL,
		SortOrder:      f.SortOrder,
	}
}

// FileResponse model
// swagger:model FileResponse
type FileResponse struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	PathName       string    `json:"path_name,"`
	FileName       string    `json:"file_name"`
	Field          string    `json:"field"`
	AttachmentID   int       `json:"attachment_id"`
	AttachmentType string    `json:"attachment_type"`
	GetURL         string    `json:"get_url"`
	UploadURL      string    `json:"upload_url,omitempty"`
	SortOrder      int       `json:"sort_order"`
}
