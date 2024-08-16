package branch

import (
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(db *gorm.DB,
	branchDB BranchDB,
	roomDB RoomDB,
	facilityDB FacilityDB,
	branchUseCase branch.BranchUseCase,
	bookingDB BookingDB,
	ratingDB RatingDB,
) *Branch {
	return &Branch{
		db:            db,
		branchDB:      branchDB,
		roomDB:        roomDB,
		facilityDB:    facilityDB,
		branchUseCase: branchUseCase,
		bookingDB:     bookingDB,
		ratingDB:      ratingDB,
	}
}

type Branch struct {
	db            *gorm.DB
	branchDB      BranchDB
	roomDB        RoomDB
	facilityDB    FacilityDB
	branchUseCase branch.BranchUseCase
	bookingDB     BookingDB
	ratingDB      RatingDB
}

type BranchDB interface {
	dbcore.Intf
	ListFeaturedDestination(db *gorm.DB) ([]string, error)
	ListFeaturedBranches(db *gorm.DB) ([]*model.Branch, error)
}

type RoomDB interface {
	dbcore.Intf
}

type FacilityDB interface {
	dbcore.Intf
	GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error)
}

type BookingDB interface {
	dbcore.Intf
}

type RatingDB interface {
	dbcore.Intf
}
