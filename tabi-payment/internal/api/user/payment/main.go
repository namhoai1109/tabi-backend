package payment

import (
	"context"
	"fmt"
	"tabi-payment/internal/model"
	"tabi-payment/internal/s2s"

	"github.com/namhoai1109/tabi/core/logger"
	paypalmodel "github.com/namhoai1109/tabi/core/paypal/model"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Payment) Create(ctx context.Context, authoUser *model.AuthoUser, req *PaymentCreationRequest) (*PaymentCreationResponse, error) {
	if authoUser.Role != model.UserRole {
		return nil, server.NewHTTPAuthorizationError("User role is invalid")
	}

	paymentInfo, err := s.s2s.GetPaymentInfo(ctx, req.RoomID)
	if err != nil {
		message := fmt.Sprintf("Failed to get payment info: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	orderReq := s.getOrderRequest(req.Price, paymentInfo.RoomName, req.Quantity, paymentInfo.PaypalPayee)

	resp, err := s.paypal.CreateOrder(ctx, orderReq)
	if err != nil {
		message := fmt.Sprintf("Failed to create Paypal order: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	if resp.Status != paypalmodel.OrderStatusPayerActionRequired {
		message := fmt.Sprintf("Order status is invalid: %s", resp.Status)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	approveRel := ""
	for _, link := range resp.Links {
		if link.Rel == paypalmodel.OrderRefPayerAction {
			approveRel = link.Href
			break
		}
	}

	if approveRel == "" {
		logger.LogError(ctx, "Approve link not found")
		return nil, server.NewHTTPInternalError("Approve link not found")
	}

	return &PaymentCreationResponse{
		ApproveLink: approveRel,
	}, nil
}

func (s *Payment) Capture(ctx context.Context, authoUser *model.AuthoUser, orderID string, req *PaymentCaptureRequest) (*PaymentCaptureResponse, error) {
	if authoUser.Role != model.UserRole {
		return nil, server.NewHTTPAuthorizationError("User role is invalid")
	}

	resp, err := s.paypal.CaptureOrder(ctx, orderID)
	if err != nil {
		message := fmt.Sprintf("Failed to capture Paypal order: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}
	logger.LogInfo(ctx, fmt.Sprintf("Capture order response: %+v", resp))
	if resp.Status != paypalmodel.OrderStatusCompleted {
		message := fmt.Sprintf("Order status is invalid: %s", resp.Status)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	if err := s.s2s.CreateBooking(ctx, &s2s.BookingRequest{
		UserID:        authoUser.ID,
		RoomID:        req.RoomID,
		TotalPrice:    req.TotalPrice,
		CheckInDate:   req.CheckInDate,
		CheckOutDate:  req.CheckOutDate,
		PaymentMethod: req.PaymentMethod,
		Quantity:      req.Quantity,
		Note:          req.Note,
	}); err != nil {
		message := fmt.Sprintf("Failed to create booking: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	return &PaymentCaptureResponse{
		Message: paypalmodel.OrderStatusCompleted,
	}, nil
}
