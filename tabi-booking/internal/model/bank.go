package model

// swagger:model Bank
type Bank struct {
	ID int `json:"id" gorm:"primaryKey"`
	// BankName      string `json:"bank_name" gorm:"type:varchar(255)"`
	BankID        int    `json:"bank_id" gorm:"type:integer"`
	BankBranch    string `json:"bank_branch" gorm:"type:varchar(128)"`
	AccountNumber string `json:"account_number" gorm:"type:varchar(32)"`
	AccountName   string `json:"account_name" gorm:"type:varchar(32)"`
	BranchID      int    `json:"branch_id"`

	Branch *Branch `gorm:"foreignKey:BranchID"`
	Base
}
