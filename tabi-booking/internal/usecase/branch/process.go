package branch

import (
	"context"
	"math"
	"sort"
	"tabi-booking/internal/model"
	"time"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

var (
	MAX_FLOAT64 float64 = 10000000000
)

func (s *Service) getPricesForBranchResponse(rooms []*model.Room, bc *PublicBranchCondition) (float64, float64) {
	now := time.Now()
	checkInDate := now
	checkOutDate := checkInDate.Add(time.Hour * 24)
	var (
		minPrice float64 = MAX_FLOAT64
		maxPrice float64 = 0
	)

	for _, room := range rooms {
		validReservation := room.ReservationReduction != nil && len(room.ReservationReduction) > 0
		validRoom := room.Status == model.RoomStatusUpdated && room.FactureReduction != nil && validReservation

		if bc != nil {
			if bc.Occupancy != nil {
				validRoom = validRoom && room.MaxOccupancy >= *bc.Occupancy
			}

			if bc.BookingDateIn != nil && len(bc.BookingDateIn) == 2 {
				checkInDate = bc.BookingDateIn[0]
				checkOutDate = bc.BookingDateIn[1]
			}
		}

		if validRoom {
			price := room.GetPriceForBookingDates(checkInDate, checkOutDate)
			if price < minPrice {
				minPrice = price
				maxPrice = room.MaxPrice
			}
		}
	}

	return minPrice, maxPrice
}

func (s *Service) getBranchResponses(branches []*model.Branch, bc *PublicBranchCondition) ([]*PublicBranch, error) {
	branchResponses := make([]*PublicBranch, 0, len(branches))

	for _, branch := range branches {
		branchResponse := &PublicBranch{
			ID:           branch.ID,
			Name:         branch.BranchName,
			ProvinceCity: branch.ProvinceCity,
			District:     branch.District,
		}

		minPrice, maxPrice := s.getPricesForBranchResponse(branch.Rooms, bc)
		if minPrice != MAX_FLOAT64 && maxPrice != 0 {
			branchResponse.MinPrice = math.Floor(minPrice)
			branchResponse.MaxPrice = math.Floor(maxPrice)
			branchResponses = append(branchResponses, branchResponse)
		}
	}

	sort.Slice(branchResponses, func(i, j int) bool {
		return branchResponses[i].MinPrice < branchResponses[j].MinPrice
	})

	return branchResponses, nil
}

func (s *Service) checkFacilityExist(db *gorm.DB, ctx context.Context, facilityIDs []int64) error {
	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}
	lq.Filter.And("id IN (?)", facilityIDs)
	list := []*model.Facility{}
	if err := s.facilityDB.List(db, &list, lq, nil); err != nil {
		logger.LogError(ctx, "Failed to get facility list")
		return server.NewHTTPInternalError("Failed to get facility list")
	}

	if len(list) != len(facilityIDs) {
		logger.LogError(ctx, "Some facilities are not exist")
		return server.NewHTTPValidationError("Some facilities are not exist")
	}

	return nil
}

func (s *Service) checkTypeExist(db *gorm.DB, ctx context.Context, typeID int) error {
	typeExist, err := s.generalTypeDB.Exist(db, `id = ?`, typeID)
	if err != nil {
		logger.LogError(ctx, "Failed to check type exist")
		return server.NewHTTPInternalError("Failed to check type exist")
	}

	if !typeExist {
		logger.LogError(ctx, "Type does not exist")
		return server.NewHTTPValidationError("Type does not exist")
	}

	return nil
}
