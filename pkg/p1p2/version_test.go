package p1p2

import (
	"testing"
)

func TestSoftwareVersionRegisterCallback(t *testing.T) {
	gotcha := false
	SoftwareVersionRegisterCallback(&ControlUnitSoftwareVersion, func(v string) {
		gotcha = true
		if v != "ID75b2" {
			t.Errorf("Got wrong control unit software version: %s", v)
		}
	})
	HeatPumpSoftwareVersion.SetValue("00000")
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	ControlUnitSoftwareVersion.SetValue("ID75b2")

	if !gotcha {
		t.Error("Callback not triggered, but shoudld")
	}
}

func TestSoftwareVersion(t *testing.T) {
	SoftwareVersionRegisterCallback(&ControlUnitSoftwareVersion, func(v string) {
		if v != "ID75B2" {
			t.Errorf("Got wrong control unit software version: %s", v)
		}
	})
	SoftwareVersionRegisterCallback(&HeatPumpSoftwareVersion, func(v string) {
		if v != "IDF6C1" {
			t.Errorf("Got wrong heat pump software version: %s", v)
		}
	})

	_, err := Decode(pkt13resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
