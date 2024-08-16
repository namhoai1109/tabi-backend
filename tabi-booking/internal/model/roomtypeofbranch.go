package model

// swagger:model RoomTypeOfBranch
type RoomTypeOfBranch struct {
	ID         int  `json:"id" gorm:"primaryKey"`
	RoomTypeID int  `json:"room_type_id"`
	BranchID   int  `json:"branch_id"`
	Linked     bool `json:"linked"`
	Base

	RoomType *RoomType `gorm:"foreignKey:RoomTypeID"`
	Branch   *Branch   `gorm:"foreignKey:BranchID"`
}
