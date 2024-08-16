package payment

import (
	"context"
	"tabi-payment/internal/s2s"

	"github.com/namhoai1109/tabi/core/paypal"
)

func New(
	paypal Paypal,
	s2s S2S,
) *Payment {
	return &Payment{
		paypal: paypal,
		s2s:    s2s,
	}
}

type Payment struct {
	paypal Paypal
	s2s    S2S
}

type Paypal interface {
	paypal.Intf
}

type S2S interface {
	CreateBooking(ctx context.Context, req *s2s.BookingRequest) error
	GetPaymentInfo(ctx context.Context, roomID int) (*s2s.PaymentInfoResponse, error)
}
