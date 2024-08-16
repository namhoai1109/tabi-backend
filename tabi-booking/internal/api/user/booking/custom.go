package booking

import (
	"tabi-booking/internal/model"
	"time"
)

// swagger: model UserBookingRequest
type UserBookingRequest struct {
	RoomID        int       `json:"room_id" validate:"required"`
	CheckInDate   time.Time `json:"check_in_date" validate:"required"`
	CheckOutDate  time.Time `json:"check_out_date" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	TotalPrice    float64   `json:"total_price" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required,min=1"`
	Note          string    `json:"note,omitempty"`
}

// swagger: model CancelBookingRequest
type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required"`
}

type UserBookingListResponse struct {
	Data  []*model.Booking `json:"data"`
	Total int64            `json:"total"`
}
