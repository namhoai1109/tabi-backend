package schedule

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	scheduleDB ScheduleDB,
) *Schedule {
	return &Schedule{
		db:         db,
		scheduleDB: scheduleDB,
	}
}

type Schedule struct {
	db         *gorm.DB
	scheduleDB ScheduleDB
}

type ScheduleDB interface {
	dbcore.Intf
}
