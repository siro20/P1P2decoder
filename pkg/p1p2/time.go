package p1p2

import (
	"fmt"
	"time"
)

type DateTime struct {
	id sensorID
	T  time.Time `json:"time"`
	Ts time.Time `json:"last_updated"`
}

func (t *DateTime) Time() string {
	return t.T.String()
}

func (t *DateTime) SetValue(newTime time.Time) {
	t.T = newTime
	t.Ts = time.Now()

	cbs, ok := timeCB[t.id]
	if ok {
		for _, cb := range cbs {
			cb(newTime)
		}
	}
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

var timeCB map[sensorID][]func(p time.Time) = map[sensorID][]func(p time.Time){}

func TimeRegisterCallback(t *DateTime, f func(t time.Time)) {
	timeCB[t.id] = append(timeCB[t.id], f)
}

func newTime() DateTime {
	id := IDcnt
	IDcnt++
	return DateTime{
		id: id,
		T:  time.Time{},
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
