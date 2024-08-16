package user

import (
	"tabi-booking/config"
	"time"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	accountDB AccountDB,
	userDB UserDB,
	jwt JWT,
	cfg *config.Configuration,
) *AuthenUser {
	return &AuthenUser{
		db:        db,
		accountDB: accountDB,
		userDB:    userDB,
		jwt:       jwt,
		cfg:       cfg,
	}
}

type AuthenUser struct {
	db        *gorm.DB
	accountDB AccountDB
	userDB    UserDB
	jwt       JWT
	cfg       *config.Configuration
}

type AccountDB interface {
	dbcore.Intf
}

type UserDB interface {
	dbcore.Intf
}

type JWT interface {
	GenerateToken(claims map[string]interface{}, expire *time.Time) (string, int, error)
}
