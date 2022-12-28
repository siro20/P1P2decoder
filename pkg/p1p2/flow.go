package p1p2

import (
	"fmt"
	"time"
)

type Flow struct {
	id              SensorID
	N               string    `json:"name"`
	V               float32   `json:"value"`
	Desc            string    `json:"description"`
	Ts              time.Time `json:"last_updated"`
	T               string    `json:"type"`
	U               string    `json:"unit"`
	I               string    `json:"icon"`
	decodeFunc      func(pkt interface{}) (float32, error)
	changeCallback  []func(Sensor, float32)
	updateCallback  []func(Sensor, float32)
	cacheCnt        int
	hysteresisCache map[int]float32
}

func (f *Flow) Unit() string {
	return f.U
}

func (f *Flow) Type() string {
	return f.T
}

func (f *Flow) Value() interface{} {
	return f.V
}

func (f *Flow) Name() string {
	return f.N
}

func (f *Flow) Icon() string {
	return f.I
}

func (f *Flow) Description() string {
	return f.Desc
}

func (f *Flow) LastUpdated() time.Time {
	return f.Ts
}

func (f *Flow) ID() SensorID {
	return f.id
}

func (f *Flow) RegisterUpdateCallback(cb func(Sensor, interface{})) error {
	f.updateCallback = append(f.updateCallback, func(s Sensor, value float32) {
		cb(f, value)
	})
	return nil
}

func (f *Flow) RegisterStateChangedCallback(cb func(Sensor, interface{})) error {
	f.changeCallback = append(f.changeCallback, func(s Sensor, value float32) {
		cb(f, value)
	})
	return nil
}

func (f *Flow) RegisterStateChangedWithHysteresisCallback(hysteresis float32, cb func(Sensor, interface{})) error {
	return fmt.Errorf("Not supported")
}

func (f *Flow) SetValue(newVal float32) {
	var oldVal float32
	oldVal = f.V
	f.V = newVal
	f.Ts = time.Now()

	// Notify event listeners
	for _, cb := range f.updateCallback {
		cb(f, newVal)
	}
	if oldVal != newVal {
		for _, cb := range f.changeCallback {
			cb(f, newVal)
		}
	}
}

func (f *Flow) Decode(pkt interface{}) {
	val, err := f.decodeFunc(pkt)
	if err == nil {
		f.SetValue(val)
	}
}

func newFlow(name string,
	desc string,
	decodeFunction func(pkt interface{}) (float32, error),
	icon string) Flow {
	id := IDcnt
	IDcnt++
	return Flow{
		id:              id,
		N:               name,
		Desc:            desc,
		V:               0.0,
		U:               "l/min",
		T:               "gauge",
		Ts:              time.Unix(0, 0),
		decodeFunc:      decodeFunction,
		I:               icon,
		changeCallback:  []func(Sensor, float32){},
		updateCallback:  []func(Sensor, float32){},
		hysteresisCache: map[int]float32{},
	}
}

var MainPumpFlow = newFlow("MainPump", "Flow in Liter / Minute",
	func(pkt interface{}) (float32, error) {
		if p, ok := pkt.(Packet13Resp); ok {
			return float32(p.FlowDeciLiterPerMin) / 10, nil
		}
		return 0.0, fmt.Errorf("Wrong message")
	},
	"mdi:gauge",
)

func init() {
	Packet13RespRegisterCallback(func(p Packet13Resp) { MainPumpFlow.Decode(p) })
}
