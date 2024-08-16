package model

// swagger:model GeneralType
type GeneralType struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Class   string `json:"class" gorm:"type:varchar(128)"`
	LabelEN string `json:"label_en" gorm:"type:text"`
	LabelVI string `json:"label_vi" gorm:"type:text"`
	DescEN  string `json:"desc_en" gorm:"type:text"`
	DescVI  string `json:"desc_vi" gorm:"type:text"`
	Order   int    `json:"order" gorm:"default:1"`
	Base
}

const (
	GeneralTypeClassHotelAccommodation     = "HOTEL_ACCOMMODATION"
	GeneralTypeClassHouseAccommodation     = "HOUSE_ACCOMMODATION"
	GeneralTypeClassApartmentAccommodation = "APARTMENT_ACCOMMODATION"
	GeneralTypeClassUniqueAccommodation    = "UNIQUE_ACCOMMODATION"
	GeneralTypeClassBed                    = "BED"
	GeneralTypeLangEN                      = "en"
	GeneralTypeLangVI                      = "vi"
)

var GeneralTypeAccommodationClass = []string{
	GeneralTypeClassHotelAccommodation,
	GeneralTypeClassHouseAccommodation,
	GeneralTypeClassApartmentAccommodation,
	GeneralTypeClassUniqueAccommodation,
}
