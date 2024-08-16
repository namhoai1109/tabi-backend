package model

//swagger:model Survey
type Survey struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	UserID     int    `json:"user_id"`
	PlaceType  string `json:"place_type"`
	Activities string `json:"activities"`
	Seasons    string `json:"seasons"`

	User *User `gorm:"foreignKey:UserID"`
}
