package user

import "time"

// swagger:model RegistrationUserReq
type RegistrationUserReq struct {
	// account info
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required"`
	// users info
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	DoB       time.Time `json:"dob" validate:"required"`
}

// swagger:model CredentialsUserReq
type CredentialsUserReq struct {
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}
