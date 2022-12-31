package p1p2

import (
	"fmt"
)

type State struct {
	GenericSensor
}

var ValveDomesticHotWater = State{
	GenericSensor: newGenericSensor("DomesticHotWater",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x80 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var ValveHeating = State{
	GenericSensor: newGenericSensor("Heating",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x01 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var ValveCooling = State{
	GenericSensor: newGenericSensor("Cooling",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x02 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var ValveMainZone = State{
	GenericSensor: newGenericSensor("MainZone",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x20 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var ValveAdditionalZone = State{
	GenericSensor: newGenericSensor("AdditionalZone",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x40 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var ValveThreeWay = State{
	GenericSensor: newGenericSensor("ThreeWay",
		"valve",
		"",
		"",
		"mdi:pipe-valve",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Valves&0x10 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateHeatingEnabled = State{
	GenericSensor: newGenericSensor("HeatingEnabled",
		"state",
		"",
		"Heating for Main Zone/Additional Zone is enabled, but not necessarily running.",
		"mdi:power",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.Heating&0x01 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateQuietMode = State{
	GenericSensor: newGenericSensor("QuietMode",
		"state",
		"",
		"",
		"mdi:volume-off",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.QuietMode&0x02 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateDHWBooster = State{
	GenericSensor: newGenericSensor("DHWBooster",
		"state",
		"",
		"",
		"mdi:power",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Req); ok {
				return p.DWHTankMode&0x02 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateDHWEnable = State{
	GenericSensor: newGenericSensor("DHWEnable",
		"state",
		"",
		"Heating for DHW is enabled, but not necessarily running.",
		"mdi:power",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Req); ok {
				return p.DHWTank&0x01 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateDHW = State{
	GenericSensor: newGenericSensor("DHW",
		"state",
		"",
		"",
		"mdi:power",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Req); ok {
				return p.DWHTankMode&0x40 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateGasEnabled = State{
	GenericSensor: newGenericSensor("GasEnabled",
		"state",
		"",
		"The gas boiler is enabled for heating, but not necessarily running.",
		"mdi:fire",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Req); ok {
				return p.OperationMode&0x80 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateBoilerRunning = State{
	GenericSensor: newGenericSensor("BoilerRunning",
		"state",
		"",
		"The gas boiler is running.",
		"mdi:fire",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.DHWActive&2 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var StateCompressor = State{
	GenericSensor: newGenericSensor("Compressor",
		"state",
		"",
		"",
		"mdi:heat-pump",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.PumpAndCompressorStatus&0x01 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var PumpMain = State{
	GenericSensor: newGenericSensor("Main",
		"pump",
		"",
		"",
		"mdi:water-pump",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				return p.PumpAndCompressorStatus&0x08 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var PumpDHWCirculation = State{
	GenericSensor: newGenericSensor("DHWCirculation",
		"pump",
		"",
		"",
		"mdi:water-pump",
		false,
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(Packet10Resp); ok {
				// FIXME
				return p.PumpAndCompressorStatus&0x01 > 0, nil
			}
			return false, fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

func init() {
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveDomesticHotWater.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveHeating.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveCooling.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveMainZone.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveAdditionalZone.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveThreeWay.decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHWBooster.decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHWEnable.decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHW.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateQuietMode.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateHeatingEnabled.decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateGasEnabled.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateCompressor.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { PumpMain.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { PumpDHWCirculation.decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateBoilerRunning.decode(p) })
}
