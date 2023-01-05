package p1p2

import (
	"fmt"
)

type SoftwareVersion struct {
	*GenericSensor
}

var ControlUnitSoftwareVersion = SoftwareVersion{
	GenericSensor: newGenericSensor("Control",
		"software",
		"",
		"Software version string",
		"mdi:information",
		"-",
		&Packet13Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet13Resp); ok {
				return fmt.Sprintf("ID%04X", p.ControlSoftwareVersion), nil
			}
			return "unknown", fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var HeatPumpSoftwareVersion = SoftwareVersion{
	GenericSensor: newGenericSensor("Heatpump",
		"software",
		"",
		"Software version string",
		"mdi:information",
		"-",
		&Packet13Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet13Resp); ok {
				return fmt.Sprintf("ID%04X", p.HeatPumpSoftwareVersion), nil
			}
			return "unknown", fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}
