package roomtype

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/rbac"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *RoomType) validateRoomType(ctx context.Context, bmID, roomTypeID int) error {

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, bmID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
	}

	exist, err := s.roomTypeOfBranchDB.Exist(s.db, `branch_id = ? AND room_type_id = ?`, branch.ID, roomTypeID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking existed room type: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking existed room type, err: %v", err))
	}

	if !exist {
		return server.NewHTTPValidationError("Room type not found")
	}

	return nil
}

func (s *RoomType) view(ctx context.Context, roomTypeID int) (*model.RoomType, error) {
	roomType := &model.RoomType{}
	if err := s.roomTypeDB.View(s.db, &roomType, `id = ?`, roomTypeID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to view room type: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view room type: %v", err))
	}

	return roomType, nil
}

func (s *RoomType) convertToRoomTypeListResponse(ctx context.Context, roomTypes []*model.RoomType) ([]*model.RoomTypeResponse, error) {
	roomTypeResponses := make([]*model.RoomTypeResponse, len(roomTypes))
	for i, roomType := range roomTypes {
		facilities, err := s.facilityDB.GetFacilityList(s.db, roomType.RoomFacilities)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get facility list: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get facility list: %v", err))
		}
		roomTypeResponse := roomType.ToResponse(facilities)
		roomTypeResponses[i] = roomTypeResponse
	}

	return roomTypeResponses, nil
}

// enforce checks permission to perform the action
func (s *RoomType) enforce(authPartner *model.AuthoPartner, action string) error {
	if !s.rbac.Enforce(authPartner.Role, model.ObjectRoomType, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
