package branch

import (
	"context"
	"fmt"
	"sort"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"
	"time"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *Branch) List(ctx context.Context, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) (*branch.PublicBranchListResponse, error) {
	return s.branchUseCase.ListPublicBranches(ctx, lq, lqBranch)
}

func (s *Branch) View(ctx context.Context, id int) (*model.PublicBranchResponse, error) {
	if err := s.checkExistBranch(ctx, id); err != nil {
		return nil, err
	}

	branch := &model.Branch{}
	preloadDB := s.db.Preload("Type").
		Preload("Ratings", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).Preload("Ratings.User").Preload("Ratings.User.Account").
		Preload("BranchManager").Preload("BranchManager.Account").
		Preload("Company").Preload("Company.Representative").Preload("Company.Representative.Account")
	if err := s.branchDB.View(preloadDB, &branch, id); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when view branch: %v", err))
		return nil, server.NewHTTPInternalError("Error when view branch")
	}

	ids := []int64{}
	ids = append(ids, branch.MainFacilities...)
	facilities, err := s.facilityDB.GetFacilityList(s.db, ids)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to get facility list: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get facility list: %v", err))
	}

	return branch.ToPublicBranchResponse(facilities), nil
}

func (s *Branch) ListRooms(ctx context.Context, branchID int, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) (*RoomListResponse, error) {
	if err := s.checkExistBranch(ctx, branchID); err != nil {
		return nil, err
	}

	rooms := []*model.Room{}
	count := int64(0)
	lq.Filter.And(`branch_id = ? AND status = ?`, branchID, model.RoomStatusUpdated)
	preloadBD := s.db.Preload("RoomType").
		Preload("BedType").
		Preload("FactureReduction").
		Preload("ReservationReduction").
		Preload("Bookings", "status = ? OR status = ?", model.BookingStatusApproved, model.BookingStatusPending)
	if err := s.roomDB.List(preloadBD, &rooms, lq, &count); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when list room: %v", err))
		return nil, server.NewHTTPInternalError("Error when list room")
	}

	resp := []*model.PublicRoom{}
	for _, room := range rooms {
		facilities, err := s.facilityDB.GetFacilityList(s.db, room.RoomType.RoomFacilities)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get facility list: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get facility list: %v", err))
		}

		checkInDate := time.Now()
		checkOutDate := checkInDate.Add(time.Hour * 24)

		if lqBranch != nil && lqBranch.BookingDateIn != nil && len(lqBranch.BookingDateIn) == 2 {
			checkInDate = lqBranch.BookingDateIn[0]
			checkOutDate = lqBranch.BookingDateIn[1]
		}

		resp = append(resp, room.ToPublicRoomResponse(facilities, checkInDate, checkOutDate))
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].CurrentPrice < resp[j].CurrentPrice
	})

	return &RoomListResponse{
		Total: count,
		Data:  resp,
	}, nil
}

func (s *Branch) ListFeaturedDestinations(ctx context.Context) ([]string, error) {
	return s.branchDB.ListFeaturedDestination(s.db)
}

func (s *Branch) ListFeaturedBranches(ctx context.Context) ([]*branch.PublicBranch, error) {
	branches, err := s.branchDB.ListFeaturedBranches(s.db)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when list featured branches: %v", err))
		return nil, server.NewHTTPInternalError("Error when list featured branches")
	}

	branchIDs := []int{}
	for _, branch := range branches {
		branchIDs = append(branchIDs, branch.ID)
	}

	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}
	lq.Filter.And("id IN (?)", branchIDs)
	res, err := s.branchUseCase.ListPublicBranches(ctx, lq, nil)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to get public branches: %s", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get public branches: %s", err))
	}

	for index := range res.Data {
		ratings, err := s.getRatings(ctx, res.Data[index].ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get ratings: %s", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get ratings: %s", err))
		}

		if len(ratings) == 0 {
			res.Data[index].StarLevel = 0
			res.Data[index].ReviewQuantity = 0
			continue
		}

		sum := 0
		for _, rating := range ratings {
			sum += rating.Rating
		}

		res.Data[index].StarLevel = float64(sum) / float64(len(ratings))
		res.Data[index].ReviewQuantity = len(ratings)
	}

	return res.Data, nil
}

func (s *Branch) ListRecommendedBranches(ctx context.Context, lq *dbcore.ListQueryCondition, userID int) (*branch.PublicBranchListResponse, error) {
	destination := []string{}

	if userID == -1 || userID == 0 {
		branches, err := s.branchDB.ListFeaturedBranches(s.db)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when list featured branches: %v", err))
			return nil, server.NewHTTPInternalError("Error when list featured branches")
		}

		for _, branch := range branches {
			destination = append(destination, branch.ProvinceCity)
			destination = append(destination, branch.District)
		}
	} else {
		lqUser := &dbcore.ListQueryCondition{
			Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
		}
		lqUser.Filter.And("user_id = ?", userID)
		bookings := []*model.Booking{}
		if err := s.bookingDB.List(s.db.Preload("Room").Preload("Room.Branch"), &bookings, lqUser, nil); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get booking list: %s", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get booking list: %s", err))
		}

		for _, booking := range bookings {
			if booking.Room != nil && booking.Room.Branch != nil {
				destination = append(destination, booking.Room.Branch.ProvinceCity)
				destination = append(destination, booking.Room.Branch.District)
			}
		}
	}

	if len(destination) != 0 {
		lq.Filter.And("province_city IN (?) OR district IN (?)", destination, destination)
	}

	return s.branchUseCase.ListPublicBranches(ctx, lq, nil)
}
