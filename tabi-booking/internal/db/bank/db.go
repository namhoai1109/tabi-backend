package bank

import (
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Bank{})}
}

type DB struct {
	*dbcore.DB
}
