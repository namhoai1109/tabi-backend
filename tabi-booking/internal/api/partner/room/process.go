package room

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/rbac"
	"github.com/namhoai1109/tabi/core/server"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

func (s *Room) validateRoomType(db *gorm.DB, ctx context.Context, bmID int, roomTypeID int) (*model.Branch, error) {
	branch := &model.Branch{}
	if err := s.branchDB.View(db, &branch, `branch_manager_id = ?`, bmID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
	}

	exist, err := s.roomTypeOfBranchDB.Exist(db, `branch_id = ? AND room_type_id = ? AND linked = true`, branch.ID, roomTypeID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking exist room type of branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist room type of branch, err: %v", err))
	}

	if !exist {
		return nil, server.NewHTTPValidationError("Room type not found")
	}

	return branch, nil
}

func (s *Room) validateBedType(db *gorm.DB, ctx context.Context, bedTypeID int) error {
	exist, err := s.generalTypeDB.Exist(db, `id = ? AND class = ?`, bedTypeID, model.GeneralTypeClassBed)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking exist bed type: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist bed type, err: %v", err))
	}

	if !exist {
		return server.NewHTTPValidationError("Bed type not found")
	}

	return nil
}

func (s *Room) toListResponse(rooms []*model.Room, count int64) *RoomListResponse {
	res := []RoomListItemResponse{}
	for _, room := range rooms {
		res = append(res, RoomListItemResponse{
			ID:                room.ID,
			RoomType:          room.RoomType.TypeName,
			RoomName:          room.RoomName,
			Status:            room.Status,
			BranchName:        room.Branch.BranchName,
			BranchManagerName: room.Branch.BranchManager.Name,
			Quantity:          room.Quantity,
			MaxPrice:          room.MaxPrice,
		})
	}

	return &RoomListResponse{
		Total: count,
		Data:  res,
	}
}

func (s *Room) toViewRoomResponse(room *model.Room, facilities []*model.Facility) *ViewRoomResponse {
	bedType := &model.GeneralType{
		ID:      room.BedType.ID,
		Class:   room.BedType.Class,
		LabelEN: room.BedType.LabelEN,
		LabelVI: room.BedType.LabelVI,
		DescEN:  room.BedType.DescEN,
		DescVI:  room.BedType.DescVI,
		Order:   room.BedType.Order,
	}
	factureReduction := &model.FactureReduction{
		ID:           room.FactureReduction.ID,
		OnlineMethod: room.FactureReduction.OnlineMethod,
		OnCashMethod: room.FactureReduction.OnCashMethod,
		NormalDay:    room.FactureReduction.NormalDay,
		Holiday:      room.FactureReduction.Holiday,
		Weekend:      room.FactureReduction.Weekend,
	}
	roomType := CustomRoomType{
		RoomType:   *room.RoomType,
		Facilities: facilities,
	}
	var reservationReductions []*model.ReservationReduction
	for _, reservation := range room.ReservationReduction {
		reservationReductions = append(reservationReductions, &model.ReservationReduction{
			ID:        reservation.ID,
			Quantity:  reservation.Quantity,
			TimeUnit:  reservation.TimeUnit,
			Reduction: reservation.Reduction,
		})
	}

	return &ViewRoomResponse{
		ID:           room.ID,
		RoomName:     room.RoomName,
		MaxOccupancy: room.MaxOccupancy,
		Status:       room.Status,
		Width:        room.Width,
		Length:       room.Length,
		MaxPrice:     room.MaxPrice,
		Quantity:     room.Quantity,

		BedType:               bedType,
		FactureReduction:      factureReduction,
		ReservationReductions: reservationReductions,
		RoomType:              roomType,
	}
}

func (s *Room) checkRoomOwnership(ctx context.Context, authoPartner *model.AuthoPartner, roomID int) error {
	if authoPartner.Role == model.AccountRoleBranchManager || authoPartner.Role == model.AccountRoleHost {
		branch := &model.Branch{}
		if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
		}

		exist, err := s.roomDB.Exist(s.db, `id = ? AND branch_id = ?`, roomID, branch.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking exist room: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist room, err: %v", err))
		}

		if !exist {
			return server.NewHTTPValidationError("Room not found or not belong to you")
		}

		return nil
	}

	company := &model.Company{}
	if err := s.companyDB.View(s.db.Preload("Branches"), &company, `representative_id = ?`, authoPartner.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting company: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when getting company, err: %v", err))
	}

	branchIDs := make([]int, 0)
	for _, branch := range company.Branches {
		branchIDs = append(branchIDs, branch.ID)
	}
	exist, err := s.roomDB.Exist(s.db, `id = ? AND branch_id IN (?)`, roomID, branchIDs)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking exist room: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist room, err: %v", err))
	}

	if !exist {
		return server.NewHTTPValidationError("Room not found or not belong to you")
	}

	return nil
}

func (s *Room) insertReservationReduction(db *gorm.DB, ctx context.Context, roomID int, reservationReductions []*ReservationReductionRequest) error {
	for _, reservationReduction := range reservationReductions {

		if !util.InSliceString(model.ReservationReductionTimeUnits, reservationReduction.TimeUnit) {
			return server.NewHTTPValidationError(fmt.Sprintf("Invalid time unit: %s", reservationReduction.TimeUnit))
		}

		reservationReduction := &model.ReservationReduction{
			RoomID:    roomID,
			Quantity:  *reservationReduction.Quantity,
			TimeUnit:  reservationReduction.TimeUnit,
			Reduction: *reservationReduction.Reduction,
		}

		if err := s.reservationReductionDB.Create(db, reservationReduction); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating reservation reduction: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when creating reservation reduction, err: %v", err))
		}
	}

	return nil
}

func (s *Room) getFactureReductionUpdate(updateMap map[string]interface{}) map[string]interface{} {
	frMap := map[string]interface{}{}
	frField := []string{"online_method", "on_cash_method", "normal_day", "holiday", "weekend"}

	for key, value := range updateMap {
		if funk.Contains(frField, key) {
			frMap[key] = *value.(*float64)
			delete(updateMap, key)
		}
	}

	return frMap
}

func (s *Room) validateRoomUpdateForBM(db *gorm.DB, ctx context.Context, roomID int, bmID int, updateMap map[string]interface{}, updates UpdateRoomRequest) (map[string]interface{}, error) {
	if updateMap["quantity"] != nil {
		if updates.Quantity < 0 {
			return nil, server.NewHTTPValidationError("Invalid quantity")
		}
	}

	if updateMap["bed_type_id"] != nil {
		if err := s.validateBedType(db, ctx, updates.BedTypeID); err != nil {
			return nil, err
		}
	}

	if updateMap["room_type_id"] != nil {
		_, err := s.validateRoomType(db, ctx, bmID, updates.RoomTypeID)
		if err != nil {
			return nil, err
		}

	}

	if updateMap["status"] != nil && updateMap["status"] != model.RoomStatusUpdated {
		return nil, server.NewHTTPValidationError("Invalid status")
	}

	frUpdate := s.getFactureReductionUpdate(updateMap)
	if len(frUpdate) > 0 {

		if err := s.factureReductionDB.Update(db, frUpdate, `room_id = ?`, roomID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when updating facture reduction: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when updating facture reduction, err: %v", err))
		}
	}

	if updates.ReservationReduction != nil && len(updates.ReservationReduction) > 0 {
		if err := s.reservationReductionDB.Delete(db, `room_id = ?`, roomID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when deleting reservation reduction: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when deleting reservation reduction, err: %v", err))
		}

		if err := s.insertReservationReduction(db, ctx, roomID, updates.ReservationReduction); err != nil {
			return nil, err
		}

		delete(updateMap, "reservation_reduction")
	}

	return updateMap, nil
}

// enforce checks permission to perform the action
func (s *Room) enforce(authPartner *model.AuthoPartner, action string) error {
	if !s.rbac.Enforce(authPartner.Role, model.ObjectRoom, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
