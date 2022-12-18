package p1p2

import (
	"fmt"
	"time"
)

type State struct {
	id         sensorID
	N          string    `json:"name"`
	V          bool      `json:"value"`
	Desc       string    `json:"description"`
	Ts         time.Time `json:"last_updated"`
	T          string    `json:"type"`
	U          string    `json:"unit"`
	decodeFunc func(pkt interface{}) (bool, error)
}

func (s *State) Unit() string {
	return s.U
}

func (s *State) Type() string {
	return s.T
}

func (s *State) Value() interface{} {
	return s.V
}

func (s *State) Name() string {
	return s.N
}

func (s *State) Description() string {
	return s.Desc
}

func (s *State) LastUpdated() time.Time {
	return s.Ts
}

func (s *State) SetValue(newValue bool) {
	oldValue := s.V
	s.V = newValue
	s.Ts = time.Now()

	// Notify event listeners
	cbs, ok := stateCB[s.id]
	if ok {
		for _, cb := range cbs {
			cb(newValue)
		}
	}

	cached, ok := stateCache[s.id]
	if (ok && (cached != newValue)) || !ok {
		// Notify change event listeners
		cbs, ok := stateChangedCB[s.id]
		if ok {
			for _, cb := range cbs {
				cb(newValue, oldValue)
			}
		}
	}
	stateCache[s.id] = newValue
}

func (s *State) Decode(pkt interface{}) {
	val, err := s.decodeFunc(pkt)
	if err == nil {
		s.SetValue(val)
	}
}

var stateCB map[sensorID][]func(s bool) = map[sensorID][]func(s bool){}
var stateChangedCB map[sensorID][]func(newVal bool, oldVal bool) = map[sensorID][]func(newVal bool, oldVal bool){}

var stateCache map[sensorID]bool = map[sensorID]bool{}

// StateRegisterCallback registers a data update callback
// Might be called for the same value again and again
func StateRegisterCallback(s State, f func(s bool)) {
	stateCB[s.id] = append(stateCB[s.id], f)
}

// StateRegisterChangeCallback registers a data change callback
func StateRegisterChangeCallback(s State, f func(newVal bool, oldVal bool)) {
	stateChangedCB[s.id] = append(stateChangedCB[s.id], f)
}

func newState(n string, d string, f func(pkt interface{}) (bool, error)) State {
	id := IDcnt
	IDcnt++
	return State{
		id:         id,
		N:          n,
		Desc:       d,
		U:          "boolean",
		T:          "gauge",
		Ts:         time.Unix(0, 0),
		decodeFunc: f,
	}
}

var ValveDomesticHotWater = newState("DomesticHotWater", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x80 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var ValveHeating = newState("Heating", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x01 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var ValveCooling = newState("Cooling", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x02 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var ValveMainZone = newState("MainZone", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x20 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var ValveAdditionalZone = newState("AdditionalZone", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x40 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var ValveThreeWay = newState("ThreeWay", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Valves&0x10 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var StatePower = newState("Power", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.Heating&0x01 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var StateQuietMode = newState("QuietMode", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.QuietMode&0x02 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var StateDHWBooster = newState("DHWBooster", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Req); ok {
			return p.DWHTankMode&0x02 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var StateDHWEnable = newState("DHWEnable", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Req); ok {
			return p.DHWTank&0x01 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})
var StateDHW = newState("DHW", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Req); ok {
			return p.DWHTankMode&0x40 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var StateGas = newState("Gas", "",
	func(pkt interface{}) (bool, error) {
		//FIXME
		if p, ok := pkt.(Packet10Req); ok {
			return p.OperationMode&0x80 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var StateCompressor = newState("Compressor", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.PumpAndCompressorStatus&0x01 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var PumpMain = newState("Main", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			return p.PumpAndCompressorStatus&0x08 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

var PumpDHWCirculation = newState("DHWCirculation", "",
	func(pkt interface{}) (bool, error) {
		if p, ok := pkt.(Packet10Resp); ok {
			// FIXME
			return p.PumpAndCompressorStatus&0x01 > 0, nil
		}
		return false, fmt.Errorf("Wrong message")
	})

func init() {
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveDomesticHotWater.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveHeating.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveCooling.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveMainZone.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveAdditionalZone.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { ValveThreeWay.Decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHWBooster.Decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHWEnable.Decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateDHW.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateQuietMode.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StatePower.Decode(p) })
	Packet10ReqRegisterCallback(func(p Packet10Req) { StateGas.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { StateCompressor.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { PumpMain.Decode(p) })
	Packet10RespRegisterCallback(func(p Packet10Resp) { PumpDHWCirculation.Decode(p) })
}
