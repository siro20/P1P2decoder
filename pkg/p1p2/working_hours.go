package p1p2

import (
	"fmt"
)

type WorkingHours struct {
	*GenericSensor
}

var WorkingHoursMainPump = WorkingHours{
	GenericSensor: newGenericSensor("MainPump",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHours{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHours); ok {
				return p.Pump.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursCompressorForHeating = WorkingHours{
	GenericSensor: newGenericSensor("CompressorForHeating",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHours{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHours); ok {
				return p.CompressorForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursCompressorForCooling = WorkingHours{
	GenericSensor: newGenericSensor("CompressorForCooling",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHours{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHours); ok {
				return p.CompressorForCooling.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursCompressorForDHW = WorkingHours{
	GenericSensor: newGenericSensor("CompressorForDHW",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHours{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHours); ok {
				return p.CompressorForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBoilerForHeating = WorkingHours{
	GenericSensor: newGenericSensor("BoilerForHeating",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursGas{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursGas); ok {
				return p.BoilerForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBoilerForDHW = WorkingHours{
	GenericSensor: newGenericSensor("BoilerForDHW",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursGas{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursGas); ok {
				return p.BoilerForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var CounterNumberOfBoilerStarts = WorkingHours{
	GenericSensor: newGenericSensor("BoilerStarts",
		"count",
		"",
		"Number of gas boiler starts.",
		"mdi:reload",
		int(0),
		&PacketB8RespOperatingHoursGas{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursGas); ok {
				return p.NumberOfBoilerStarts.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var CounterNumberOfCompressorStarts = WorkingHours{
	GenericSensor: newGenericSensor("CompressorStarts",
		"count",
		"",
		"Number of compressor starts.",
		"mdi:reload",
		int(0),
		&PacketB8RespOperatingHoursCompressor{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursCompressor); ok {
				return p.NumberOfCompressorStarts.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBackupHeater1ForHeating = WorkingHours{
	GenericSensor: newGenericSensor("BackupHeater1ForHeating",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursHeater{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursHeater); ok {
				return p.BackupHeater1ForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBackupHeater1ForDHW = WorkingHours{
	GenericSensor: newGenericSensor("BackupHeater1ForDHW",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursHeater{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursHeater); ok {
				return p.BackupHeater1ForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBackupHeater2ForHeating = WorkingHours{
	GenericSensor: newGenericSensor("BackupHeater2ForHeating",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursHeater{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursHeater); ok {
				return p.BackupHeater2ForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var WorkingHoursBackupHeater2ForDHW = WorkingHours{
	GenericSensor: newGenericSensor("BackupHeater2ForDHW",
		"working_hours",
		"h",
		"Timer",
		"mdi:calendar-clock",
		int(0),
		&PacketB8RespOperatingHoursHeater{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(PacketB8RespOperatingHoursHeater); ok {
				return p.BackupHeater2ForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}
