package branchmanager

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	accountDB AccountDB,
	branchManagerDB BranchManagerDB,
	branchDB BranchDB,
	companyDB CompanyDB,
) *BranchManagerService {
	return &BranchManagerService{
		db:              db,
		accountDB:       accountDB,
		branchManagerDB: branchManagerDB,
		branchDB:        branchDB,
		companyDB:       companyDB,
	}
}

type BranchManagerService struct {
	db              *gorm.DB
	accountDB       AccountDB
	branchManagerDB BranchManagerDB
	branchDB        BranchDB
	companyDB       CompanyDB
}

type AccountDB interface {
	dbcore.Intf
}

type BranchManagerDB interface {
	dbcore.Intf
}

type BranchDB interface {
	dbcore.Intf
	CheckExitedBranchManager(db *gorm.DB, branchID int) (bool, error)
}

type CompanyDB interface {
	dbcore.Intf
	CheckExistedBranchOfRepresentative(db *gorm.DB, rpID, branchID int) (bool, error)
}
