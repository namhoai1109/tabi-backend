package booking

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/rbac"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	bookingDB BookingDB,
	roomDB RoomDB,
	userDB UserDB,
	branchDB BranchDB,
	rbac rbac.Intf,
) *Booking {
	return &Booking{
		db:        db,
		bookingDB: bookingDB,
		roomDB:    roomDB,
		userDB:    userDB,
		branchDB:  branchDB,
		rbac:      rbac,
	}
}

type Booking struct {
	db        *gorm.DB
	bookingDB BookingDB
	roomDB    RoomDB
	userDB    UserDB
	branchDB  BranchDB
	rbac      rbac.Intf
}

type BookingDB interface {
	dbcore.Intf
}

type BranchDB interface {
	dbcore.Intf
}

type RoomDB interface {
	dbcore.Intf
}

type UserDB interface {
	dbcore.Intf
}
