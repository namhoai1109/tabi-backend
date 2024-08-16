package facility

import (
	dbcore "github.com/namhoai1109/tabi/core/db"

	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	facilityDB FacilityDB,
) *Facility {
	return &Facility{
		db:         db,
		facilityDB: facilityDB,
	}
}

type Facility struct {
	db         *gorm.DB
	facilityDB FacilityDB
}

type FacilityDB interface {
	dbcore.Intf
}
