package p1p2

import (
	"fmt"
	"time"
)

type DateTime struct {
	id             SensorID
	T              time.Time `json:"time"`
	Ts             time.Time `json:"last_updated"`
	changeCallback []func(Sensor, time.Time)
	updateCallback []func(Sensor, time.Time)
}

func (t *DateTime) Time() string {
	return t.T.String()
}

func (t *DateTime) Description() string {
	return ""
}

func (t *DateTime) Icon() string {
	return "mdi:clock"
}

func (t *DateTime) Name() string {
	return "SystemTime"
}

func (t *DateTime) Type() string {
	return "time"
}

func (t *DateTime) Unit() string {
	return "time"
}

func (t *DateTime) Value() interface{} {
	return t.Time()
}

func (t *DateTime) ID() SensorID {
	return t.id
}

func (t *DateTime) SetValue(newTime time.Time) {
	var oldTime time.Time
	oldTime = t.T
	t.T = newTime
	t.Ts = time.Now()

	for _, cb := range t.updateCallback {
		cb(t, newTime)
	}

	if !oldTime.Equal(newTime) {
		for _, cb := range t.changeCallback {
			cb(t, newTime)
		}
	}
}

func (t *DateTime) RegisterUpdateCallback(f func(Sensor, interface{})) error {
	t.updateCallback = append(t.updateCallback, func(s Sensor, value time.Time) {
		f(s, value)
	})
	return nil
}

func (t *DateTime) RegisterStateChangedCallback(f func(Sensor, interface{})) error {
	t.changeCallback = append(t.changeCallback, func(s Sensor, value time.Time) {
		f(s, value)
	})
	return nil
}

func (t *DateTime) RegisterStateChangedWithHysteresisCallback(hysteresis float32, f func(s Sensor, value interface{})) error {
	return fmt.Errorf("Not supported")
}

func (t *DateTime) LastUpdated() time.Time {
	return t.Ts
}

func (t *DateTime) Decode(pkt interface{}) (time.Time, error) {
	if p, ok := pkt.(Packet12Req); ok {
		rfc3339 := fmt.Sprintf("20%02d-%02d-%02dT%02d:%02d:00Z00:00", p.DateYear,
			p.DateMonth, p.DateDayOfMonth, p.TimeHours, p.TimeMinutes)
		return time.Parse(time.RFC3339, rfc3339)
	}
	return time.Time{}, fmt.Errorf("Wrong message")
}

func newTime() DateTime {
	id := IDcnt
	IDcnt++
	return DateTime{
		id:             id,
		T:              time.Time{},
		changeCallback: []func(Sensor, time.Time){},
		updateCallback: []func(Sensor, time.Time){},
	}
}

// SystemTime returns the current time without seconds or timezone
var SystemTime = newTime()

func init() {
	Packet12ReqRegisterCallback(func(p Packet12Req) {
		t, err := SystemTime.Decode(p)
		if err != nil {
			SystemTime.SetValue(t)
		}
	})
}
