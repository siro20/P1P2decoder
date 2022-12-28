package p1p2

import "testing"

func TestStateRegisterCallback(t *testing.T) {
	gotcha := false

	ValveDomesticHotWater.RegisterUpdateCallback(func(s Sensor, value interface{}) {
		gotcha = true
	})
	ValveHeating.SetValue(true)
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	ValveDomesticHotWater.SetValue(true)

	if !gotcha {
		t.Error("Callback not invoked")
	}

}
