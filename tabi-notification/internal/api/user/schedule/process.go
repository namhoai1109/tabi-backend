package schedule

import "tabi-notification/internal/model"

func (s *Schedule) toResponse(schedule *model.Schedule) *ScheduleResponse {
	resp := &ScheduleResponse{
		ID:         schedule.ID,
		IsNotified: schedule.IsNotified,
		ScheduleCreation: ScheduleCreation{
			BookingID:               schedule.BookingID,
			StartTime:               schedule.StartTime,
			EndTime:                 schedule.EndTime,
			DestinationID:           schedule.DestinationID,
			DestinationName:         schedule.DestinationName,
			DestinationCategory:     schedule.DestinationCategory,
			DestinationLocation:     schedule.DestinationLocation,
			DestinationWebsite:      schedule.DestinationWebsite,
			DestinationImage:        schedule.DestinationImage,
			DestinationOpeningHours: schedule.DestinationOpeningHours,
			DestinationLongitude:    schedule.DestinationLongitude,
			DestinationLatitude:     schedule.DestinationLatitude,
		},
	}

	return resp
}

func (s *Schedule) toListResponse(list []*model.Schedule, count int64) *ScheduleListResponse {
	resp := []*ScheduleResponse{}
	for _, item := range list {
		resp = append(resp, s.toResponse(item))
	}

	return &ScheduleListResponse{
		Total: count,
		Data:  resp,
	}
}
