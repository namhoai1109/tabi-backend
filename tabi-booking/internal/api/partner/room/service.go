package room

import (
	"context"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/rbac"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	roomDB RoomDB,
	roomTypeOfBranchDB RoomTypeOfBranchDB,
	branchDB BranchDB,
	generalTypeDB GeneralTypeDB,
	factureReductionDB FactureReductionDB,
	reservationReductionDB ReservationReductionDB,
	companyDB CompanyDB,
	facilityDB FacilityDB,
	bookingDB BookingDB,
	userDB UserDB,
	rbac rbac.Intf,
) *Room {
	return &Room{
		db:                     db,
		roomDB:                 roomDB,
		roomTypeOfBranchDB:     roomTypeOfBranchDB,
		branchDB:               branchDB,
		generalTypeDB:          generalTypeDB,
		factureReductionDB:     factureReductionDB,
		reservationReductionDB: reservationReductionDB,
		companyDB:              companyDB,
		facilityDB:             facilityDB,
		bookingDB:              bookingDB,
		userDB:                 userDB,
		rbac:                   rbac,
	}
}

type Room struct {
	db                     *gorm.DB
	roomDB                 RoomDB
	roomTypeOfBranchDB     RoomTypeOfBranchDB
	branchDB               BranchDB
	generalTypeDB          GeneralTypeDB
	factureReductionDB     FactureReductionDB
	reservationReductionDB ReservationReductionDB
	companyDB              CompanyDB
	facilityDB             FacilityDB
	bookingDB              BookingDB
	userDB                 UserDB
	rbac                   rbac.Intf
}

type CompanyDB interface {
	dbcore.Intf
}

type RoomDB interface {
	dbcore.Intf
}

type RoomTypeOfBranchDB interface {
	dbcore.Intf
	GetRoomTypeIDsOfBranch(db *gorm.DB, ctx context.Context, lq *dbcore.ListQueryCondition, branchID int) ([]int, error)
}

type BranchDB interface {
	dbcore.Intf
	GetBranchIDByBranchManagerID(db *gorm.DB, managerID int) (int, error)
	GetBranchIDsByRepresentativeID(db *gorm.DB, companyID int) ([]int, error)
}

type GeneralTypeDB interface {
	dbcore.Intf
}

type FactureReductionDB interface {
	dbcore.Intf
}

type ReservationReductionDB interface {
	dbcore.Intf
}

type FacilityDB interface {
	dbcore.Intf
	GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error)
}

type BookingDB interface {
	dbcore.Intf
}

type UserDB interface {
	dbcore.Intf
}
