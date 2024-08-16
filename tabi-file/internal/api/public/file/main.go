package file

import (
	"context"
	"fmt"
	"tabi-file/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *File) ListFile(ctx context.Context, lq *dbcore.ListQueryCondition) ([]*model.FileResponse, error) {
	var files []*model.File
	if err := s.fileDB.List(s.db.WithContext(ctx), &files, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when listing files: %v", err))
		return nil, server.NewHTTPInternalError("Error listing file").SetInternal(err)
	}

	var fileResponses []*model.FileResponse
	for _, file := range files {
		fileResponses = append(fileResponses, file.ToFileResponse())
	}

	return fileResponses, nil
}

func (s *File) UploadFilePreSigned(ctx context.Context, request UploadPresignedRequest) ([]*model.FileResponse, error) {

	// Make sure the files is provided.
	if len(request.Data) == 0 {
		logger.LogError(ctx, "There is no file to upload")
		return nil, server.NewHTTPValidationError("There is no file to upload")
	}

	var size int64
	var filesResponse []*model.FileResponse

	trxErr := s.db.Transaction(func(db *gorm.DB) error {
		for index, file := range request.Data {
			size = int64(file.FileSize)

			filename, err := s.validateFileName(file.FileName, index)
			if err != nil {
				logger.LogError(ctx, fmt.Sprintf("filename is invalid with error: %s", err))
				return err
			}

			pathName := fmt.Sprintf("/%s/%d/%s/", file.AttachmentType, file.AttachmentID, file.Field)
			key := fmt.Sprint(pathName, filename)

			// use the pre-signed URL to put an object to S3
			putURL, err := s.s3.PreparePresignedURL(key, s.cfg.S3PublicBucketName, s.cfg.ExpireTime)
			if err != nil {
				logger.LogError(ctx, fmt.Sprintf("cannot create PUT pre-signed URL at file %d with error: %s", index+1, err))
				return server.NewHTTPInternalError(fmt.Sprintf("cannot create PUT pre-signed URL at file %d with error: %s", index+1, err))
			}
			if putURL == nil {
				logger.LogError(ctx, "cannot create PUT pre-signed URL")
				return server.NewHTTPInternalError("cannot create PUT pre-signed URL")
			}

			new := &model.File{
				PathName:       pathName,
				FileName:       filename,
				ContentType:    file.ContentType,
				AttachmentID:   file.AttachmentID,
				AttachmentType: file.AttachmentType,
				Field:          file.Field,
				FileSize:       size,
				FullURL:        s.mappingFileURL(key),
				UploadURL:      *putURL,
				SortOrder:      file.SortOrder,
			}

			if err := s.createFile(db, ctx, new, index); err != nil {
				return err
			}

			filesResponse = append(filesResponse, new.ToFileResponse())
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return filesResponse, nil
}
