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

func TestTemperatureRegisterChangeCallback(t *testing.T) {
	gotcha := false
	TempLeavingWater.SetValue(0.0)

	TemperatureRegisterChangeCallback(TempLeavingWater, func(newVal float32, oldVal float32) {
		if newVal == 1.0 && oldVal == 0.0 {
			gotcha = true
		}
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
