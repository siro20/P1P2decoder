package p1p2

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TemperatureDBEntry struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Temperature float32    `json:"temperature"`
	Sensor      sensorID   `json:"-"`
}

type StateDBEntry struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	State     bool       `json:"state"`
	Sensor    sensorID   `json:"-"`
}

type DB struct {
	dbGorm    *gorm.DB
	lastState map[int]bool
}

func (db *DB) registerState(s State) {
	StateRegisterChangeCallback(s, func(val bool) {
		db.dbGorm.Create(&StateDBEntry{State: val, Sensor: s.id})
	})
}

func (db *DB) registerTempCallback(t Temperature) {
	TemperatureRegisterChangeCallback(t, func(val float32) {
		db.dbGorm.Create(&TemperatureDBEntry{Temperature: val, Sensor: t.id})
	})
}

func OpenDB(path string) (p1p2db *DB, err error) {
	p1p2db = &DB{}
	p1p2db.dbGorm, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("Failed to open DB: %v", err)
		return
	}

	// Migrate the schema
	p1p2db.dbGorm.AutoMigrate(&TemperatureDBEntry{})
	p1p2db.dbGorm.AutoMigrate(&StateDBEntry{})

	p1p2db.registerTempCallback(TempLeavingWater)
	p1p2db.registerTempCallback(TempExternalSensor)
	p1p2db.registerTempCallback(TempDomesticHotWater)
	p1p2db.registerTempCallback(TempDomesticHotWaterTarget)
	p1p2db.registerTempCallback(TempMainZoneTarget)
	p1p2db.registerTempCallback(TempOutside)
	p1p2db.registerTempCallback(TempReturnWater)
	p1p2db.registerTempCallback(TempGasBoiler)
	p1p2db.registerTempCallback(TempRefrigerant)
	p1p2db.registerTempCallback(TempActualRoom)
	p1p2db.registerTempCallback(TempDeltaT)

	p1p2db.registerState(ValveDomesticHotWater)
	p1p2db.registerState(ValveHeating)
	p1p2db.registerState(ValveCooling)
	p1p2db.registerState(ValveMainZone)
	p1p2db.registerState(ValveAdditionalZone)
	p1p2db.registerState(ValveThreeWay)
	p1p2db.registerState(StatePower)
	p1p2db.registerState(StateQuietMode)
	p1p2db.registerState(StateDHWBooster)
	p1p2db.registerState(StateDHWEnable)
	p1p2db.registerState(StateDHW)
	p1p2db.registerState(StateGas)
	p1p2db.registerState(StateCompressor)
	p1p2db.registerState(PumpMain)
	p1p2db.registerState(PumpDHWCirculation)

	return
}

func (db *DB) GetTemperature(t Temperature) (e TemperatureDBEntry, err error) {
	err = db.dbGorm.Where("sensor = ?", t.id).First(&e).Error
	return
}

func (db *DB) GetTemperatures(t Temperature, Since time.Time) (e []TemperatureDBEntry, err error) {
	err = db.dbGorm.Where("sensor = ? AND created_at BETWEEN ? AND ?", t.id, Since, time.Now()).Find(&e).Error
	return
}

func (db *DB) GetState(s State) (e StateDBEntry, err error) {
	err = db.dbGorm.Where("sensor = ?", s.id).First(&e).Error
	return
}

func (db *DB) GetStates(s State, Since time.Time) (e []StateDBEntry, err error) {
	err = db.dbGorm.Where("sensor = ? AND created_at BETWEEN ? AND ?", s.id, Since, time.Now()).Find(&e).Error
	return
}
