package p1p2

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_DB_GetTempeature(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "go-test")
	if err != nil {
		t.SkipNow()
	}
	defer os.Remove(file.Name())
	db, err := OpenDB(file.Name())
	if err != nil {
		t.Errorf("Opening the DB failed with: %v", err)
		return
	}
	_, err = db.GetTemperature(TempLeavingWater)
	if err == nil {
		t.Errorf("Expected to receive an error")
	}
	TempLeavingWater.SetValue(1.0)
	dbEntry, err := db.GetTemperature(TempLeavingWater)
	if err != nil {
		t.Errorf("Failed reading a DB entry: %v", err)
		return
	}
	if dbEntry.Sensor != TempLeavingWater.id {
		t.Errorf("Entry has wrong id: %v", dbEntry.Sensor)
	}
	if dbEntry.Temperature > 1.0 || dbEntry.Temperature < 1.0 {
		t.Errorf("Entry has wrong value: %v", dbEntry.Temperature)
	}
	if time.Since(dbEntry.CreatedAt).Seconds() > 60 {
		t.Errorf("Entry has wrong CreatedAt: %v", dbEntry.CreatedAt)
	}
}

func Test_DB_GetState(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "go-test")
	if err != nil {
		t.SkipNow()
	}
	defer os.Remove(file.Name())
	db, err := OpenDB(file.Name())
	if err != nil {
		t.Errorf("Opening the DB failed with: %v", err)
		return
	}
	_, err = db.GetState(ValveCooling)
	if err == nil {
		t.Errorf("Expected to receive an error")
	}

	ValveCooling.SetValue(true)
	dbEntry, err := db.GetState(ValveCooling)
	if err != nil {
		t.Errorf("Failed reading a DB entry: %v", err)
		return
	}
	if dbEntry.Sensor != ValveCooling.id {
		t.Errorf("Entry has wrong id: %v", dbEntry.Sensor)
	}
	if !dbEntry.State {
		t.Errorf("Entry has wrong value: %v", dbEntry.State)
	}
	if time.Since(dbEntry.CreatedAt).Seconds() > 60 {
		t.Errorf("Entry has wrong CreatedAt: %v", dbEntry.CreatedAt)
	}
	_, err = db.GetState(ValveHeating)
	if err == nil {
		t.Errorf("Expected to receive an error")
		return
	}
}

func Test_DB_GetTempeatures(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "go-test")
	if err != nil {
		t.SkipNow()
	}
	defer os.Remove(file.Name())
	db, err := OpenDB(file.Name())
	if err != nil {
		t.Errorf("Opening the DB failed with: %v", err)
		return
	}

	TempLeavingWater.SetValue(1.0)
	time.Sleep(time.Millisecond)
	TempLeavingWater.SetValue(2.0)
	time.Sleep(time.Millisecond)
	TempLeavingWater.SetValue(3.0)

	dbEntries, err := db.GetTemperatures(TempLeavingWater, time.Now().Add(-time.Minute))
	if err != nil {
		t.Errorf("Failed reading a DB entry: %v", err)
		return
	}
	if len(dbEntries) != 3 {
		t.Errorf("Got wrong number of entries: %d", len(dbEntries))
		return
	}
	for i := range dbEntries {
		if dbEntries[i].Temperature > float32(i+1) || dbEntries[i].Temperature < float32(i+1) {
			t.Errorf("Entry %d has wrong value: %v", i, dbEntries[i].Temperature)
		}
		if dbEntries[i].Sensor != TempLeavingWater.id {
			t.Errorf("Entry %d has wrong id: %v", i, dbEntries[i].Sensor)
		}
		if time.Since(dbEntries[i].CreatedAt).Seconds() > 60 {
			t.Errorf("Entry %d has wrong CreatedAt: %v", i, dbEntries[i].CreatedAt)
		}
	}

}

func Test_DB_GetStates(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "go-test")
	if err != nil {
		t.SkipNow()
	}
	defer os.Remove(file.Name())
	db, err := OpenDB(file.Name())
	if err != nil {
		t.Errorf("Opening the DB failed with: %v", err)
		return
	}

	ValveCooling.SetValue(true)
	time.Sleep(time.Millisecond)
	ValveCooling.SetValue(false)
	time.Sleep(time.Millisecond)
	ValveCooling.SetValue(true)

	dbEntries, err := db.GetStates(ValveCooling, time.Now().Add(-time.Minute))
	if err != nil {
		t.Errorf("Failed reading a DB entry: %v", err)
		return
	}
	if len(dbEntries) != 3 {
		t.Errorf("Got wrong number of entries: %d", len(dbEntries))
		return
	}
	for i := range dbEntries {
		if dbEntries[i].State != (i != 1) {
			t.Errorf("Entry %d has wrong value: %v", i, dbEntries[i].State)
		}
		if dbEntries[i].Sensor != ValveCooling.id {
			t.Errorf("Entry %d has wrong id: %v", i, dbEntries[i].Sensor)
		}
		if time.Since(dbEntries[i].CreatedAt).Seconds() > 60 {
			t.Errorf("Entry %d has wrong CreatedAt: %v", i, dbEntries[i].CreatedAt)
		}
	}
}
