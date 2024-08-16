package schedule

import (
	"tabi-notification/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
)

// NewDB returns a new user database instance
func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Schedule{})}
}

// DB represents the client for user table
type DB struct {
	*dbcore.DB
}
