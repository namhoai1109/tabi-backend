package roomtype

import (
	"context"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/rbac"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	roomTypeDB RoomTypeDB,
	roomTypeOfBranchDB RoomTypeOfBranchDB,
	branchDB BranchDB,
	companyDB CompanyDB,
	facilityDB FacilityDB,
	rbac rbac.Intf,
) *RoomType {
	return &RoomType{
		db:                 db,
		roomTypeDB:         roomTypeDB,
		roomTypeOfBranchDB: roomTypeOfBranchDB,
		branchDB:           branchDB,
		companyDB:          companyDB,
		facilityDB:         facilityDB,
		rbac:               rbac,
	}
}

type RoomType struct {
	db                 *gorm.DB
	roomTypeDB         RoomTypeDB
	roomTypeOfBranchDB RoomTypeOfBranchDB
	branchDB           BranchDB
	companyDB          CompanyDB
	facilityDB         FacilityDB
	rbac               rbac.Intf
}

type RoomTypeDB interface {
	dbcore.Intf
}

type RoomTypeOfBranchDB interface {
	dbcore.Intf
	GetRoomTypes(db *gorm.DB, ctx context.Context, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomType, error)
}

type BranchDB interface {
	dbcore.Intf
}

type FacilityDB interface {
	dbcore.Intf
	GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error)
}

type CompanyDB interface {
	dbcore.Intf
	CheckExistedBranchOfRepresentative(db *gorm.DB, rpID, branchID int) (bool, error)
	GetBranchIDsOfRP(db *gorm.DB, rpID int) ([]int, error)
	GetRoomTypeIDsOfRP(db *gorm.DB, rpID int) ([]int, error)
}
