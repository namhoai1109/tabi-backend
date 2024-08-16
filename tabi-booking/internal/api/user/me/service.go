package me

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	userDB UserDB,
) *Me {
	return &Me{
		db:     db,
		userDB: userDB,
	}
}

type Me struct {
	db     *gorm.DB
	userDB UserDB
}

type UserDB interface {
	dbcore.Intf
}
