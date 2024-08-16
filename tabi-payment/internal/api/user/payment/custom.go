package payment

import "time"

// swagger:model PaymentCreationRequest
type PaymentCreationRequest struct {
	RoomID   int     `json:"room_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required"`
}

// swagger:model PaymentCreationResponse
type PaymentCreationResponse struct {
	ApproveLink string `json:"approve_link"`
}

// swagger:model PaymentCaptureResponse
type PaymentCaptureResponse struct {
	Message string `json:"message"`
}

//swagger:model PaymentCaptureRequest
type PaymentCaptureRequest struct {
	RoomID        int       `json:"room_id" validate:"required"`
	CheckInDate   time.Time `json:"check_in_date" validate:"required"`
	CheckOutDate  time.Time `json:"check_out_date" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	TotalPrice    float64   `json:"total_price" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required,min=1"`
	Note          string    `json:"note,omitempty"`
}
