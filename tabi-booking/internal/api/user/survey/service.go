package survey

import (
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func New(db *gorm.DB,
	surveyDB SurveyDB,
) *Survey {
	return &Survey{
		db:       db,
		surveyDB: surveyDB,
	}
}

type Survey struct {
	db       *gorm.DB
	surveyDB SurveyDB
}

type SurveyDB interface {
	dbcore.Intf
}
