package model

// swagger:model Representative
type Representative struct {
	ID        int      `json:"id" gorm:"primary_key"`
	Name      string   `json:"name" gorm:"type:varchar(32)"`
	AccountID int      `json:"account_id"`
	Account   *Account `gorm:"foreignKey:AccountID"`
	Base
}

type RepresentativeResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func (r *Representative) ToRepresentativeResponse() *RepresentativeResponse {
	resp := &RepresentativeResponse{
		ID:   r.ID,
		Name: r.Name,
	}

	if r.Account != nil {
		resp.Username = r.Account.Username
		resp.Email = r.Account.Email
		resp.Phone = r.Account.Phone
		resp.Role = r.Account.Role
	}

	return resp
}
