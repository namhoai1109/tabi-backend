package booking

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"time"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Booking) Create(ctx context.Context, autho *model.AuthoUser, req *UserBookingRequest) error {
	if autho.Role != model.AccountRoleClient {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	room := &model.Room{}
	if err := s.roomDB.View(s.db.Preload("Bookings"), room, req.RoomID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when view room: %v", err))
		return server.NewHTTPValidationError("Room not found")
	}
	if room.CountAvailableRoom(req.CheckInDate, req.CheckOutDate) < req.Quantity {
		return server.NewHTTPValidationError("Room is not available or quantity is not enough")
	}

	status := model.BookingStatusPending
	if req.PaymentMethod == model.BookingPaymentMethodOnline {
		status = model.BookingStatusApproved
	}

	booking := &model.Booking{
		UserID:        autho.ID,
		RoomID:        req.RoomID,
		CheckInDate:   &req.CheckInDate,
		CheckOutDate:  &req.CheckOutDate,
		TotalPrice:    req.TotalPrice,
		Status:        status,
		Note:          req.Note,
		PaymentMethod: req.PaymentMethod,
		Quantity:      req.Quantity,
	}

	if err := s.bookingDB.Create(s.db, booking); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when create booking: %v", err))
		return server.NewHTTPInternalError("Error when create booking")
	}

	return nil
}

func (s *Booking) List(ctx context.Context, autho *model.AuthoUser, lq *dbcore.ListQueryCondition, count *int64) ([]*model.Booking, error) {
	if autho.Role != model.AccountRoleClient {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	bookings := []*model.Booking{}
	lq.Filter.And("user_id = ?", autho.ID)

	preloadDB := s.db.Preload("Room").
		Preload("Room.Branch").
		Preload("Room.RoomType").
		Preload("Room.BedType")

	if err := s.bookingDB.List(preloadDB, &bookings, lq, count); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when list booking: %v", err))
		return nil, server.NewHTTPInternalError("Error when list booking")
	}

	return bookings, nil
}

func (s *Booking) Cancel(ctx context.Context, autho *model.AuthoUser, bookingID int, req CancelBookingRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if autho.Role != model.AccountRoleClient {
			return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
		}

		booking := &model.Booking{}
		if err := s.bookingDB.View(tx.Preload("Room").Preload("Room.Branch"), &booking, `user_id = ? AND id = ? AND status IN (?) AND payment_method = ?`,
			autho.ID, bookingID,
			[]string{
				model.BookingStatusApproved,
				model.BookingStatusPending},
			model.BookingPaymentMethodCash); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when view booking: %v", err))
			return server.NewHTTPValidationError("Booking not found")
		}

		branch := booking.Room.Branch
		if branch == nil {
			return server.NewHTTPInternalError("Branch not found")
		}

		cancellationTime := time.Now()
		cancellationTimeOfBranch := branch.GetCancellationTime()
		if cancellationTime.After(booking.CheckInDate.Add(-cancellationTimeOfBranch)) {
			return server.NewHTTPValidationError("\t\t\t\t\t\t You can not cancel this booking.\n Please contact branch for more information!")
		}

		if err := s.bookingDB.Update(tx, map[string]interface{}{
			"status": model.BookingStatusCancel,
			"reason": req.Reason,
		}, bookingID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when cancel booking: %v", err))
			return server.NewHTTPInternalError("Error when cancel booking")
		}

		return nil
	})
}
