package roomtype

import (
	"context"
	"fmt"
	"strings"
	"tabi-booking/internal/model"
	structutil "tabi-booking/internal/util/struct"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *RoomType) List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomTypeResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}

	lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
		"type_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			val = "%" + strings.ToLower(val.(string)) + "%"
			return []interface{}{
				`lower(type_name) like ?`, val}
		}})

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
	}
	lq.Filter.And(`branch_id = ? AND linked = true`, branch.ID)
	roomTypes, err := s.roomTypeOfBranchDB.GetRoomTypes(s.db, ctx, lq, count)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting room types: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting room types, err: %v", err))
	}

	return s.convertToRoomTypeListResponse(ctx, roomTypes)
}

func (s *RoomType) ListAll(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomTypeResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
	}

	lq.Filter.And(`branch_id = ?  AND linked = false`, branch.ID)
	roomTypes, err := s.roomTypeOfBranchDB.GetRoomTypes(s.db, ctx, lq, count)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting room types: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting room types, err: %v", err))
	}

	return s.convertToRoomTypeListResponse(ctx, roomTypes)
}

func (s *RoomType) Create(ctx context.Context, authoPartner *model.AuthoPartner, req *RoomTypeCreateRequest) (*model.RoomType, error) {
	if err := s.enforce(authoPartner, model.ActionCreate); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return nil, server.NewHTTPInternalError("Branch not found")
	}

	if len(req.RoomFacilities) == 0 {
		return nil, server.NewHTTPValidationError("Room facilities must not be empty")
	}

	checkInTimeAfter12pm := req.CheckInTime.Hour() > 12 || (req.CheckInTime.Hour() == 12 && req.CheckInTime.Minute() > 0)
	checkOutTimeBefore12pm := req.CheckOutTime.Hour() < 12 || (req.CheckOutTime.Hour() == 12 && req.CheckOutTime.Minute() == 0)
	if !checkInTimeAfter12pm || !checkOutTimeBefore12pm {
		return nil, server.NewHTTPValidationError("Check in time must be after 12:00 and check out time must be before or equal 12:00")
	}

	creationRoomType := &model.RoomType{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {

		// create room type
		creationRoomType = &model.RoomType{
			TypeName:         req.TypeName,
			CheckInTime:      req.CheckInTime,
			CheckOutTime:     req.CheckOutTime,
			IncludeBreakfast: req.IncludeBreakfast,
			RoomFacilities:   req.RoomFacilities,
		}
		if err := s.roomTypeDB.Create(tx, creationRoomType); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating room type: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when creating room type, err: %v", err))
		}

		// link room type to branch
		roomTypeOfBranch := &model.RoomTypeOfBranch{
			RoomTypeID: creationRoomType.ID,
			BranchID:   branch.ID,
			Linked:     true,
		}

		if err := s.roomTypeOfBranchDB.Create(tx, roomTypeOfBranch); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating room type of branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when creating room type of branch, err: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return creationRoomType, nil
}

func (s *RoomType) Update(ctx context.Context, authoPartner *model.AuthoPartner, roomTypeID int, update RoomTypeUpdateRequest) (*model.RoomType, error) {
	if err := s.enforce(authoPartner, model.ActionUpdate); err != nil {
		return nil, err
	}

	if err := s.validateRoomType(ctx, authoPartner.ID, roomTypeID); err != nil {
		return nil, err
	}

	updates := structutil.ToMap(update)
	if len(updates) == 0 {
		return nil, server.NewHTTPValidationError("Data is empty")
	}

	if err := s.roomTypeDB.Update(s.db, updates, roomTypeID); err != nil {
		errMsg := fmt.Sprintf("Failed to update room type: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return s.view(ctx, roomTypeID)
}

func (s *RoomType) ChangeLinkStatus(ctx context.Context, authoPartner *model.AuthoPartner, req LinkRoomTypeRequest) (*string, error) {
	if err := s.enforce(authoPartner, model.ActionUpdate); err != nil {
		return nil, err
	}

	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		branch := &model.Branch{}
		if err := s.branchDB.View(tx, &branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
		}

		exist, err := s.roomTypeOfBranchDB.Exist(tx, `room_type_id = ? AND branch_id = ?`, *req.RoomTypeID, branch.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking existed room type of branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Error when checking existed room type of branch, err: %v", err))
		}

		if !exist {
			creation := &model.RoomTypeOfBranch{
				RoomTypeID: *req.RoomTypeID,
				BranchID:   branch.ID,
				Linked:     *req.LinkStatus,
			}

			if err := s.roomTypeOfBranchDB.Create(tx, creation); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when creating room type of branch: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Error when creating room type of branch, err: %v", err))
			}
		} else {
			if err := s.roomTypeOfBranchDB.Update(tx, map[string]interface{}{
				"linked": *req.LinkStatus,
			}, `room_type_id = ? AND branch_id = ?`, *req.RoomTypeID, branch.ID); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when updating room type of branch: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Error when updating room type of branch, err: %v", err))
			}
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	resp := "LINKED"
	if !*req.LinkStatus {
		resp = "UNLINKED"
	}

	return &resp, nil
}
