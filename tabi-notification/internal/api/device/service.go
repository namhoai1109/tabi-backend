package device

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	deviceDB DeviceDB,
) *Device {
	return &Device{
		db:       db,
		deviceDB: deviceDB,
	}
}

type Device struct {
	db       *gorm.DB
	deviceDB DeviceDB
}

type DeviceDB interface {
	dbcore.Intf
}
