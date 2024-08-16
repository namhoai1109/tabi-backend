package bank

// swagger:model BankRegisterRequest
type BankRegisterRequest struct {
	BankID        int    `json:"bank_id" validate:"required"`
	BankBranch    string `json:"bank_branch" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
}

// swagger:model CustomBankResponse
type CustomBankResponse struct {
	ID            int    `json:"id"`
	BankID        int    `json:"bank_id"`
	BankBranch    string `json:"bank_branch"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BranchID      int    `json:"branch_id"`
}

// swagger:model BankListResponse
type BankListResponse struct {
	Total int64                 `json:"total"`
	Data  []*CustomBankResponse `json:"data"`
}

type BankUpdateRequest struct {
	BankID        int    `json:"bank_id" validate:"omitempty"`
	BankBranch    string `json:"bank_branch" validate:"omitempty"`
	AccountNumber string `json:"account_number" validate:"omitempty"`
	AccountName   string `json:"account_name" validate:"omitempty"`
}
