package facility

import (
	"context"
	"tabi-booking/config"
	facilitydb "tabi-booking/internal/db/facility"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"
	dbutil "tabi-booking/internal/util/db"
	"tabi-booking/internal/util/maps"

	"gorm.io/gorm"
)

// Run executes the seed
func Run(ctx context.Context) (respErr error) {
	cfg, err := config.Load()
	checkErr(err)

	db, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	checkErr(err)
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	sqlDB, err := db.DB()
	checkErr(err)
	defer sqlDB.Close()

	facilityDB := facilitydb.NewDB()

	path := "seed/data/facility/"
	if cfg.Stage == "development" {
		path = "./functions/seed/data/facility/"
	}

	trxErr := db.Transaction(func(tx *gorm.DB) error {
		dataFromDB := []*model.Facility{}
		checkErr(facilityDB.List(tx, &dataFromDB, nil, nil))
		lenDataFromDB := len(dataFromDB)

		if lenDataFromDB > 0 {
			for _, data := range dataFromDB {
				checkErr(facilityDB.Delete(tx, `id = ?`, data.ID))
			}
		}

		checkErr(insertDataFromFile(tx, facilityDB, path, "main/main_facility_en.json", "main/main_facility_vi.json", "main/map_main.json", model.FacilityTypeMain))
		checkErr(insertDataFromFile(tx, facilityDB, path, "room/room_facility_en.json", "room/room_facility_vi.json", "room/map_room.json", model.FacilityTypeRoom))
		return nil
	})

	return trxErr
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertDataFromFile(db *gorm.DB, facilityDB *facilitydb.DB, path, filenameEN, filenameVI, mapFile, facilityType string) error {

	facilityEN, err := util.ReadFile(path + filenameEN)
	checkErr(err)
	keyEN := maps.Keys(facilityEN)

	facilityVI, err := util.ReadFile(path + filenameVI)
	checkErr(err)

	mapping, err := util.ReadFile(path + mapFile)
	checkErr(err)

	for i := 0; i < len(keyEN); i++ {
		kEN := keyEN[i]
		kVI := mapping[kEN].(string)

		valuesEN := util.InterfaceToArrayString(facilityEN[kEN].([]interface{}))
		valuesVI := util.InterfaceToArrayString(facilityVI[kVI].([]interface{}))

		facilities := []*model.Facility{}
		for j := 0; j < len(valuesEN); j++ {
			facilities = append(facilities, &model.Facility{
				Type:    facilityType,
				ClassEN: kEN,
				ClassVI: kVI,
				NameEN:  valuesEN[j],
				NameVI:  valuesVI[j],
			})
		}

		checkErr(facilityDB.Create(db, &facilities))
	}

	return nil
}
