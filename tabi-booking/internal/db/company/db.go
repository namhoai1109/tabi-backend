package company

import (
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Company{})}
}

type DB struct {
	*dbcore.DB
}

func (d *DB) GetCompanyByRepID(db *gorm.DB, rpID int) (*model.Company, error) {
	company := &model.Company{}
	if err := d.View(db, &company, `representative_id = ?`, rpID); err != nil {
		return nil, err
	}

	return company, nil
}

func (d *DB) CheckExistedBranchOfRepresentative(db *gorm.DB, rpID, branchID int) (bool, error) {
	company := &model.Company{}
	if err := d.View(db.Preload("Branches"), &company, `representative_id = ?`, rpID); err != nil {
		return false, err
	}

	for _, branch := range company.Branches {
		if branch.ID == branchID {
			return true, nil
		}
	}

	return false, nil
}

func (d *DB) CheckBankAccountOwnership(db *gorm.DB, rpID, bankID int) (bool, error) {
	company := &model.Company{}
	if err := d.View(db.Preload("Branches.Banks"), &company, `representative_id = ?`, rpID); err != nil {
		return false, err
	}

	bankIDs := []int{}

	// Might affect performance, refactor later
	for _, branch := range company.Branches {
		for _, bank := range branch.Banks {
			bankIDs = append(bankIDs, bank.ID)
		}
	}

	for _, id := range bankIDs {
		if id == bankID {
			return true, nil
		}
	}

	return false, nil
}

func (d *DB) GetBranchIDsOfRP(db *gorm.DB, rpID int) ([]int, error) {
	company := &model.Company{}
	if err := d.View(db.Preload("Branches"), &company, `representative_id = ?`, rpID); err != nil {
		return nil, err
	}

	branchIDs := []int{}
	for _, branch := range company.Branches {
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs, nil
}

func (d *DB) GetRoomTypeIDsOfRP(db *gorm.DB, rpID int) ([]int, error) {
	company := &model.Company{}
	if err := d.View(db.Preload("Branches.RoomTypes"), &company, `representative_id = ?`, rpID); err != nil {
		return nil, err
	}

	roomTypeIDs := []int{}
	for _, branch := range company.Branches {
		for _, roomType := range branch.RoomTypes {
			roomTypeIDs = append(roomTypeIDs, roomType.ID)
		}
	}

	return roomTypeIDs, nil
}
