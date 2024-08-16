package survey

import (
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Survey{})}
}

type DB struct {
	*dbcore.DB
}
