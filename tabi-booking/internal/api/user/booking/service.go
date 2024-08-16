package booking

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	roomDB RoomDB,
	bookingDB BookingDB,
) *Booking {
	return &Booking{
		db:        db,
		roomDB:    roomDB,
		bookingDB: bookingDB,
	}
}

type Booking struct {
	db        *gorm.DB
	roomDB    RoomDB
	bookingDB BookingDB
}

type RoomDB interface {
	dbcore.Intf
}

type BookingDB interface {
	dbcore.Intf
}
