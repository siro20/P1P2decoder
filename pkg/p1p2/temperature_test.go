package p1p2

import "testing"

func TestTemperatureRegisterCallback(t *testing.T) {
	gotcha := false
	err := TempLeavingWater.RegisterStateChangedCallback(func(s Sensor, value interface{}) {
		gotcha = true
	})
	if err != nil {
		t.Errorf("RegisterStateChangedCallback failed with %v", err)
	}
	TempDomesticHotWater.SetValue(float32(1.0))
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	TempLeavingWater.SetValue(float32(1.0))

	if !gotcha {
		t.Error("Callback not invoked")
	}
}

func TestTemperatureRegisterChangeCallback(t *testing.T) {
	gotcha := false
	TempLeavingWater.SetValue(float32(0.0))

	err := TempLeavingWater.RegisterStateChangedWithHysteresisCallback(1.0, func(s Sensor, value interface{}) {
		newVal, ok := value.(float32)
		if !ok {
			return
		}
		if newVal == 1.0 {
			gotcha = true
		}
	})
	if err != nil {
		t.Errorf("RegisterStateChangedWithHysteresisCallback failed with %v", err)
	}
	TempDomesticHotWater.SetValue(float32(1.0))
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}

	TempLeavingWater.SetValue(float32(1.0))

	if !gotcha {
		t.Error("Callback not invoked")
	}
}

func TestTemperatureRegisterChangeCallbackWithHysteresis(t *testing.T) {
	gotcha := false
	gotcha2 := false
	TempLeavingWater.SetValue(float32(0.0))

	err := TempLeavingWater.RegisterStateChangedWithHysteresisCallback(1.0, func(s Sensor, value interface{}) {
		gotcha = true
	})
	if err != nil {
		t.Errorf("RegisterStateChangedWithHysteresisCallback failed with %v", err)
	}
	err = TempLeavingWater.RegisterStateChangedWithHysteresisCallback(0.5, func(s Sensor, value interface{}) {
		gotcha2 = true
	})
	if err != nil {
		t.Errorf("RegisterStateChangedWithHysteresisCallback failed with %v", err)
	}
	gotcha = false
	gotcha2 = false
	TempLeavingWater.SetValue(float32(0.5))
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}
	if !gotcha2 {
		t.Error("Callback not invoked")
	}
	gotcha = false
	gotcha2 = false
	TempLeavingWater.SetValue(float32(1.0))
	if !gotcha {
		t.Error("Callback not invoked")
	}
	if !gotcha2 {
		t.Error("Callback not invoked")
	}
	gotcha = false
	gotcha2 = false
	TempLeavingWater.SetValue(float32(1.5))
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}
	if !gotcha2 {
		t.Error("Callback not invoked")
	}

	gotcha = false
	gotcha2 = false
	TempLeavingWater.SetValue(float32(2.0))
	if !gotcha {
		t.Error("Callback not invoked")
	}
	if !gotcha2 {
		t.Error("Callback not invoked")
	}
}
