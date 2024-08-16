package partner

import "github.com/lib/pq"

// swagger:model RpRegistrationReq
type RpRegistrationReq struct {
	// account info
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email,omitempty"`
	// representative info
	FullName string `json:"full_name" validate:"required"`
	// company info
	CompanyName string `json:"company_name" validate:"required"`
	ShortName   string `json:"short_name" validate:"required"`
	Description string `json:"description" validate:"required"`
	WebsiteURL  string `json:"website_url" validate:"required"`
	TaxNumber   string `json:"tax_number" validate:"required"`
}

//swagger:model HstRegistrationReq
type HstRegistrationReq struct {
	// account info
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email,omitempty"`
	// host info
	FullName string `json:"full_name" validate:"required"`
	//branch info
	BranchName            string        `json:"branch_name" validate:"required"`
	Address               string        `json:"address" validate:"required"`
	ProvinceCity          string        `json:"province_city" validate:"required"`
	District              string        `json:"district" validate:"required"`
	Ward                  string        `json:"ward" validate:"required"`
	Latitude              string        `json:"latitude,omitempty"`
	Longitude             string        `json:"longitude,omitempty"`
	Description           string        `json:"description,omitempty"`
	ReceptionArea         bool          `json:"reception_area"`
	MainFacilities        pq.Int64Array `json:"main_facilities" validate:"required"`
	TypeID                int           `json:"type_id" validate:"required"`
	CancellationTimeUnit  string        `json:"cancellation_time_unit" validate:"required,min=1"`
	CancellationTimeValue int           `json:"cancellation_time_value" validate:"required"`
	GeneralPolicy         string        `json:"general_policy" validate:"required"`
	WebsiteURL            string        `json:"website_url" validate:"required"`
	TaxNumber             string        `json:"tax_number" validate:"required"`
}

// swagger:model CredentialsPartnerReq
type CredentialsPartnerReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenReq represents refresh token request data
// swagger:model
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

//swagger:model HSTRegisterResponse
type HSTRegisterResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	BranchID     int    `json:"branch_id,omitempty"`
}
