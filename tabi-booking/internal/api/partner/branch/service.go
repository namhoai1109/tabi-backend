package branch

import (
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/rbac"

	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	accountDB AccountDB,
	branchManagerDB BranchManagerDB,
	companyDB CompanyDB,
	branchDB BranchDB,
	bankDB BankDB,
	facilityDB FacilityDB,
	generalTypeDB GeneralTypeDB,
	branchUseCase branch.BranchUseCase,
	rbac rbac.Intf,
) *Branch {
	return &Branch{
		db:              db,
		accountDB:       accountDB,
		branchManagerDB: branchManagerDB,
		companyDB:       companyDB,
		branchDB:        branchDB,
		bankDB:          bankDB,
		facilityDB:      facilityDB,
		generalTypeDB:   generalTypeDB,
		branchUseCase:   branchUseCase,
		rbac:            rbac,
	}
}

type Branch struct {
	db              *gorm.DB
	accountDB       AccountDB
	branchManagerDB BranchManagerDB
	companyDB       CompanyDB
	branchDB        BranchDB
	bankDB          BankDB
	facilityDB      FacilityDB
	generalTypeDB   GeneralTypeDB
	branchUseCase   branch.BranchUseCase
	rbac            rbac.Intf
}

type CompanyDB interface {
	dbcore.Intf
	GetCompanyByRepID(db *gorm.DB, rpID int) (*model.Company, error)
}

type BranchDB interface {
	dbcore.Intf
	AnalyzeBranchesRevenue(db *gorm.DB, branchIDs []int, year int) ([]*model.RevenueAnalysisData, error)
	AnalyzeBookingRequestQuantity(db *gorm.DB, branchIDs []int, year int) ([]*model.BookingRequestQuantityAnalysisData, error)
}

type BankDB interface {
	dbcore.Intf
}

type FacilityDB interface {
	dbcore.Intf
	GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error)
}

type GeneralTypeDB interface {
	dbcore.Intf
}

type AccountDB interface {
	dbcore.Intf
}

type BranchManagerDB interface {
	dbcore.Intf
}
