package p1p2

import (
	"fmt"
	"time"
)

type Temperature struct {
	id         sensorID
	N          string    `json:"name"`
	V          float32   `json:"value"`
	Desc       string    `json:"description"`
	Ts         time.Time `json:"last_updated"`
	T          string    `json:"type"`
	U          string    `json:"unit"`
	RangeMin   float32   `json:"range_min"`
	RangeMax   float32   `json:"range_max"`
	SetPoint   bool      `json:"setpoint"`
	decodeFunc func(pkt interface{}) (float32, error)
}

func (t *Temperature) Unit() string {
	return t.U
}

func (t *Temperature) Type() string {
	return t.T
}

func (t *Temperature) Value() interface{} {
	return t.V
}

func (t *Temperature) Name() string {
	return t.N
}

func (t *Temperature) Description() string {
	return t.Desc
}

func (t *Temperature) LastUpdated() time.Time {
	return t.Ts
}

func (t *Temperature) SetValue(newVal float32) {
	oldValue := t.V
	t.V = newVal
	t.Ts = time.Now()

	// Notify event listeners
	cbs, ok := temperatureCB[t.id]
	if ok {
		for _, cb := range cbs {
			cb(newVal)
		}
	}

	cached, ok := temperatureCache[t.id]
	if (ok && cached != newVal) || !ok {
		// Notify change event listeners
		cbs, ok := temperatureChangedCB[t.id]
		if ok {
			for _, cb := range cbs {
				cb(newVal, oldValue)
			}
		}
	}
	temperatureCache[t.id] = newVal
}

func (t *Temperature) Decode(pkt interface{}) {
	val, err := t.decodeFunc(pkt)
	if err == nil {
		t.SetValue(val)
	}
}

var temperatureCB map[sensorID][]func(p float32) = map[sensorID][]func(p float32){}
var temperatureChangedCB map[sensorID][]func(newVal float32, oldVal float32) = map[sensorID][]func(newVal float32, oldVal float32){}

var temperatureCache map[sensorID]float32 = map[sensorID]float32{}

// TemperatureRegisterCallback registers a data update callback
func TemperatureRegisterCallback(t Temperature, f func(t float32)) {
	temperatureCB[t.id] = append(temperatureCB[t.id], f)
}

// TemperatureRegisterChangeCallback registers a data change callback
func TemperatureRegisterChangeCallback(t Temperature, f func(newVal float32, oldVal float32)) {
	temperatureChangedCB[t.id] = append(temperatureChangedCB[t.id], f)
}

func newTemperature(name string,
	desc string,
	decodeFunction func(pkt interface{}) (float32, error),
	min float32,
	max float32,
	setpoint bool) Temperature {
	id := IDcnt
	IDcnt++
	return Temperature{
		id:         id,
		N:          name,
		Desc:       desc,
		V:          0.0,
		U:          "°C",
		T:          "gauge",
		Ts:         time.Unix(0, 0),
		decodeFunc: decodeFunction,
		RangeMin:   min,
		RangeMax:   max,
		SetPoint:   setpoint,
	}
}

var TempLeavingWater = newTemperature("LeavingWater", "Water temperature sent to the heat emitters.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.LWTtemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 90, false,
)

var TempDomesticHotWater = newTemperature("DomesticHotWater", "Actual domestic hot water temperature.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.DHWtemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 60, false,
)

var TempDomesticHotWaterTarget = newTemperature("DomesticHotWaterTarget", "Target temperature for domestic hot water.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.DHWTankTargetTemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 60, true,
)
var TempMainZoneTarget = newTemperature("MainZoneTarget", "Target temperature for main zone hot water.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet14Resp); ok {
			return p.MainZoneTargetTemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 90, true,
)
var TempAdditionalZoneTarget = newTemperature("AdditionalZoneTarget", "Target temperature for additional zone hot water.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet14Resp); ok {
			return p.AddZoneargetTemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 90, true,
)
var TempOutside = newTemperature("Outside", "Outside air temperature.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.Outsidetemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	-30, 50, false,
)
var TempReturnWater = newTemperature("ReturnWater", "Water temperature received back from the heat emitters.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.RWT.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 90, false,
)
var TempGasBoiler = newTemperature("GasBoiler", "Water temperature in the gas boiler.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.GasBoiler.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	10, 90, false,
)
var TempRefrigerant = newTemperature("Refrigerant", "Temperature of the refrigant.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.Refrigerant.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	-30, 50, false,
)
var TempActualRoom = newTemperature("ActualRoom", "Room temperature of the main control.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.ActualRoomtemperature.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	0, 40, false,
)
var TempExternalSensor = newTemperature("ExternalSensor", "External sensor or averaged outside temperature.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet11Resp); ok {
			return p.ExternalTemperatureSensor.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	-30, 50, false,
)
var TempDeltaT = newTemperature("DeltaT", "Delta between LWT and RWT.",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet14Req); ok {
			return p.DeltaT.Decode(), nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	-20, 20, false,
)

func init() {
	Packet14ReqRegisterCallback(func(p Packet14Req) { TempDeltaT.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempExternalSensor.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempActualRoom.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempRefrigerant.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempGasBoiler.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempReturnWater.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempOutside.Decode(p) })
	Packet14RespRegisterCallback(func(p Packet14Resp) { TempAdditionalZoneTarget.Decode(p) })
	Packet14RespRegisterCallback(func(p Packet14Resp) { TempMainZoneTarget.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { TempDomesticHotWaterTarget.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempLeavingWater.Decode(p) })
	Packet11RespRegisterCallback(func(p Packet11Resp) { TempDomesticHotWater.Decode(p) })
}
