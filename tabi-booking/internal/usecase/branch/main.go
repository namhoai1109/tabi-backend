package branch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"tabi-booking/internal/model"
	"time"

	"github.com/imdatngo/gowhere"
	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Service) ListPublicBranches(ctx context.Context, lq *dbcore.ListQueryCondition, bc *PublicBranchCondition) (*PublicBranchListResponse, error) {
	branches := []*model.Branch{}

	preloadDB := s.db
	if bc != nil && bc.Occupancy != nil {
		preloadDB = preloadDB.Preload("Rooms", `max_occupancy >= ?`, *bc.Occupancy)
	} else {
		preloadDB = preloadDB.Preload("Rooms")
	}
	preloadDB = preloadDB.Preload("Rooms.FactureReduction").
		Preload("Rooms.ReservationReduction").
		Preload("Rooms.Bookings", "status = ? OR status = ?", model.BookingStatusApproved, model.BookingStatusPending)

	lq.Filter.SetCustomConditions(map[string]gowhere.CustomConditionFn{
		"destination": func(_ string, val interface{}, _ *gowhere.Config) interface{} {
			val = strings.ToLower(val.(string))
			return []interface{}{
				`lower(province_city) = ?
				OR lower(district) = ?`, val, val}
		},
	})
	lq.Filter.And("is_active = true")

	count := int64(0)
	if err := s.branchDB.List(preloadDB, &branches, lq, &count); err != nil {
		messageErr := fmt.Sprintf("Error when list branch: %v", err)
		logger.LogError(ctx, messageErr)
		return nil, server.NewHTTPInternalError(messageErr)
	}

	brancheResponse, err := s.getBranchResponses(branches, bc)
	if err != nil {
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when get branch responses: %v", err))
	}

	return &PublicBranchListResponse{
		Data:  brancheResponse,
		Total: count,
	}, nil
}

func ReqListQuery(c echo.Context) (*PublicBranchCondition, error) {
	lq := &PublicBranchFilter{}
	if err := c.Bind(lq); err != nil {
		return nil, server.NewHTTPValidationError(err.Error())
	}

	resp := &PublicBranchCondition{}

	if lq.BookingDateIn != "" {
		dates := []time.Time{}
		err := json.Unmarshal([]byte(lq.BookingDateIn), &dates)
		if err != nil {
			return nil, server.NewHTTPValidationError("Invalid filter, expecting JSON of time range for booking_date__in")
		}

		if len(dates) != 2 {
			return nil, server.NewHTTPValidationError("Invalid filter, booking_date__in must have 2 dates")
		}

		resp.BookingDateIn = dates
	}

	if lq.Occupancy != 0 {
		resp.Occupancy = &lq.Occupancy
	}

	return resp, nil
}

func (s *Service) CreateBranch(db *gorm.DB, ctx context.Context, branchCreation BranchCreationRequest, companyID *int) (*model.Branch, error) {
	if err := s.checkFacilityExist(db, ctx, branchCreation.MainFacilities); err != nil {
		return nil, err
	}

	if err := s.checkTypeExist(db, ctx, branchCreation.TypeID); err != nil {
		return nil, err
	}

	branch := &model.Branch{
		BranchName:            branchCreation.BranchName,
		Address:               branchCreation.Address,
		FullAddress:           branchCreation.Address + ", " + branchCreation.Ward + ", " + branchCreation.District + ", " + branchCreation.ProvinceCity,
		Ward:                  branchCreation.Ward,
		District:              branchCreation.District,
		ProvinceCity:          branchCreation.ProvinceCity,
		Latitude:              branchCreation.Latitude,
		Longitude:             branchCreation.Longitude,
		Description:           branchCreation.Description,
		ReceptionArea:         branchCreation.ReceptionArea,
		MainFacilities:        branchCreation.MainFacilities,
		TypeID:                branchCreation.TypeID,
		CancellationTimeUnit:  branchCreation.CancellationTimeUnit,
		CancellationTimeValue: branchCreation.CancellationTimeValue,
		GeneralPolicy:         branchCreation.GeneralPolicy,
		TaxNumber:             branchCreation.TaxNumber,
		WebsiteURL:            branchCreation.WebsiteURL,
	}

	if companyID != nil {
		branch.CompanyID = *companyID
	}

	if branchCreation.BranchManagerID != nil {
		branch.BranchManagerID = *branchCreation.BranchManagerID
	}

	if err := s.branchDB.Create(db, &branch); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to create branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to create branch: %v", err))
	}

	return branch, nil
}
