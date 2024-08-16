package bank

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/rbac"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	bankDB BankDB,
	companyDB CompanyDB,
	rbac rbac.Intf,
) *BankService {
	return &BankService{
		db:        db,
		bankDB:    bankDB,
		companyDB: companyDB,
		rbac:      rbac,
	}
}

type BankService struct {
	db        *gorm.DB
	bankDB    BankDB
	companyDB CompanyDB
	rbac      rbac.Intf
}

type BankDB interface {
	dbcore.Intf
}

type CompanyDB interface {
	dbcore.Intf
	CheckExistedBranchOfRepresentative(db *gorm.DB, rpID, branchID int) (bool, error)
	CheckBankAccountOwnership(db *gorm.DB, rpID, bankID int) (bool, error)
}
