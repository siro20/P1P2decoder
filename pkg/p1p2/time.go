package p1p2

import (
	"fmt"
	"time"
)

type DateTime struct {
	*GenericSensor
}

var SystemTime = DateTime{
	GenericSensor: newGenericSensor("LogTime",
		"time",
		"",
		"System time without daylight saving time",
		"mdi:clock",
		time.Unix(0, 0),
		&PacketF031Req{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketF031Req); ok {
				rfc3339 := fmt.Sprintf("20%02d-%02d-%02dT%02d:%02d:%02dZ", p.DateYear,
					p.DateMonth, p.DateDayOfMonth, p.TimeHours, p.TimeMinutes, p.TimeSeconds)
				return time.Parse(time.RFC3339, rfc3339)
			}

			return time.Unix(0, 0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var UITime = DateTime{
	GenericSensor: newGenericSensor("UITime",
		"time",
		"",
		"Time shown in UI",
		"mdi:clock",
		time.Unix(0, 0),
		&Packet12Req{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet12Req); ok {
				rfc3339 := fmt.Sprintf("20%02d-%02d-%02dT%02d:%02d:00Z", p.DateYear,
					p.DateMonth, p.DateDayOfMonth, p.TimeHours, p.TimeMinutes)
				return time.Parse(time.RFC3339, rfc3339)
			}

			return time.Unix(0, 0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}
