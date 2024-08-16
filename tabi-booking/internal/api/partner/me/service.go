package me

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	representativeDB RepresentativeDB,
	branchManagerDB BranchManagerDB,
	companyDB CompanyDB,
	branchDB BranchDB,
) *Me {
	return &Me{
		db:               db,
		representativeDB: representativeDB,
		branchManagerDB:  branchManagerDB,
		companyDB:        companyDB,
		branchDB:         branchDB,
	}
}

type Me struct {
	db               *gorm.DB
	representativeDB RepresentativeDB
	branchManagerDB  BranchManagerDB
	companyDB        CompanyDB
	branchDB         BranchDB
}

type RepresentativeDB interface {
	dbcore.Intf
}

type BranchManagerDB interface {
	dbcore.Intf
}

type CompanyDB interface {
	dbcore.Intf
}

type BranchDB interface {
	dbcore.Intf
}
