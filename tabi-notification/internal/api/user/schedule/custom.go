package schedule

import "time"

// swagger:model ScheduleCreation
type ScheduleCreation struct {
	BookingID               int       `json:"booking_id" validate:"required"`
	StartTime               time.Time `json:"start_time" validate:"required"`
	EndTime                 time.Time `json:"end_time" validate:"required"`
	DestinationID           int       `json:"destination_id" validate:"required"`
	DestinationName         string    `json:"destination_name" validate:"required"`
	DestinationCategory     string    `json:"destination_category" validate:"required"`
	DestinationLocation     string    `json:"destination_location" validate:"required"`
	DestinationWebsite      string    `json:"destination_website,omitempty"`
	DestinationImage        string    `json:"destination_image,omitempty"`
	DestinationOpeningHours string    `json:"destination_opening_hours,omitempty"`
	DestinationLongitude    string    `json:"destination_longitude,omitempty"`
	DestinationLatitude     string    `json:"destination_latitude,omitempty"`
}

// swagger:model ScheduleCreationRequest
type ScheduleCreationRequest struct {
	Schedules []*ScheduleCreation `json:"schedules" validate:"required"`
}

// swagger:model ScheduleResponse
type ScheduleResponse struct {
	ID         int  `json:"id"`
	IsNotified bool `json:"is_notified"`
	ScheduleCreation
}

// swagger:model ScheduleListResponse
type ScheduleListResponse struct {
	Total int64               `json:"total"`
	Data  []*ScheduleResponse `json:"data"`
}

type ScheduleDeleteRequest struct {
	IDs []int `json:"ids" validate:"required"`
}
