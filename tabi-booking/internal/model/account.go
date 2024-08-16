package model

import "time"

// swagger:model Account
type Account struct {
	ID           int        `json:"id" gorm:"primary_key"`
	Username     string     `json:"username" gorm:"type:text"`
	Password     string     `json:"password" gorm:"type:text"`
	Phone        string     `json:"phone" gorm:"type:varchar(11)"`
	Email        string     `json:"email" gorm:"type:varchar(128)"`
	Role         string     `json:"role" gorm:"type:varchar(3)"`
	RefreshToken string     `json:"refresh_token" gorm:"type:varchar(255);unique_index"`
	LastLogin    *time.Time `json:"last_login"`
	Base
}

const (
	AccountRoleRepresentative = "REP"
	AccountRoleBranchManager  = "BMA"
	AccountRoleHost           = "HST"
	AccountRoleClient         = "CLI"
)

var AccountRolesForPartner = []string{
	AccountRoleRepresentative,
	AccountRoleBranchManager,
	AccountRoleHost,
}
