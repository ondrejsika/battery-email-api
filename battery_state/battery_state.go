package battery_state

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BatteryState struct {
	gorm.Model
	Device string
	Level  string
}

type BatteryStateDB struct {
	DB *gorm.DB
}

func Init(driver string, connection string) *BatteryStateDB {
	var bs BatteryStateDB
	if driver == "none" {
		return nil
	}
	if driver == "sqlite" {
		db, err := gorm.Open(sqlite.Open(connection), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
		}
		bs.DB = db
	}
	if driver == "postgres" {
		db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
		}
		bs.DB = db
	}
	bs.DB.AutoMigrate(&BatteryState{})
	return &bs
}

func (bs *BatteryStateDB) Add(device string, level string) {
	bs.DB.Create(&BatteryState{
		Device: device,
		Level:  level,
	})
}

func (bs *BatteryStateDB) Get(device string) BatteryState {
	var b BatteryState
	bs.DB.Order("created_at DESC").First(&b, "device = ?", device)
	return b
}
