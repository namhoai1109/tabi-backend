package migration

import (
	"fmt"
	"strings"
	"tabi-booking/config"
	"tabi-booking/internal/model"
	dbutil "tabi-booking/internal/util/db"
	"time"

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
			ID: "202311091504",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Account{},
					&model.Representative{},
					&model.Company{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"account",
					"representative",
					"company",
				)
			},
		},
		{
			ID: "202311261638",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Branch{},
					&model.Bank{},
					&model.BranchManager{},
					&model.GeneralType{},
					&model.Facility{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"branch",
					"bank",
					"branch_manager",
					"general_type",
					"facility",
				)
			},
		},
		{
			ID: "202312161126",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Branch{},
					&model.Bank{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"branch",
					"bank",
				)
			},
		},
		{
			ID: "202312262108",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.User{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"user",
				)
			},
		},
		{
			ID: "202401131609",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Bank{},
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "bank" DROP COLUMN bank_name;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("bank"); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "bank" ADD COLUMN bank_name text;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202301150902",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Room{},
					&model.RoomType{},
					&model.RoomTypeOfBranch{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"room",
					"room_type",
					"room_type_of_branch",
				)
			},
		},
		{
			ID: "202301231006",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Bank{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"bank",
				)
			},
		},
		{
			ID: "202401281529",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.ReservationReduction{},
					&model.FactureReduction{},
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "room" ADD COLUMN max_price FLOAT;`,
					`ALTER TABLE "room" DROP COLUMN price_vn;`,
					`ALTER TABLE "room" DROP COLUMN price_us;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable(
					"reservation_reduction",
					"facture_reduction",
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "room" ADD COLUMN price_vn FLOAT;`,
					`ALTER TABLE "room" ADD COLUMN price_us FLOAT;`,
					`ALTER TABLE "room" DROP COLUMN max_price;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202301261610",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "room" ADD COLUMN branch_id INTEGER;`,
					`ALTER TABLE "room" ADD CONSTRAINT fk_branch_rooms FOREIGN KEY (branch_id) REFERENCES branch(id);`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE room DROP CONSTRAINT fk_branch_rooms;`,
					`ALTER TABLE room DROP COLUMN branch_id;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "20230128121",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.SavedBranch{},
				); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"saved_branch",
				)
			},
		},
		{
			ID: "202403091601",
			Migrate: func(tx *gorm.DB) error {
				type Booking struct {
					ID            int        `json:"id" gorm:"primaryKey"`
					UserID        int        `json:"user_id"`
					RoomID        int        `json:"room_id"`
					CheckInDate   *time.Time `json:"check_in_date" gorm:"type:timestamp,default:NULL"`
					CheckOutDate  *time.Time `json:"check_out_date" gorm:"type:timestamp,default:NULL"`
					PaymentMethod string     `json:"payment_method" gorm:"type:varchar(6)"`
					TotalPrice    float64    `json:"total_price"`
					Status        string     `json:"status" gorm:"type:varchar(3)"`
					Note          string     `json:"note" gorm:"type:text"`

					model.Base
				}

				if err := tx.AutoMigrate(
					&Booking{},
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "booking" ADD CONSTRAINT fk_user_bookings FOREIGN KEY (user_id) REFERENCES "user"(id);`,
					`ALTER TABLE "booking" ADD CONSTRAINT fk_room_bookings FOREIGN KEY (room_id) REFERENCES "room"(id);`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("booking"); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "booking" DROP CONSTRAINT fk_user_bookings;`,
					`ALTER TABLE "booking" DROP CONSTRAINT fk_room_bookings;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202403111556",
			Migrate: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "room" DROP COLUMN available;`,
					`ALTER TABLE "room" DROP COLUMN check_in_date;`,
					`ALTER TABLE "room" DROP COLUMN check_out_date;`,
					`ALTER TABLE "room_type" DROP COLUMN number_rooms;`,
					`ALTER TABLE "room" ADD COLUMN quantity INT DEFAULT 1;`,
					`ALTER TABLE "branch" ADD COLUMN is_active BOOLEAN DEFAULT false;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "room" ADD COLUMN available BOOLEAN;`,
					`ALTER TABLE "room" ADD COLUMN check_in_date TIMESTAMP;`,
					`ALTER TABLE "room" ADD COLUMN check_out_date TIMESTAMP;`,
					`ALTER TABLE "room_type" ADD COLUMN number_rooms INT;`,
					`ALTER TABLE "room" DROP COLUMN quantity;`,
					`ALTER TABLE "branch" DROP COLUMN is_active;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202403221533",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(
					&model.Survey{},
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "booking" ADD COLUMN quantity INTEGER DEFAULT 1;`,
					`ALTER TABLE "room" ADD COLUMN cancellation_time_value INTEGER;`,
					`ALTER TABLE "room" ADD COLUMN cancellation_time_unit VARCHAR(5);`,
					`ALTER TABLE "room" ADD COLUMN other_policy TEXT;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("survey"); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "booking" DROP COLUMN quantity;`,
					`ALTER TABLE "room" DROP COLUMN cancellation_time_value;`,
					`ALTER TABLE "room" DROP COLUMN cancellation_time_unit;`,
					`ALTER TABLE "room" DROP COLUMN other_policy;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202404010923",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "booking" ADD COLUMN reason TEXT DEFAULT NULL;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "booking" DROP COLUMN reason;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202404021623",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "room" DROP COLUMN cancellation_time_value;`,
					`ALTER TABLE "room" DROP COLUMN cancellation_time_unit;`,
					`ALTER TABLE "room" DROP COLUMN other_policy;`,

					`ALTER TABLE "branch" ADD COLUMN cancellation_time_value INTEGER;`,
					`ALTER TABLE "branch" ADD COLUMN cancellation_time_unit VARCHAR(5);`,
					`ALTER TABLE "branch" ADD COLUMN general_policy TEXT;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "room" ADD COLUMN cancellation_time_value INTEGER;`,
					`ALTER TABLE "room" ADD COLUMN cancellation_time_unit VARCHAR(5);`,
					`ALTER TABLE "room" ADD COLUMN other_policy TEXT;`,

					`ALTER TABLE "branch" DROP COLUMN cancellation_time_value;`,
					`ALTER TABLE "branch" DROP COLUMN cancellation_time_unit;`,
					`ALTER TABLE "branch" DROP COLUMN general_policy;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202404181630",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "survey" DROP COLUMN time;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "survey" ADD COLUMN time text;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202404200930",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "branch_manager" ALTER COLUMN representative_id SET DEFAULT NULL;`,
					`ALTER TABLE "branch" ALTER COLUMN company_id SET DEFAULT NULL;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "branch_manager" ALTER COLUMN representative_id DROP DEFAULT;`,
					`ALTER TABLE "branch" ALTER COLUMN company_id DROP DEFAULT;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202405310742",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "survey" DROP COLUMN budget;`,
					`ALTER TABLE "survey" DROP COLUMN place_size;`,
					`ALTER TABLE "survey" ADD COLUMN seasons text;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "survey" ADD COLUMN budget text;`,
					`ALTER TABLE "survey" ADD COLUMN place_size text;`,
					`ALTER TABLE "survey" DROP COLUMN seasons;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202406020828",
			Migrate: func(tx *gorm.DB) error {
				type Rating struct {
					ID       int    `json:"id" gorm:"primaryKey"`
					UserID   int    `json:"user_id"`
					BranchID int    `json:"branch_id"`
					Rating   int    `json:"rating"`
					Comment  string `json:"comment"`
					model.Base
				}

				if err := tx.AutoMigrate(
					&Rating{},
				); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "rating" ADD CONSTRAINT fk_user_ratings FOREIGN KEY (user_id) REFERENCES "user"(id);`,
					`ALTER TABLE "rating" ADD CONSTRAINT fk_branch_ratings FOREIGN KEY (branch_id) REFERENCES "branch"(id);`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Migrator().DropTable("rating"); err != nil {
					return err
				}

				changes := []string{
					`ALTER TABLE "rating" DROP CONSTRAINT fk_user_ratings;`,
					`ALTER TABLE "rating" DROP CONSTRAINT fk_branch_ratings;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
		{
			ID: "202406161102",
			Migrate: func(tx *gorm.DB) error {

				changes := []string{
					`ALTER TABLE "branch" ADD COLUMN website_url text DEFAULT NULL;`,
					`ALTER TABLE "branch" ADD COLUMN tax_number text DEFAULT NULL;`,
				}

				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},

			Rollback: func(tx *gorm.DB) error {
				changes := []string{
					`ALTER TABLE "branch" DROP COLUMN website_url;`,
					`ALTER TABLE "branch" DROP COLUMN tax_number;`,
				}
				return migrationcore.ExecMultiple(tx, strings.Join(changes, " "))
			},
		},
	})

	return nil
}
