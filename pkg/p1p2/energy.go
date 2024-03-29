package p1p2

import (
	"fmt"
)

type Energy struct {
	*GenericSensor
}

var EnergyConsumedBackUpHeaterForHeating = Energy{
	GenericSensor: newGenericSensor("ConsumedBackUpHeaterForHeating",
		"energy",
		"kWh",
		"Energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.BackUpHeaterForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyConsumedBackUpHeaterForDHW = Energy{
	GenericSensor: newGenericSensor("ConsumedBackUpHeaterForDHW",
		"energy",
		"kWh",
		"Energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.BackUpHeaterForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyConsumedCompressorForHeating = Energy{
	GenericSensor: newGenericSensor("ConsumedCompressorForHeating",
		"energy",
		"kWh",
		"Compressor energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.CompressorForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyConsumedCompressorForCooling = Energy{
	GenericSensor: newGenericSensor("ConsumedCompressorForCooling",
		"energy",
		"kWh",
		"Compressor energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.CompressorForCooling.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyConsumedCompressorForDHW = Energy{
	GenericSensor: newGenericSensor("ConsumedCompressorForDHW",
		"energy",
		"kWh",
		"Compressor energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.CompressorForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyConsumedTotal = Energy{
	GenericSensor: newGenericSensor("ConsumedTotal",
		"energy",
		"kWh",
		"Total energy meter",
		"mdi:transmission-tower-export",
		int(0),
		&PacketB8RespEnergyConsumed{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyConsumed); ok {
				return p.Total.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyProducedForHeating = Energy{
	GenericSensor: newGenericSensor("ProducedForHeating",
		"energy",
		"kWh",
		"Produced heat for heating energy meter",
		"mdi:meter-electric",
		int(0),
		&PacketB8RespEnergyProduced{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyProduced); ok {
				return p.ForHeating.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyProducedForCooling = Energy{
	GenericSensor: newGenericSensor("ProducedForCooling",
		"energy",
		"kWh",
		"Produced cooling energy meter",
		"mdi:meter-electric",
		int(0),
		&PacketB8RespEnergyProduced{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyProduced); ok {
				return p.ForCooling.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyProducedForDHW = Energy{
	GenericSensor: newGenericSensor("ProducedForDHW",
		"energy",
		"kWh",
		"Produced DHW heat energy meter",
		"mdi:meter-electric",
		int(0),
		&PacketB8RespEnergyProduced{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyProduced); ok {
				return p.ForDHW.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}

var EnergyProducedTotal = Energy{
	GenericSensor: newGenericSensor("ProducedTotal",
		"energy",
		"kWh",
		"Produced total heat energy meter",
		"mdi:meter-electric",
		int(0),
		&PacketB8RespEnergyProduced{},
		func(pkt interface{}) (interface{}, error) {
			if p, ok := pkt.(*PacketB8RespEnergyProduced); ok {
				return p.Total.Decode(), nil
			}
			return int(0), fmt.Errorf("Wrong message")
		},
		0,
		0,
		false),
}
