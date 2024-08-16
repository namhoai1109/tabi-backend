package migration

import (
	"fmt"
	"tabi-file/config"
	"tabi-file/internal/model"
	dbutil "tabi-file/internal/util/db"

	migrationcore "github.com/namhoai1109/tabi/core/migration"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Run executes the migration
func Run() (respErr error) {
	fmt.Println("Start migration function...")
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := dbutil.New(cfg.DbDsn, true)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				respErr = fmt.Errorf("%s", x)
			case error:
				respErr = x
			default:
				respErr = fmt.Errorf("unknown error: %+v", x)
			}
		}
	}()

	fmt.Println("db connected: " + db.Name())

	initSQL := "CREATE TABLE IF NOT EXISTS migrations (id VARCHAR(255) PRIMARY KEY)"
	if err := db.Exec(initSQL).Error; err != nil {
		return err
	}

	migrationcore.Run(db, []*gormigrate.Migration{
		{
			ID: "202311250854",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&model.File{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("file")
			},
		},
	})

	return nil
}
