package device

import (
	"context"
	"fmt"
	"tabi-notification/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Device) List(ctx context.Context, lq *dbcore.ListQueryCondition) ([]*model.Device, error) {
	devices := []*model.Device{}
	if err := s.deviceDB.List(s.db, &devices, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error listing devices: %v", err))
		return nil, server.NewHTTPInternalError("error listing devices")
	}

	return devices, nil
}
