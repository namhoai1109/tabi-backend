package model

// swagger:model BranchManager
type BranchManager struct {
	ID               int             `json:"id" gorm:"primaryKey"`
	Name             string          `json:"name" gorm:"type:varchar(32)"`
	AccountID        int             `json:"account_id"`
	RepresentativeID int             `json:"representative_id" gorm:"default:NULL"`
	Account          *Account        `gorm:"foreignKey:AccountID"`
	Representative   *Representative `gorm:"foreignKey:RepresentativeID"`
	Base
}

// swagger:model BranchManagerResponse
type BranchManagerResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func (b *BranchManager) ToResponse() *BranchManagerResponse {
	resp := &BranchManagerResponse{
		ID:   b.ID,
		Name: b.Name,
	}

	if b.Account != nil {
		resp.Username = b.Account.Username
		resp.Phone = b.Account.Phone
		resp.Email = b.Account.Email
		resp.Role = b.Account.Role
	}

	return resp
}
