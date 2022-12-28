package p1p2

import (
	"testing"
)

func TestSoftwareVersionRegisterCallback(t *testing.T) {
	gotcha := false
	ControlUnitSoftwareVersion.RegisterUpdateCallback(func(s Sensor, v interface{}) {
		gotcha = true
		version, ok := v.(string)
		if ok && version != "ID75b2" {
			t.Errorf("Got wrong control unit software version: %s", v)
		}
	})
	HeatPumpSoftwareVersion.SetValue("00000")
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	gotcha = false
	ControlUnitSoftwareVersion.SetValue("ID75b2")

	if !gotcha {
		t.Error("Callback not triggered, but shoudld")
	}
}

func TestSoftwareVersion(t *testing.T) {
	ControlUnitSoftwareVersion.RegisterUpdateCallback(func(s Sensor, v interface{}) {
		version, ok := v.(string)
		if ok && version != "ID75B2" {
			t.Errorf("Got wrong control unit software version: %s", v)
		}
	})
	HeatPumpSoftwareVersion.RegisterUpdateCallback(func(s Sensor, v interface{}) {
		version, ok := v.(string)
		if ok && version != "IDF6C1" {
			t.Errorf("Got wrong heat pump software version: %s", v)
		}
	})

	_, err := Decode(pkt13resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
