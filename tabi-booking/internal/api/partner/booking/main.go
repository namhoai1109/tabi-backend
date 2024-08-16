package booking

import (
	"context"
	"fmt"
	"strings"
	"tabi-booking/internal/model"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Booking) List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.BookingResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}

	lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
		"first_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			val = "%" + strings.ToLower(val.(string)) + "%"
			return []interface{}{
				`lower("user".first_name)like ?`, val}
		},
		"room_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			val = "%" + strings.ToLower(val.(string)) + "%"
			return []interface{}{
				`lower(room.room_name) like ?`, val}
		},
	})

	roomIDs, err := s.getRoomIDs(ctx, authoPartner.ID)
	if err != nil {
		return nil, err
	}
	lq.Sort = append(lq.Sort, "created_at DESC")
	lq.Filter.And("room_id IN (?)", roomIDs)
	bookings := []*model.Booking{}
	preloadDB := s.db.Preload("Room").Preload("User").Preload("User.Account").
		Joins("JOIN room ON room.id = booking.room_id").
		Joins("JOIN \"user\" ON \"user\".id = booking.user_id")

	if err := s.bookingDB.List(preloadDB, &bookings, lq, count); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when list booking: %v", err))
		return nil, server.NewHTTPInternalError("Error when list booking")
	}

	resp := []*model.BookingResponse{}
	for _, booking := range bookings {
		resp = append(resp, booking.ToResponse())
	}

	return resp, nil
}

// Create creates a new booking from payment service
func (s *Booking) Create(ctx context.Context, authoPartner *model.AuthoPartner, req *BookingRequest) error {
	if err := s.enforce(authoPartner, model.ActionCreate); err != nil {
		return err
	}

	exist, err := s.userDB.Exist(s.db, req.UserID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check user exist: %v", err))
		return server.NewHTTPInternalError("Error when check user exist")
	}

	if !exist {
		return server.NewHTTPValidationError("User not found")
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
		UserID:        req.UserID,
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

func (s *Booking) Approve(ctx context.Context, authoPartner *model.AuthoPartner, id int) error {
	if err := s.enforce(authoPartner, model.ActionUpdateAll); err != nil {
		return err
	}

	roomIDs, err := s.getRoomIDs(ctx, authoPartner.ID)
	if err != nil {
		return err
	}

	exist, err := s.bookingDB.Exist(s.db, `id = ? AND status = ? AND room_id IN (?)`, id, model.BookingStatusPending, roomIDs)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check booking exist: %v", err))
		return server.NewHTTPInternalError("Error when check booking exist")
	}

	if !exist {
		return server.NewHTTPValidationError("Booking not found or already approved")
	}

	if err := s.bookingDB.Update(s.db, map[string]interface{}{
		"status": model.BookingStatusApproved,
	}, id); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when update booking: %v", err))
		return server.NewHTTPInternalError("Error when update booking")
	}

	return nil
}

func (s *Booking) Reject(ctx context.Context, authoPartner *model.AuthoPartner, id int, req RejectReasonRequest) error {
	if err := s.enforce(authoPartner, model.ActionUpdateAll); err != nil {
		return err
	}

	roomIDs, err := s.getRoomIDs(ctx, authoPartner.ID)
	if err != nil {
		return err
	}

	exist, err := s.bookingDB.Exist(s.db, `id = ? AND status = ? AND room_id IN (?)`, id, model.BookingStatusPending, roomIDs)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check booking exist: %v", err))
		return server.NewHTTPInternalError("Error when check booking exist")
	}

	if !exist {
		return server.NewHTTPValidationError("Booking not found or already rejected")
	}

	if err := s.bookingDB.Update(s.db, map[string]interface{}{
		"status": model.BookingStatusRejected,
		"reason": req.Reason,
	}, id); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when update booking: %v", err))
		return server.NewHTTPInternalError("Error when update booking")
	}

	return nil
}
