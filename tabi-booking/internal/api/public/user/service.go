package user

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(db *gorm.DB,
	userDB UserDB,
	surveyDB SurveyDB,
) *User {
	return &User{
		db:       db,
		userDB:   userDB,
		surveyDB: surveyDB,
	}
}

type User struct {
	db       *gorm.DB
	userDB   UserDB
	surveyDB SurveyDB
}

type UserDB interface {
	dbcore.Intf
}

type SurveyDB interface {
	dbcore.Intf
}
