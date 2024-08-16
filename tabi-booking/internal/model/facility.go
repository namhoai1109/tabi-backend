package model

// swagger:model Facility
type Facility struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Type    string `json:"type" gorm:"type:varchar(10)"`
	ClassEN string `json:"class_en" gorm:"type:varchar(128)"`
	ClassVI string `json:"class_vi" gorm:"type:varchar(128)"`
	NameEN  string `json:"name_en" gorm:"type:text"`
	NameVI  string `json:"name_vi" gorm:"type:text"`
	Base
}

const (
	FacilityTypeMain = "MAIN"
	FacilityTypeRoom = "ROOM"
	FacilityLangEN   = "en"
	FacilityLangVI   = "vi"
)
