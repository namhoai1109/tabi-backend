package roomtypeofbranch

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.RoomTypeOfBranch{})}
}

type DB struct {
	*dbcore.DB
}

func (d *DB) GetRoomTypes(db *gorm.DB, ctx context.Context, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomType, error) {
	roomTypes := []*model.RoomType{}
	joinedDB := db.Table("room_type_of_branch").
		Select("DISTINCT ON (room_type_of_branch.room_type_id) room_type_of_branch.*, room_type.*").
		Joins("INNER JOIN room_type ON room_type.id = room_type_of_branch.room_type_id").
		Where("room_type.deleted_at IS NULL")
	if err := d.List(joinedDB, &roomTypes, lq, nil); err != nil {
		return nil, err
	}

	fullList := []*model.RoomType{}
	lq.PerPage = 0
	lq.Page = 0
	countedDB := db.Table("room_type_of_branch").
		Select("DISTINCT ON (room_type_of_branch.room_type_id) room_type_of_branch.*, room_type.*").
		Joins("INNER JOIN room_type ON room_type.id = room_type_of_branch.room_type_id").
		Where("room_type.deleted_at IS NULL")
	if err := d.List(countedDB, &fullList, lq, nil); err != nil {
		return nil, err
	}
	fmt.Println("fullList", len(fullList))
	if count != nil {
		*count = int64(len(fullList))
	}

	return roomTypes, nil
}

func (d *DB) GetRoomTypeIDsOfBranch(db *gorm.DB, ctx context.Context, lq *dbcore.ListQueryCondition, branchID int) ([]int, error) {
	roomTypeIDs := []int{}
	joinedDB := db.Table("room_type_of_branch").
		Select("room_type_id").
		Where("branch_id = ? AND linked = true", branchID)
	if err := d.List(joinedDB, &roomTypeIDs, lq, nil); err != nil {
		return nil, err
	}

	return roomTypeIDs, nil
}
