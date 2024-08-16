package schedule

import (
	"context"
	"fmt"
	"tabi-notification/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Schedule) List(ctx context.Context, lq *dbcore.ListQueryCondition) ([]*model.Schedule, error) {
	schedules := []*model.Schedule{}
	if err := s.scheduleDB.List(s.db, &schedules, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error listing schedules: %v", err))
		return nil, server.NewHTTPInternalError("error listing schedules")
	}

	return schedules, nil
}

func (s *Schedule) MarkNotified(ctx context.Context, id int) error {
	if err := s.scheduleDB.Update(s.db, map[string]interface{}{
		"is_notified": true,
	}, id); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error marking schedule as notified: %v", err))
		return server.NewHTTPInternalError("error marking schedule as notified")
	}

	return nil
}
