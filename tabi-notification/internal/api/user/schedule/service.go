package schedule

import (
	"tabi-notification/config"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(db *gorm.DB, scheduleDB ScheduleDB, cfg *config.Configuration) *Schedule {
	return &Schedule{db: db, scheduleDB: scheduleDB, cfg: cfg}
}

type Schedule struct {
	db         *gorm.DB
	scheduleDB ScheduleDB
	cfg        *config.Configuration
}

type ScheduleDB interface {
	dbcore.Intf
}
