package file

import (
	"tabi-file/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
)

// NewDB returns a new file database instance
func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.File{})}
}

// DB represents the client for file table
type DB struct {
	*dbcore.DB
}
