package branch

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"
	structutil "tabi-booking/internal/util/struct"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"github.com/thoas/go-funk"

	"gorm.io/gorm"
)

func (s *Branch) View(ctx context.Context, authoPartner *model.AuthoPartner, branchID int) (*model.BranchResponse, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	facilities := []*model.Facility{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = s.checkExistBranch(tx, ctx, authoPartner, branchID); err != nil {
			return err
		}

		query := tx.Preload("Banks").Preload("Type").Preload("BranchManager").Preload("BranchManager.Account").
			Preload("Ratings", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at DESC")
			}).Preload("Ratings.User").Preload("Ratings.User.Account")
		if err = s.branchDB.View(query, &branch, `id = ?`, branchID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to view branch: %v", err))
		}

		ids := []int64{}
		ids = append(ids, branch.MainFacilities...)
		facilities, err = s.facilityDB.GetFacilityList(tx, ids)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get facility list: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to get facility list: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return branch.ToBranchResponse(facilities), nil
}

func (s *Branch) List(
	ctx context.Context,
	authoPartner *model.AuthoPartner,
	lq *dbcore.ListQueryCondition) (*BranchListResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}

	var count int64 = 0
	branches := []*model.Branch{}

	trxErr := s.db.Transaction(func(db *gorm.DB) error {
		company, err := s.companyDB.GetCompanyByRepID(db, authoPartner.ID)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get company while listing branches: %v", err)
			logger.LogError(ctx, errMsg)
			return server.NewHTTPInternalError(errMsg)
		}

		lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
			"branch_name": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(branch_name) like ?`, val}
			},
			"address": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(address) like ?`, val}
			},
			"province_city": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(province_city) like ?`, val}
			},
			"district": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(district) like ?`, val}
			},
			"ward": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
				val = "%" + strings.ToLower(val.(string)) + "%"
				return []interface{}{
					`lower(ward) like ?`, val}
			},
		})

		lq.Filter.And("company_id = ?", company.ID)
		tx := db.Preload("Type")
		if err := s.branchDB.List(tx, &branches, lq, &count); err != nil {
			errMsg := fmt.Sprintf("Failed to list branches: %v", err)
			logger.LogError(ctx, errMsg)
			return server.NewHTTPInternalError(errMsg)
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.toBranchListResponse(branches, count), nil
}

func (s *Branch) CreateBranch(ctx context.Context, authoPartner *model.AuthoPartner, branchCreation *branch.BranchCreationRequest) (*model.BranchResponse, error) {
	if err := s.enforce(authoPartner, model.ActionCreate); err != nil {
		return nil, err
	}

	if !funk.Contains(model.BranchCancellationTimeUnits, branchCreation.CancellationTimeUnit) {
		return nil, server.NewHTTPValidationError("Cancellation time unit is invalid")
	}

	branch := &model.Branch{}
	trxErr := s.db.Transaction(func(db *gorm.DB) error {
		company, err := s.companyDB.GetCompanyByRepID(db, authoPartner.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get company: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to get company: %v", err))
		}

		branch, err = s.branchUseCase.CreateBranch(db, ctx, *branchCreation, &company.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to create branch: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to create branch: %v", err))
		}
		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.View(ctx, authoPartner, branch.ID)
}

func (s *Branch) Update(ctx context.Context,
	authoPartner *model.AuthoPartner,
	branchID int,
	data *BranchUpdateRequest) (*model.BranchResponse, error) {
	if err := s.enforce(authoPartner, model.ActionUpdateAll); err != nil {
		return nil, err
	}

	if reflect.ValueOf(*data).IsZero() {
		return nil, server.NewHTTPValidationError("Data is empty")
	}

	updates := structutil.ToMap(*data)
	if updates["cancellation_time_unit"] != nil && !funk.Contains(model.BranchCancellationTimeUnits, updates["cancellation_time_unit"]) {
		return nil, server.NewHTTPValidationError("Cancellation time unit is invalid")
	}

	if err := s.checkExistBranch(s.db, ctx, authoPartner, branchID); err != nil {
		return nil, err
	}

	trxErr := s.db.Transaction(func(db *gorm.DB) error {
		if updates["email"] != nil && authoPartner.Role == model.RoleHost {
			bm := &model.BranchManager{}
			if err := s.branchManagerDB.View(db.Preload("Account"), bm, `id = ?`, authoPartner.ID); err != nil {
				errMsg := fmt.Sprintf("Failed to get branch manager: %v", err)
				logger.LogError(ctx, errMsg)
				return server.NewHTTPInternalError(errMsg)
			}

			if err := s.accountDB.Update(db, map[string]interface{}{
				"email": updates["email"],
			}, `id = ?`, bm.Account.ID); err != nil {
				errMsg := fmt.Sprintf("Failed to update account: %v", err)
				logger.LogError(ctx, errMsg)
				return server.NewHTTPInternalError(errMsg)
			}
		}
		delete(updates, "email")

		if err := s.branchDB.Update(db, updates, branchID); err != nil {
			errMsg := fmt.Sprintf("Failed to update branch: %v", err)
			logger.LogError(ctx, errMsg)
			return server.NewHTTPInternalError(errMsg)
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.View(ctx, authoPartner, branchID)
}

func (s *Branch) Activate(ctx context.Context, authoPartner *model.AuthoPartner) (*ActivateBranchResponse, error) {
	if err := s.enforce(authoPartner, model.ActionUpdate); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	preloadDB := s.db.Preload("Rooms").
		Preload("Rooms.FactureReduction").
		Preload("Rooms.ReservationReduction")
	if err := s.branchDB.View(preloadDB, &branch, `branch_manager_id = ? AND is_active = false`, authoPartner.ID); err != nil {
		errMsg := fmt.Sprintf("Branch is not found or branch is already activated: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	haveValidRoom := false
	for _, room := range branch.Rooms {
		validReservation := room.ReservationReduction != nil && len(room.ReservationReduction) > 0
		validRoom := room.Status == model.RoomStatusUpdated && room.FactureReduction != nil && validReservation
		if validRoom {
			haveValidRoom = true
			break
		}
	}

	if !haveValidRoom {
		return nil, server.NewHTTPValidationError("Branch does not have any valid room")
	}

	if err := s.branchDB.Update(s.db, map[string]interface{}{"is_active": true}, branch.ID); err != nil {
		errMsg := fmt.Sprintf("Failed to activate branch: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return &ActivateBranchResponse{
		Activated: true,
	}, nil
}

func (s *Branch) AnalyzeRevenues(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.RevenueAnalysisData, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
		errMsg := fmt.Sprintf("Failed to get branch: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	resp, err := s.branchDB.AnalyzeBranchesRevenue(s.db, []int{branch.ID}, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to analyze branches revenue: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return resp, nil
}

func (s *Branch) AnalyzeBookingRequestQuantity(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.BookingRequestQuantityAnalysisData, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	if err := s.branchDB.View(s.db, branch, `branch_manager_id = ?`, authoPartner.ID); err != nil {
		errMsg := fmt.Sprintf("Failed to get branch: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	resp, err := s.branchDB.AnalyzeBookingRequestQuantity(s.db, []int{branch.ID}, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to analyze branches revenue: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return resp, nil
}
