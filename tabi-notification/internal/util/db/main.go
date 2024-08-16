package dbutil

import (
	"log"
	"os"
	"time"

	dbcore "github.com/namhoai1109/tabi/core/db"

	"github.com/imdatngo/gowhere"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// New creates new database connection to the database server
func New(dbPsn string, enableLog bool) (*gorm.DB, error) {
	// Add your DB related stuffs here, such as:
	// - gorm.DefaultTableNameHandler
	// - gowhere.DefaultConfig
	gowhere.DefaultConfig.Dialect = gowhere.DialectPostgreSQL

	config := new(gorm.Config)

	namingStrategy := schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	if enableLog {
		config.Logger = newLogger
	}

	config.NamingStrategy = namingStrategy

	return dbcore.New(dbPsn, config)
}

// NewDB creates new DB instance
func NewDB(model interface{}) *dbcore.DB {
	return &dbcore.DB{Model: model}
}
