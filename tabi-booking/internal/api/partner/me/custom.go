package me

// swagger:model PartnerInfoResponse
type PartnerInfoResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	PlaceID  int    `json:"place_id"`
}
