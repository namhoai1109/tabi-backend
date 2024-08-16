package model

//swagger:model SavedBranch
type SavedBranch struct {
	ID       int `json:"id" gorm:"primaryKey"`
	UserID   int `json:"user_id"`
	BranchID int `json:"branch_id"`

	User   *User   `gorm:"foreignKey:UserID"`
	Branch *Branch `gorm:"foreignKey:BranchID"`
	Base
}
