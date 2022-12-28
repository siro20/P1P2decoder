package p1p2

import (
	"testing"
	"time"
)

func TestTimeRegisterUpdateCallback(t *testing.T) {
	gotcha := false
	now := time.Now()

	SystemTime.RegisterUpdateCallback(func(s Sensor, value interface{}) {
		_, ok := value.(time.Time)
		if ok {
			gotcha = true
		}
	})

	gotcha = false
	SystemTime.SetValue(now)
	if !gotcha {
		t.Error("Callback not triggered")
	}
	gotcha = false
	SystemTime.SetValue(now)
	if !gotcha {
		t.Error("Callback not triggered")
	}

}

func TestTimeRegisterChangeCallback(t *testing.T) {
	now := time.Now()
	SystemTime.SetValue(now)

	gotcha := false
	SystemTime.RegisterStateChangedCallback(func(s Sensor, value interface{}) {
		_, ok := value.(time.Time)
		if ok {
			gotcha = true
		}
	})

	gotcha = false
	SystemTime.SetValue(now)
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}
	gotcha = false
	SystemTime.SetValue(now)
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}
}
