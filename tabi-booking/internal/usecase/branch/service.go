package branch

import (
	"context"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

type BranchUseCase interface {
	ListPublicBranches(ctx context.Context, lq *dbcore.ListQueryCondition, bc *PublicBranchCondition) (*PublicBranchListResponse, error)
	CreateBranch(db *gorm.DB, ctx context.Context, branchCreation BranchCreationRequest, companyID *int) (*model.Branch, error)
}

func New(db *gorm.DB,
	branchDB BranchDB,
	roomDB RoomDB,
	facilityDB FacilityDB,
	generalTypeDB GeneralTypeDB,
	bankDB BankDB,
) BranchUseCase {
	return &Service{
		db:            db,
		branchDB:      branchDB,
		roomDB:        roomDB,
		facilityDB:    facilityDB,
		generalTypeDB: generalTypeDB,
		bankDB:        bankDB,
	}
}

type Service struct {
	db            *gorm.DB
	branchDB      BranchDB
	roomDB        RoomDB
	facilityDB    FacilityDB
	generalTypeDB GeneralTypeDB
	bankDB        BankDB
}

type BranchDB interface {
	dbcore.Intf
}

type RoomDB interface {
	dbcore.Intf
}

type FacilityDB interface {
	dbcore.Intf
	GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error)
}

type GeneralTypeDB interface {
	dbcore.Intf
}

type BankDB interface {
	dbcore.Intf
}
