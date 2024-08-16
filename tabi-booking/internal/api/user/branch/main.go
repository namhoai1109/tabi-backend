package branch

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"
	"tabi-booking/internal/util"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Branch) Save(ctx context.Context, autho *model.AuthoUser, branchID int, req SaveBranchRequest) (*SaveBranchResponse, error) {
	if autho.Role != model.AccountRoleClient {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	trxErr := s.db.Transaction(func(tx *gorm.DB) error {

		existedBranch, err := s.branchDB.Exist(tx, `id = ?`, branchID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking branch exist: %v", err))
			return server.NewHTTPInternalError("Error when checking branch exist")
		}

		if !existedBranch {
			return server.NewHTTPValidationError("Branch not found")
		}

		existedRecord, err := s.savedBranchDB.Exist(tx, `branch_id = ? AND user_id = ?`, branchID, autho.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking saved branch exist: %v", err))
			return server.NewHTTPInternalError("Error when checking saved branch exist")
		}

		if existedRecord && !*req.Save {
			if err := s.savedBranchDB.Delete(tx, `branch_id = ? AND user_id = ?`, branchID, autho.ID); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when deleting saved branch: %v", err))
				return server.NewHTTPInternalError("Error when deleting saved branch")
			}
		} else if !existedRecord && *req.Save {
			if err := s.savedBranchDB.Create(tx, &model.SavedBranch{
				BranchID: branchID,
				UserID:   autho.ID,
			}); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Error when creating saved branch: %v", err))
				return server.NewHTTPInternalError("Error when creating saved branch")
			}
		} else {
			return server.NewHTTPValidationError("Invalid request")
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return &SaveBranchResponse{
		Message: util.TernaryOperator(*req.Save, "Branch saved", "Branch unsaved").(string),
	}, nil
}

func (s *Branch) ListSaved(ctx context.Context, autho *model.AuthoUser, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) ([]*branch.PublicBranch, error) {
	if autho.Role != model.AccountRoleClient {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	lqID := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}

	lqID.Filter.And("user_id = ?", autho.ID)

	branchIDs := []int{}
	if err := s.savedBranchDB.List(s.db.Table("saved_branch").Select("branch_id"), &branchIDs, lqID, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting saved branch list: %v", err))
		return nil, server.NewHTTPInternalError("Error when getting saved branch list")
	}

	lq.Filter.And("id IN (?)", branchIDs)

	resp, err := s.branchUseCase.ListPublicBranches(ctx, lq, lqBranch)
	if err != nil {
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when listing saved branches: %v", err))
	}

	return resp.Data, nil
}

func (s *Branch) Rate(ctx context.Context, autho *model.AuthoUser, branchID int, req RatingBranchRequest) error {
	if autho.Role != model.AccountRoleClient {
		return server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		branch := &model.Branch{}
		if err := s.branchDB.View(tx.Preload("Rooms"), &branch, `id = ? AND is_active = true`, branchID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking branch exist: %v", err))
			return server.NewHTTPInternalError("Branch not found")
		}

		existedRoom := false
		for _, room := range branch.Rooms {
			if room.ID == req.RoomID {
				existedRoom = true
				break
			}
		}

		if !existedRoom {
			return server.NewHTTPValidationError("Room not found")
		}

		existed, err := s.bookingDB.Exist(tx, `id = ? AND user_id = ? AND room_id = ?`, req.BookingID, autho.ID, req.RoomID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when checking booking exist: %v", err))
			return server.NewHTTPInternalError("Error when checking booking exist")
		}

		if !existed {
			return server.NewHTTPValidationError("Booking not found")
		}

		creation := &model.Rating{
			BranchID: branchID,
			UserID:   autho.ID,
			Rating:   req.Rating,
			Comment:  req.Comment,
		}

		if err := s.ratingDB.Create(tx, &creation); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when creating rating: %v", err))
			return server.NewHTTPInternalError("Error when creating rating")
		}

		if err := s.bookingDB.Update(tx, map[string]interface{}{
			"status": model.BookingStatusCompleted,
		}, `id = ?`, req.BookingID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when updating booking status: %v", err))
			return server.NewHTTPInternalError("Error when updating booking status")
		}

		return nil
	})
}
