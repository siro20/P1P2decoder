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
	Sensor      SensorID   `json:"-"`
}

type StateDBEntry struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	State     bool       `json:"state"`
	Sensor    SensorID   `json:"-"`
}

type DB struct {
	dbGorm            *gorm.DB
	lastState         map[int]bool
	samplingDurations []time.Duration
	deletePolicy      time.Duration
}

func (db *DB) cleanTemperatures(t Temperature, Since time.Time) error {
	return db.dbGorm.Order("created_at").Where("sensor = ? AND created_at < ?", t.id, Since).Delete(&TemperatureDBEntry{}).Error
}

func (db *DB) cleanStates(s State, Since time.Time) error {
	return db.dbGorm.Order("created_at").Where("sensor = ? AND created_at < ?", s.id, Since).Delete(&StateDBEntry{}).Error
}

func (db *DB) registerState(s State) {
	s.RegisterStateChangedCallback(func(sensor Sensor, value interface{}) {
		newVal, ok := value.(bool)
		if !ok {
			return
		}
		oldVal := !newVal
		db.dbGorm.Create(&StateDBEntry{State: oldVal, Sensor: sensor.ID()})
		time.Sleep(time.Millisecond)
		db.dbGorm.Create(&StateDBEntry{State: newVal, Sensor: sensor.ID()})
		if db.deletePolicy != 0 {
			db.cleanStates(s, time.Now().Add(-db.deletePolicy))
		}
	})
}

func (db *DB) registerTempCallback(t Temperature) {
	t.RegisterStateChangedCallback(func(s Sensor, value interface{}) {
		newVal, ok := value.(float32)
		if !ok {
			return
		}
		db.dbGorm.Create(&TemperatureDBEntry{Temperature: newVal, Sensor: s.ID()})
		if db.deletePolicy != 0 {
			db.cleanTemperatures(t, time.Now().Add(-db.deletePolicy))
		}
	})
}

// SetDeletionPolicy sets the duration after samples should be deleted from database
// Default: never
func (db *DB) SetDeletionPolicy(before time.Duration) {
	db.deletePolicy = before
}

// SetAverageSamplingIntervals sets the intervals to average sampling values on
// Default: Minute, Hour, Day
// Currently unused
func (db *DB) SetAverageSamplingIntervals(dur []time.Duration) {
	db.samplingDurations = dur
}

func OpenDB(path string) (p1p2db *DB, err error) {
	p1p2db = &DB{
		samplingDurations: []time.Duration{time.Minute, time.Hour, time.Hour * 24},
		deletePolicy:      0,
	}
	p1p2db.dbGorm, err = gorm.Open(sqlite.Open(path), &gorm.Config{
		PrepareStmt: true,
	})
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
	err = db.dbGorm.Limit(1).Order("created_at desc").First(&e, "sensor = ?", t.id).Error
	return
}

func (db *DB) GetTemperatures(t Temperature, Since time.Time) (e []TemperatureDBEntry, err error) {
	err = db.dbGorm.Order("created_at").Where("sensor = ? AND created_at BETWEEN ? AND ?", t.id, Since, time.Now()).Find(&e).Error
	return
}

func (db *DB) GetState(s State) (e StateDBEntry, err error) {
	err = db.dbGorm.Limit(1).Order("created_at desc").First(&e, "sensor = ?", s.id).Error
	return
}

func (db *DB) GetStates(s State, Since time.Time) (e []StateDBEntry, err error) {
	err = db.dbGorm.Order("created_at").Where("sensor = ? AND created_at BETWEEN ? AND ?", s.id, Since, time.Now()).Find(&e).Error
	return
}
