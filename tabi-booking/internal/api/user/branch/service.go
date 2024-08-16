package branch

import (
	"tabi-booking/internal/usecase/branch"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	branchDB BranchDB,
	savedBranchDB SavedBranchDB,
	branchUseCase branch.BranchUseCase,
	ratingDB RatingDB,
	bookingDB BookingDB,
) *Branch {
	return &Branch{
		db:            db,
		branchDB:      branchDB,
		savedBranchDB: savedBranchDB,
		branchUseCase: branchUseCase,
		ratingDB:      ratingDB,
		bookingDB:     bookingDB,
	}
}

type Branch struct {
	db            *gorm.DB
	branchDB      BranchDB
	savedBranchDB SavedBranchDB
	branchUseCase branch.BranchUseCase
	ratingDB      RatingDB
	bookingDB     BookingDB
}

type BranchDB interface {
	dbcore.Intf
}

type SavedBranchDB interface {
	dbcore.Intf
}

type RatingDB interface {
	dbcore.Intf
}

type BookingDB interface {
	dbcore.Intf
}
