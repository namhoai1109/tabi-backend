package generaltype

import (
	dbcore "github.com/namhoai1109/tabi/core/db"

	"gorm.io/gorm"
)

func New(db *gorm.DB, generalTypeDB GeneralTypeDB) *GeneralType {
	return &GeneralType{
		db:            db,
		generalTypeDB: generalTypeDB,
	}
}

type GeneralType struct {
	db            *gorm.DB
	generalTypeDB GeneralTypeDB
}

type GeneralTypeDB interface {
	dbcore.Intf
}
