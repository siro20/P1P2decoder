package p1p2

import (
	"fmt"
)

type Temperature struct {
	*GenericSensor
}

var TempLeavingWater = Temperature{
	GenericSensor: newGenericSensor("LeavingWater",
		"temperature",
		"°C",
		"Water temperature sent to the heat emitters.",
		"mdi:thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.LWTtemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		90,
		false),
}

var TempDomesticHotWater = Temperature{
	GenericSensor: newGenericSensor("DomesticHotWater",
		"temperature",
		"°C",
		"Actual domestic hot water temperature.",
		"mdi:home-thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.DHWtemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		60,
		false),
}

var TempDomesticHotWaterTarget = Temperature{
	GenericSensor: newGenericSensor("DomesticHotWaterTarget",
		"temperature",
		"°C",
		"Target temperature for domestic hot water.",
		"mdi:thermometer-check",
		float32(0.0),
		&Packet10Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet10Resp); ok {
				return p.DHWTankTargetTemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		60,
		true),
}

var TempMainZoneTarget = Temperature{
	GenericSensor: newGenericSensor("MainZoneTarget",
		"temperature",
		"°C",
		"Target temperature for main zone hot water.",
		"mdi:thermometer-check",
		float32(0.0),
		&Packet14Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet14Resp); ok {
				return p.MainZoneTargetTemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		90,
		true),
}

var TempAdditionalZoneTarget = Temperature{
	GenericSensor: newGenericSensor("AdditionalZoneTarget",
		"temperature",
		"°C",
		"Target temperature for additional zone hot water.",
		"mdi:thermometer-check",
		float32(0.0),
		&Packet14Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet14Resp); ok {
				return p.AddZoneargetTemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		90,
		true),
}

var TempOutside = Temperature{
	GenericSensor: newGenericSensor("Outside",
		"temperature",
		"°C",
		"Outside air temperature.",
		"mdi:sun-thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.Outsidetemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		-30,
		50,
		false),
}

var TempReturnWater = Temperature{
	GenericSensor: newGenericSensor("ReturnWater",
		"temperature",
		"°C",
		"Water temperature received back from the heat emitters.",
		"mdi:thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.RWT.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		90,
		false),
}

var TempGasBoiler = Temperature{
	GenericSensor: newGenericSensor("GasBoiler",
		"temperature",
		"°C",
		"Water temperature in the gas boiler.",
		"mdi:thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.GasBoiler.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		10,
		90,
		false),
}

var TempRefrigerant = Temperature{
	GenericSensor: newGenericSensor("Refrigerant",
		"temperature",
		"°C",
		"Temperature of the refrigant.",
		"mdi:snowflake-thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.Refrigerant.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		-30,
		60,
		false),
}

var TempActualRoom = Temperature{
	GenericSensor: newGenericSensor("ActualRoom",
		"temperature",
		"°C",
		"Room temperature of the main control.",
		"mdi:thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.ActualRoomtemperature.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		0,
		40,
		false),
}

var TempExternalSensor = Temperature{
	GenericSensor: newGenericSensor("ExternalSensor",
		"temperature",
		"°C",
		"External sensor or averaged outside temperature.",
		"mdi:sun-thermometer",
		float32(0.0),
		&Packet11Resp{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet11Resp); ok {
				return p.ExternalTemperatureSensor.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		-30,
		50,
		false),
}

var TempDeltaT = Temperature{
	GenericSensor: newGenericSensor("DeltaT",
		"temperature",
		"°C",
		"Delta between LWT and RWT.",
		"mdi:thermometer",
		float32(0.0),
		&Packet14Req{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*Packet14Req); ok {
				return p.DeltaT.Decode(), nil
			}
			return float32(0.0), fmt.Errorf("Wrong message")
		},
		-20,
		20,
		false),
}
