package p1p2

import (
	"fmt"
	"reflect"
	"time"
)

type SensorID int

var IDcnt SensorID

type Sensor interface {
	ID() SensorID
	Unit() string
	// Type can be one of "temperature", "gauge", "state", "valve", "software", "pump", "working_hours", "energy", "count"
	Type() string
	Value() interface{} // Usually the direct value. Can be a sum if Divisor is > 1
	Name() string
	Icon() string
	Description() string
	LastUpdated() time.Time

	// Invoke function on Value update
	RegisterUpdateCallback(f func(s Sensor, value interface{})) error
	// Invoke function on Value change
	RegisterStateChangedCallback(f func(s Sensor, value interface{})) error
	// Invoke function on Value change with hysteresis
	// Not supported by all sensors
	RegisterStateChangedWithHysteresisCallback(hysteresis float32, f func(s Sensor, value interface{})) error
}

type GenericSensor struct {
	id              SensorID
	changeCallback  []func(Sensor, interface{})
	updateCallback  []func(Sensor, interface{})
	decodeFunction  func(pkt interface{}) (interface{}, error)
	cacheCnt        int
	hysteresisCache map[int]float32
	unsetValue      bool
	N               string      `json:"name"`
	Desc            string      `json:"description"`
	Ts              time.Time   `json:"last_updated"`
	T               string      `json:"type"`
	U               string      `json:"unit"`
	RangeMin        float32     `json:"range_min"`
	RangeMax        float32     `json:"range_max"`
	SetPoint        bool        `json:"setpoint"`
	I               string      `json:"icon"`
	V               interface{} `json:"value"`
}

func (g *GenericSensor) LastUpdated() time.Time {
	return g.Ts
}

func (g *GenericSensor) Description() string {
	return g.Desc
}

func (g *GenericSensor) Icon() string {
	return g.I
}

func (g *GenericSensor) Name() string {
	return g.N
}

func (g *GenericSensor) Type() string {
	return g.T
}

func (g *GenericSensor) Unit() string {
	return g.U
}

func (g *GenericSensor) ID() SensorID {
	return g.id
}

func (g *GenericSensor) Value() interface{} {
	if g.unsetValue {
		return nil
	}
	return g.V
}

func (g *GenericSensor) decode(pkt interface{}) error {
	newValue, err := g.decodeFunction(pkt)
	if reflect.TypeOf(newValue) != reflect.TypeOf(g.V) {
		return fmt.Errorf("Decode function for sensor %s returned incorrect type %s, expected %s",
			g.Name(), reflect.TypeOf(newValue), reflect.TypeOf(g.V))
	}
	if err != nil {
		return fmt.Errorf("Decode for sensor %s failed with %v", g.Name(), err)
	}
	return g.SetValue(newValue)
}

func (g *GenericSensor) RegisterUpdateCallback(f func(Sensor, interface{})) error {
	g.updateCallback = append(g.updateCallback, f)
	return nil
}

func (g *GenericSensor) RegisterStateChangedCallback(f func(Sensor, interface{})) error {
	g.changeCallback = append(g.changeCallback, f)
	return nil
}

func (g *GenericSensor) RegisterStateChangedWithHysteresisCallback(hysteresis float32, f func(Sensor, interface{})) error {
	switch g.V.(type) {
	case float32:
		if hysteresis == 0.0 {
			return g.RegisterStateChangedCallback(f)
		} else {
			id := g.cacheCnt
			g.cacheCnt++

			g.hysteresisCache[id] = 0.0
			g.changeCallback = append(g.changeCallback, func(_ Sensor, value interface{}) {
				newVal, ok := value.(float32)
				if !ok {
					return
				}
				oldVal := g.hysteresisCache[id]
				if oldVal+hysteresis <= newVal || oldVal-hysteresis >= newVal {
					f(g, newVal)
					g.hysteresisCache[id] = newVal
				}
			})
		}
	case float64:
		if hysteresis == 0.0 {
			return g.RegisterStateChangedCallback(f)
		} else {
			id := g.cacheCnt
			g.cacheCnt++

			g.hysteresisCache[id] = 0.0
			g.changeCallback = append(g.changeCallback, func(_ Sensor, value interface{}) {
				newVal, ok := value.(float64)
				if !ok {
					return
				}
				oldVal := g.hysteresisCache[id]
				if float64(oldVal+hysteresis) <= newVal || float64(oldVal-hysteresis) >= newVal {
					f(g, newVal)
					g.hysteresisCache[id] = float32(newVal)
				}
			})
		}
	default:
		return fmt.Errorf("Not supported")
	}

	return nil
}

func (g *GenericSensor) SetValue(newValue interface{}) error {
	oldValue := g.V
	g.Ts = time.Now()
	g.V = newValue

	for _, cb := range g.updateCallback {
		cb(g, newValue)
	}
	if newValue != oldValue || g.unsetValue {
		for _, cb := range g.changeCallback {
			cb(g, newValue)
		}
	}
	g.unsetValue = false
	return nil
}

func newGenericSensor(sensorName string,
	sensorType string,
	sensorUnit string,
	sensorDescription string,
	sensorIcon string,
	initialValue interface{},
	decodeFunc func(pkt interface{}) (interface{}, error),
	minValue float32,
	maxValue float32,
	setpoint bool) GenericSensor {

	id := IDcnt
	IDcnt++
	return GenericSensor{
		N:               sensorName,
		Desc:            sensorDescription,
		Ts:              time.Unix(0, 0),
		T:               sensorType,
		U:               sensorUnit,
		id:              id,
		I:               sensorIcon,
		RangeMin:        minValue,
		RangeMax:        maxValue,
		SetPoint:        setpoint,
		V:               initialValue,
		changeCallback:  []func(Sensor, interface{}){},
		updateCallback:  []func(Sensor, interface{}){},
		hysteresisCache: map[int]float32{},
		decodeFunction:  decodeFunc,
		unsetValue:      true,
	}
}
