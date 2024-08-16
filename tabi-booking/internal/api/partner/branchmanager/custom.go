package branchmanager

// swagger:model BranchManagerRegisterRequest
type BranchManagerRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Name     string `json:"name" validate:"required"`
	BranchID int    `json:"branch_id" validate:"required"`
}

// swagger:model BranchManagerRegisterResponse
type BranchManagerRegisterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
}
