package model

// Device represents the device model
// swagger:model
type Device struct {
	ID        int    `json:"id" gorm:"primary_key"`
	UserID    int    `json:"user_id"`
	Brand     string `json:"brand" gorm:"type:text"`
	Model     string `json:"model" gorm:"type:text"`
	OS        string `json:"os" gorm:"type:text"`
	OSVersion string `json:"os_version" gorm:"type:text"`
	PushToken string `json:"push_token" gorm:"type:text"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	Base
}
