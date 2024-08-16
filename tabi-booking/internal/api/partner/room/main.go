package room

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"
	structutil "tabi-booking/internal/util/struct"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Room) Create(ctx context.Context, authoPartner *model.AuthoPartner, req CreateRoomRequest) (*model.Room, error) {
	if err := s.enforce(authoPartner, model.ActionCreate); err != nil {
		return nil, err
	}

	room := &model.Room{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {

		// * insert room
		branch, err := s.validateRoomType(tx, ctx, authoPartner.ID, req.RoomTypeID)
		if err != nil {
			return err
		}

		if err := s.validateBedType(tx, ctx, req.BedTypeID); err != nil {
			return err
		}

		room = &model.Room{
			RoomTypeID:   req.RoomTypeID,
			BedTypeID:    req.BedTypeID,
			Status:       model.RoomStatusPending,
			RoomName:     req.RoomName,
			MaxOccupancy: req.MaxOccupancy,
			MaxPrice:     req.MaxPrice,
			Width:        req.Width,
			Length:       req.Length,
			Quantity:     req.Quantity,
			BranchID:     branch.ID,
		}

		// if bm is host, room status is updated
		if authoPartner.Role == model.AccountRoleHost {
			room.Status = model.RoomStatusUpdated
		}

		if err := s.roomDB.Create(tx, room); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating room: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when creating room, err: %v", err))
		}

		// * insert facture reduction
		factureReduction := &model.FactureReduction{
			RoomID:       room.ID,
			OnlineMethod: *req.OnlineMethod,
			OnCashMethod: *req.OnCashMethod,
			NormalDay:    *req.NormalDay,
			Holiday:      *req.Holiday,
			Weekend:      *req.Weekend,
		}

		if err := s.factureReductionDB.Create(tx, factureReduction); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating facture reduction: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when creating facture reduction, err: %v", err))
		}

		// * insert reservation reduction
		return s.insertReservationReduction(tx, ctx, room.ID, req.ReservationReduction)
	})

	return room, trxErr
}

func (s *Room) List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition) (*RoomListResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}

	list := []*model.Room{}
	count := int64(0)
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		if authoPartner.Role == model.AccountRoleBranchManager || authoPartner.Role == model.AccountRoleHost {
			branchID, err := s.branchDB.GetBranchIDByBranchManagerID(tx, authoPartner.ID)
			if err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when getting branch id: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch id, err: %v", err))
			}
			lq.Filter.And("branch_id = ?", branchID)
		} else {
			branchIDs, err := s.branchDB.GetBranchIDsByRepresentativeID(tx, authoPartner.ID)
			if err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when getting branch ids: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch ids, err: %v", err))
			}
			lq.Filter.And("branch_id IN (?)", branchIDs)
		}

		lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
			"room_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(room_name) like ?`, val}
			},
			"branch_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(branch.branch_name) like ?`, val}
			},
			"room_type": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				return []interface{}{
					`room_type.type_name = ?`, val.(string)}
			},
			"room_type__in": func(_ string, val interface{}, _ *gowhere.Config) interface{} {

				types := []string{}
				if err := json.Unmarshal([]byte(val.(string)), &types); err != nil {
					return nil
				}

				return []interface{}{
					`room_type.type_name IN (?)`, types}
			},
		})

		query := tx.Preload("RoomType").Preload("Branch").Preload("Branch.BranchManager").
			Joins("JOIN branch ON branch.id = room.branch_id").
			Joins("JOIN room_type ON room_type.id = room.room_type_id")
		if err := s.roomDB.List(query, &list, lq, &count); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when listing room: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when listing room, err: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.toListResponse(list, count), nil
}

func (s *Room) View(ctx context.Context, authoPartner *model.AuthoPartner, roomID int) (*ViewRoomResponse, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	if err := s.checkRoomOwnership(ctx, authoPartner, roomID); err != nil {
		return nil, err
	}

	room := &model.Room{}
	var facilities []*model.Facility
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Preload("RoomType").
			Preload("BedType").
			Preload("FactureReduction").
			Preload("ReservationReduction")

		if err := s.roomDB.View(query, room, `id = ?`, roomID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when getting room: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when getting room, err: %v", err))
		}

		values, err := s.facilityDB.GetFacilityList(tx, room.RoomType.RoomFacilities)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when getting facility list: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when getting facility list, err: %v", err))
		}
		facilities = values

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.toViewRoomResponse(room, facilities), nil
}

func (s *Room) ListBookings(ctx context.Context, authoPartner *model.AuthoPartner, roomID int, lq *dbcore.ListQueryCondition, count *int64) ([]*model.BookingResponse, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	if err := s.checkRoomOwnership(ctx, authoPartner, roomID); err != nil {
		return nil, err
	}

	lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
		"date__in": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			year, month := util.GetMonthYear(val.(string))
			monthStr := fmt.Sprintf("%02d", month)
			yearStr := fmt.Sprintf("%04d", year)
			return []interface{}{
				`to_char(check_in_date, 'YYYY-MM') = ?`, fmt.Sprintf("%s-%s", yearStr, monthStr)}
		},
	})

	bookings := []*model.Booking{}
	lq.Sort = append(lq.Sort, "check_in_date ASC")
	lq.Filter.And("room_id = ? AND status = ?", roomID, model.BookingStatusApproved)
	preloadDB := s.db.Preload("User").Preload("User.Account")
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

func (s *Room) Update(ctx context.Context, authoPartner *model.AuthoPartner, updates UpdateRoomRequest, roomID int) (*ViewRoomResponse, error) {
	if err := s.enforce(authoPartner, model.ActionUpdate); err != nil {
		return nil, err
	}

	if authoPartner.Role == model.AccountRoleBranchManager || authoPartner.Role == model.AccountRoleHost {
		if err := s.checkRoomOwnership(ctx, authoPartner, roomID); err != nil {
			return nil, err
		}

		trxErr := s.db.Transaction(func(tx *gorm.DB) error {
			updateMap := structutil.ToMap(updates)

			newUpdates, err := s.validateRoomUpdateForBM(tx, ctx, roomID, authoPartner.ID, updateMap, updates)
			if err != nil {
				return err
			}

			if len(newUpdates) > 0 {
				if err := s.roomDB.Update(tx, newUpdates, `id = ?`, roomID); err != nil {
					logger.LogError(ctx, fmt.Sprintf("Error when updating room: %v", err))
					return server.NewHTTPInternalError(fmt.Sprintf("Error when updating room, err: %v", err))
				}
			}

			return nil
		})

		if trxErr != nil {
			return nil, trxErr
		}
	} else if authoPartner.Role == model.AccountRoleRepresentative {
		if err := s.checkRoomOwnership(ctx, authoPartner, roomID); err != nil {
			return nil, err
		}

		if updates.Status == model.RoomStatusUpdated {
			return nil, server.NewHTTPValidationError("Invalid status")
		}

		if err := s.roomDB.Update(s.db, map[string]interface{}{
			"status": updates.Status,
		}, `id = ?`, roomID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when updating room: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when updating room, err: %v", err))
		}
	}

	return s.View(ctx, authoPartner, roomID)
}

func (s *Room) GetPaymentInfo(ctx context.Context, authoPartner *model.AuthoPartner, roomID int) (*PaymentInfoResponse, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	room := &model.Room{}
	preloadDB := s.db.Preload("Branch").Preload("Branch.BranchManager").Preload("Branch.BranchManager.Account")

	if err := s.roomDB.View(preloadDB, &room, `id = ?`, roomID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting room: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting room, err: %v", err))
	}
	roomName := room.RoomName
	paypalPayee := ""
	bm := room.Branch.BranchManager
	if bm.Account.Role == model.AccountRoleHost {
		paypalPayee = bm.Account.Email
	} else {
		company := &model.Company{}
		preloadDB = s.db.Preload("Representative").
			Preload("Representative.Account")
		if err := s.companyDB.View(preloadDB, &company, `id = ?`, room.Branch.CompanyID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when getting company: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting company, err: %v", err))
		}
		paypalPayee = company.Representative.Account.Email
	}

	return &PaymentInfoResponse{
		RoomName:    roomName,
		PaypalPayee: paypalPayee,
	}, nil
}
