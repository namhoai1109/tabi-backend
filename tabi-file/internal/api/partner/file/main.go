package file

import (
	"context"
	"fmt"
	"tabi-file/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"github.com/thoas/go-funk"
)

func (s *File) ListFile(ctx context.Context, autho *model.AuthoPartner, lq *dbcore.ListQueryCondition) ([]*model.FileResponse, error) {
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

func (s *File) DeleteFile(ctx context.Context, autho *model.AuthoPartner, IDs []int) error {
	if !funk.Contains(model.PartnerRoles, autho.Role) {
		return server.NewHTTPAuthorizationError("You are not allowed to delete file")
	}

	if len(IDs) == 0 {
		logger.LogWarn(ctx, "There is no file to delete")
		return server.NewHTTPValidationError("There is no file to delete")
	}

	if err := s.fileDB.Delete(s.db, `id IN (?)`, IDs); err != nil {
		logger.LogError(ctx, fmt.Sprintf("cannot delete file with error: %s", err))
		return server.NewHTTPInternalError(fmt.Sprintf("cannot delete file with error: %s", err))
	}

	return nil
}
