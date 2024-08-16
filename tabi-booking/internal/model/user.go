package model

import "time"

// User represents the user model
// swagger:model
type User struct {
	ID          int        `json:"id" gorm:"primary_key"`
	FirstName   string     `json:"first_name" gorm:"type:text"`
	LastName    string     `json:"last_name" gorm:"type:text"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	AccountID   int        `json:"account_id"`
	Account     *Account   `gorm:"foreignKey:AccountID"`
	Booking     []*Booking `gorm:"foreignKey:UserID"`
	Survey      *Survey    `gorm:"foreignKey:UserID"`
	Base
}

// swagger:model UserResponse
type UserResponse struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	DoB       time.Time `json:"dob"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
}

func (u *User) ToResponse() *UserResponse {
	resp := &UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		DoB:       *u.DateOfBirth,
	}

	if u.Account != nil {
		resp.Username = u.Account.Username
		resp.Email = u.Account.Email
		resp.Phone = u.Account.Phone
	}

	return resp
}
