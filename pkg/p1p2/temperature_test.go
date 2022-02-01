package p1p2

import "testing"

func TestTemperatureRegisterCallback(t *testing.T) {
	gotcha := false
	TemperatureRegisterCallback(TempLeavingWater, func(v float32) {
		gotcha = true
	})
	TempDomesticHotWater.SetValue(1.0)
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	TempLeavingWater.SetValue(1.0)

	if !gotcha {
		t.Error("Callback not invoked")
	}
}
