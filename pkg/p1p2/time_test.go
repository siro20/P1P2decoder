package p1p2

import (
	"testing"
	"time"
)

func TestTimeRegisterCallback(t *testing.T) {
	gotcha := false
	TimeRegisterCallback(&SystemTime, func(v time.Time) {
		gotcha = true
	})
	SystemTime.SetValue(time.Now())
	if gotcha {
		t.Error("Callback triggered, but shouldn't")
	}
}
