package user

import (
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.User{})}
}

// DB represents the client for account table
type DB struct {
	*dbcore.DB
}
