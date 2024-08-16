package facility

import (
	"tabi-booking/internal/model"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Facility{})}
}

// DB represents the client for account table
type DB struct {
	*dbcore.DB
}

func (d *DB) GetFacilityList(db *gorm.DB, ids []int64) ([]*model.Facility, error) {
	facilities := []*model.Facility{}
	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}

	lq.Filter.And("id IN (?)", ids)
	lq.Sort = []string{"class_en ASC", "name_en ASC", "class_vi ASC", "name_vi ASC"}
	if err := d.List(db, &facilities, lq, nil); err != nil {
		return nil, err
	}

	return facilities, nil
}
