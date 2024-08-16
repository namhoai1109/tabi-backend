package migration

import (
	"fmt"
	"strings"
	"tabi-notification/config"
	"tabi-notification/internal/model"
	dbutil "tabi-notification/internal/util/db"

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
			ID: "202404302056",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Device{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("device")
			},
		},
		{
			ID: "202405072035",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Schedule{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("schedule")
			},
		},
		{
			ID: "202405090938",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "device" DROP COLUMN device_uuid;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "device" DROP COLUMN device_type;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202405212120",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "schedule" DROP COLUMN time;
						ALTER TABLE "schedule"
							ADD COLUMN start_time timestamp,
							ADD COLUMN end_time timestamp,
							ADD COLUMN booking_id int`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "schedule" DROP COLUMN start_time;
						ALTER TABLE "schedule" DROP COLUMN end_time;
						ALTER TABLE "schedule" DROP COLUMN booking_id;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202407020807",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "schedule" ADD COLUMN destination_longitude TEXT;`,
					`ALTER TABLE "schedule" ADD COLUMN destination_latitude TEXT;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "schedule" DROP COLUMN destination_longitude;`,
					`ALTER TABLE "schedule" DROP COLUMN destination_latitude;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
	})

	return nil
}
