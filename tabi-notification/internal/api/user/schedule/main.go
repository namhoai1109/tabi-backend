package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"tabi-notification/internal/model"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Schedule) Create(ctx context.Context, autho *model.AuthoUser, req ScheduleCreationRequest) error {
	if autho.Role != model.ClientRole {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Schedules {
			schedule := model.Schedule{
				UserID:                  autho.ID,
				BookingID:               item.BookingID,
				StartTime:               item.StartTime,
				EndTime:                 item.EndTime,
				DestinationID:           item.DestinationID,
				DestinationName:         item.DestinationName,
				DestinationCategory:     item.DestinationCategory,
				DestinationLocation:     item.DestinationLocation,
				DestinationWebsite:      item.DestinationWebsite,
				DestinationImage:        item.DestinationImage,
				DestinationOpeningHours: item.DestinationOpeningHours,
				DestinationLongitude:    item.DestinationLongitude,
				DestinationLatitude:     item.DestinationLatitude,
			}

			exist, err := s.scheduleDB.Exist(tx, `user_id = ? AND booking_id = ? AND start_time = ? AND destination_id = ?`, autho.ID, item.BookingID, item.StartTime, item.DestinationID)
			if err != nil {
				logger.LogError(ctx, fmt.Sprintf("error checking if schedule exist: %v", err))
				return server.NewHTTPInternalError("error checking if schedule exist or overlap")
			}

			if exist {
				return server.NewHTTPValidationError("schedule already exist")
			}

			if err := s.scheduleDB.Create(tx, &schedule); err != nil {
				logger.LogError(ctx, fmt.Sprintf("error creating schedule: %v", err))
				return server.NewHTTPInternalError("error creating schedule")
			}
		}

		return nil
	})

	if trxErr != nil {
		return trxErr
	}

	return nil
}

func (s *Schedule) View(ctx context.Context, autho *model.AuthoUser, id int) (*ScheduleResponse, error) {
	if autho.Role != model.ClientRole {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	schedule := model.Schedule{}
	if err := s.scheduleDB.View(s.db, &schedule, `id = ? AND user_id = ?`, id, autho.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error getting schedule: %v", err))
		return nil, server.NewHTTPInternalError("error getting schedule")
	}

	return s.toResponse(&schedule), nil
}

func (s *Schedule) List(ctx context.Context, autho *model.AuthoUser, bookingID int, lq *dbcore.ListQueryCondition) (*ScheduleListResponse, error) {
	if autho.Role != model.ClientRole {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	schedules := []*model.Schedule{}
	count := int64(0)
	lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
		"time_range": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			times := []string{}

			if err := json.Unmarshal([]byte(val.(string)), &times); err != nil {
				logger.LogError(ctx, fmt.Sprintf("error unmarshal times: %v", err))
				return nil
			}

			from, to := times[0], times[1]

			return []interface{}{
				`start_time >= ? AND end_time <= ?`, from, to,
			}
		},
	})
	lq.Filter.And(`user_id = ? AND booking_id = ?`, autho.ID, bookingID)
	if err := s.scheduleDB.List(s.db, &schedules, lq, &count); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error listing schedules: %v", err))
		return nil, server.NewHTTPInternalError("error listing schedules")
	}

	return s.toListResponse(schedules, count), nil
}

func (s *Schedule) Delete(ctx context.Context, autho *model.AuthoUser, id int) error {
	if autho.Role != model.ClientRole {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	if err := s.scheduleDB.Delete(s.db, `id = ? AND user_id = ?`, id, autho.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error deleting schedule: %v", err))
		return server.NewHTTPInternalError("error deleting schedule")
	}

	return nil
}

func (s *Schedule) Update(ctx context.Context, autho *model.AuthoUser, id int, req ScheduleCreation) error {
	if autho.Role != model.ClientRole {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	if err := s.scheduleDB.Update(s.db, req, `id = ? AND user_id = ?`, id, autho.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error updating schedule: %v", err))
		return server.NewHTTPInternalError("error updating schedule")
	}

	return nil
}

func (s *Schedule) DeleteIDs(ctx context.Context, autho *model.AuthoUser, ids []int) error {
	if autho.Role != model.ClientRole {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	if err := s.scheduleDB.Delete(s.db, `id IN (?) AND user_id = ?`, ids, autho.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error deleting schedules: %v", err))
		return server.NewHTTPInternalError("error deleting schedules")
	}

	return nil
}
