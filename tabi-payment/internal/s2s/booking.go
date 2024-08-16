package s2s

import (
	"context"
	"encoding/json"
	"fmt"
	structutil "tabi-payment/internal/util/struct"
	"time"

	"github.com/namhoai1109/tabi/core/logger"
)

//swagger:model BookingRequest
type BookingRequest struct {
	UserID        int       `json:"user_id" validate:"required"`
	RoomID        int       `json:"room_id" validate:"required"`
	CheckInDate   time.Time `json:"check_in_date" validate:"required"`
	CheckOutDate  time.Time `json:"check_out_date" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	TotalPrice    float64   `json:"total_price" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required"`
	Note          string    `json:"note,omitempty"`
}

type PaymentInfoResponse struct {
	RoomName    string `json:"room_name"`
	PaypalPayee string `json:"paypal_payee"`
}

func (s *S2S) CreateBooking(ctx context.Context, req *BookingRequest) error {
	creation := structutil.ToMap(req)
	resp, err := s.s2s.Post(ctx, creation, s.cfg.BookingEndpoint, "/partner/bookings")
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to create booking: %v", err))
		return fmt.Errorf("failed to create booking: %v", err)
	}

	return s.s2s.BuildError(resp)
}

func (s *S2S) GetPaymentInfo(ctx context.Context, roomID int) (*PaymentInfoResponse, error) {

	resp, err := s.s2s.Get(ctx, s.cfg.BookingEndpoint, fmt.Sprintf("/partner/rooms/%d/payment-info", roomID))
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to get payment info: %v", err))
		return nil, fmt.Errorf("failed to get payment info: %v", err)
	}

	data := PaymentInfoResponse{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, err
	}

	return &data, nil
}
