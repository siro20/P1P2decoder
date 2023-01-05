package p1p2

import (
	"fmt"
)

type Flow struct {
	*GenericSensor
}

var MainPumpFlow = Flow{
	GenericSensor: newGenericSensor("MainPump",
		"gauge",
		"l/min",
		"Flow in Liter / Minute",
		"mdi:gauge",
		float32(0.0),
		&Packet13Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet13Resp); ok {
				return float32(p.FlowDeciLiterPerMin) / 10, nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}
