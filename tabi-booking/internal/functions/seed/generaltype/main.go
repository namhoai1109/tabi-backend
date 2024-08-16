package generaltype

import (
	"context"
	"fmt"
	"tabi-booking/config"
	generaltypedb "tabi-booking/internal/db/generaltype"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"
	dbutil "tabi-booking/internal/util/db"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

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

	generalTypeDB := generaltypedb.NewDB()
	dataFromDB := []*model.GeneralType{}
	checkErr(generalTypeDB.List(db, &dataFromDB, nil, nil))
	dataFromFile, err := readAccommodationDataFromFile(ctx, cfg)

	if len(dataFromDB) == 0 {
		checkErr(err)
		checkErr(generalTypeDB.Create(db, &dataFromFile))
	} else {
		checkErr(updateInsertData(db, generalTypeDB, dataFromDB, dataFromFile))
	}

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readAccommodationDataFromFile(ctx context.Context, cfg *config.Configuration) ([]*model.GeneralType, error) {
	path := "seed/data/generaltype/"
	if cfg.Stage == "development" {
		path = "./functions/seed/data/generaltype/"
	}

	dataEN, err := util.ReadFile(path + "accommodation_type_en.json")
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to read file: %s", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to read file: %s", err))
	}
	dataVI, err := util.ReadFile(path + "accommodation_type_vi.json")
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to read file: %s", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to read file: %s", err))
	}

	resp := []*model.GeneralType{}
	for key, value := range dataEN {
		class := key
		valueEN := value.(map[string]interface{})
		valueVI := dataVI[key].(map[string]interface{})
		parent := &model.GeneralType{
			Class:   class,
			LabelEN: valueEN["label"].(string),
			LabelVI: valueVI["label"].(string),
			DescEN:  valueEN["description"].(string),
			DescVI:  valueVI["description"].(string),
			Order:   1,
		}
		resp = append(resp, parent)

		childsEN := valueEN["child"].([]interface{})
		childsVI := valueVI["child"].([]interface{})
		for iEN, childEN := range childsEN {
			for iVI, childVI := range childsVI {
				if iEN == iVI {
					childEN := childEN.(map[string]interface{})
					childVI := childVI.(map[string]interface{})
					resp = append(resp, &model.GeneralType{
						Class:   class,
						LabelEN: childEN["label"].(string),
						LabelVI: childVI["label"].(string),
						DescEN:  childEN["description"].(string),
						DescVI:  childVI["description"].(string),
						Order:   2,
					})
				}
			}
		}
	}

	return resp, nil
}

func updateInsertData(db *gorm.DB, generalTypeDB *generaltypedb.DB, dataFromDB []*model.GeneralType, dataFromFile []*model.GeneralType) error {
	if trxErr := db.Transaction(func(tx *gorm.DB) error {
		lenDataFromDB := len(dataFromDB)
		lenDataFromFile := len(dataFromFile)
		if lenDataFromDB > lenDataFromFile {
			for index, dataFile := range dataFromFile {
				if err := generalTypeDB.Update(tx, getUpdatesData(dataFile), `id = ?`, dataFromDB[index].ID); err != nil {
					return err
				}
			}
			for _, data := range dataFromDB[lenDataFromFile:] {
				if err := generalTypeDB.Delete(tx, data); err != nil {
					return err
				}
			}
		} else if lenDataFromDB < lenDataFromFile {
			for index, data := range dataFromDB {
				if err := generalTypeDB.Update(tx, getUpdatesData(dataFromFile[index]), `id = ?`, data.ID); err != nil {
					return err
				}
			}
			for _, dataFile := range dataFromFile[lenDataFromDB:] {
				if err := generalTypeDB.Create(tx, dataFile); err != nil {
					return err
				}
			}
		} else {
			for index, data := range dataFromDB {
				if err := generalTypeDB.Update(tx, getUpdatesData(dataFromFile[index]), `id = ?`, data.ID); err != nil {
					return err
				}
			}
		}

		return nil
	}); trxErr != nil {
		return trxErr
	}

	return nil
}

func getUpdatesData(data *model.GeneralType) map[string]interface{} {
	return map[string]interface{}{
		"class":    data.Class,
		"label_en": data.LabelEN,
		"label_vi": data.LabelVI,
		"desc_en":  data.DescEN,
		"desc_vi":  data.DescVI,
	}
}
