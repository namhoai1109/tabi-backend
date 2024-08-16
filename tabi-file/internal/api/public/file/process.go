package file

import (
	"context"
	"fmt"
	"tabi-file/internal/model"
	"tabi-file/internal/util/filehelper"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

func (s *File) validateFileName(filename string, index int) (string, error) {

	if filename == "" {
		return "", server.NewHTTPValidationError(fmt.Sprintf("File name is required for file %d", index+1))
	}

	newFilename := filename
	if len(filename) > 255 {
		newFilename = filename[:255]
	}

	// Validate the file extension
	if filehelper.IsOverType(filename) {
		return "", server.NewHTTPValidationError(fmt.Sprintf("%v invalid file extension", filename))
	}

	// Validate the file name
	if filehelper.AvoidCharacters(filename) {
		return "", server.NewHTTPValidationError(fmt.Sprintf("%v contains invalid characters", filename))
	}

	return newFilename, nil
}

// MappingFileURL to get url by path
func (s *File) mappingFileURL(key string) string {
	return fmt.Sprintf("https://%s.s3.ap-southeast-1.amazonaws.com%s", s.cfg.S3PublicBucketName, key)
}

func (s *File) createFile(db *gorm.DB, ctx context.Context, file *model.File, index int) error {
	condition := map[string]interface{}{
		"attachment_id":   file.AttachmentID,
		"attachment_type": file.AttachmentType,
		"field":           file.Field,
		"sort_order":      file.SortOrder,
	}

	exist, err := s.fileDB.Exist(db, condition)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("cannot check exist file %d with error: %s", index+1, err))
		return server.NewHTTPInternalError(fmt.Sprintf("cannot check exist file %d with error: %s", index+1, err))
	}

	if exist {
		if err := s.fileDB.Delete(db, map[string]interface{}{
			"attachment_id":   file.AttachmentID,
			"attachment_type": file.AttachmentType,
			"field":           file.Field,
			"sort_order":      file.SortOrder,
		}); err != nil {
			logger.LogError(ctx, fmt.Sprintf("cannot delete file %d with error: %s", index+1, err))
			return server.NewHTTPInternalError(fmt.Sprintf("cannot delete file %d with error: %s", index+1, err))
		}
	}

	if err := s.fileDB.Create(db, file); err != nil {
		logger.LogError(ctx, fmt.Sprintf("cannot create file %d with error: %s", index+1, err))
		return server.NewHTTPInternalError(fmt.Sprintf("cannot create file %d with error: %s", index+1, err))
	}

	return nil
}
