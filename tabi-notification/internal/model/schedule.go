package model

import "time"

type Schedule struct {
	ID        int       `json:"id" gorm:"primary_key"`
	UserID    int       `json:"user_id"`
	BookingID int       `json:"booking_id"`
	StartTime time.Time `json:"start_time" gorm:"type:timestamp"`
	EndTime   time.Time `json:"end_time" gorm:"type:timestamp"`

	DestinationID           int    `json:"destination_id"`
	DestinationName         string `json:"destination_name"`
	DestinationCategory     string `json:"destination_category"`
	DestinationLocation     string `json:"destination_location"`
	DestinationWebsite      string `json:"destination_website"`
	DestinationImage        string `json:"destination_image"`
	DestinationOpeningHours string `json:"destination_opening_hours"`
	DestinationLongitude    string `json:"destination_longitude"`
	DestinationLatitude     string `json:"destination_latitude"`

	IsNotified bool `json:"is_notified"`
	Base
}
