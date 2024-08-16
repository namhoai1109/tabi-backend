package company

import (
	"tabi-booking/config"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"

	"gorm.io/gorm"
)

type Company struct {
	db               *gorm.DB
	accountDB        AccountDB
	representativeDB RepresentativeDB
	companyDB        CompanyDB
	branchDB         BranchDB
	cfg              *config.Configuration
}

func New(
	db *gorm.DB,
	accountDB AccountDB,
	representativeDB RepresentativeDB,
	companyDB CompanyDB,
	branchDB BranchDB,
	cfg *config.Configuration,
) *Company {
	return &Company{
		db:               db,
		accountDB:        accountDB,
		representativeDB: representativeDB,
		companyDB:        companyDB,
		branchDB:         branchDB,
		cfg:              cfg,
	}
}

type CompanyDB interface {
	dbcore.Intf
}

type AccountDB interface {
	dbcore.Intf
}

type RepresentativeDB interface {
	dbcore.Intf
}

type BranchDB interface {
	dbcore.Intf
	AnalyzeBranchesRevenue(db *gorm.DB, branchIDs []int, year int) ([]*model.RevenueAnalysisData, error)
	AnalyzeBookingRequestQuantity(db *gorm.DB, branchIDs []int, year int) ([]*model.BookingRequestQuantityAnalysisData, error)
}
