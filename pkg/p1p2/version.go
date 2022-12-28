package p1p2

import (
	"fmt"
	"time"
)

type SoftwareVersion struct {
	id             SensorID
	N              string    `json:"name"`
	S              string    `json:"version"`
	Ts             time.Time `json:"last_updated"`
	decodeFunc     func(pkt interface{}) (string, error)
	changeCallback []func(Sensor, string)
	updateCallback []func(Sensor, string)
}

func (s *SoftwareVersion) Unit() string {
	return ""
}

func (s *SoftwareVersion) Type() string {
	return "software"
}

func (s *SoftwareVersion) Version() string {
	return s.S
}

func (s *SoftwareVersion) Name() string {
	return s.N
}

func (s *SoftwareVersion) Icon() string {
	return "mdi:information"
}

func (s *SoftwareVersion) Description() string {
	return ""
}

func (s *SoftwareVersion) Value() interface{} {
	return s.Version()
}

func (s *SoftwareVersion) ID() SensorID {
	return s.id
}

func (s *SoftwareVersion) RegisterUpdateCallback(f func(Sensor, interface{})) error {
	s.updateCallback = append(s.updateCallback, func(s Sensor, value string) {
		f(s, value)
	})
	return nil
}

func (s *SoftwareVersion) RegisterStateChangedCallback(f func(Sensor, interface{})) error {
	s.changeCallback = append(s.changeCallback, func(s Sensor, value string) {
		f(s, value)
	})
	return nil
}

func (s *SoftwareVersion) RegisterStateChangedWithHysteresisCallback(hysteresis float32, f func(s Sensor, value interface{})) error {
	return fmt.Errorf("Not supported")
}

func (s *SoftwareVersion) SetValue(newVersion string) {
	var oldversion string

	oldversion = s.S
	s.S = newVersion
	s.Ts = time.Now()

	for _, cb := range s.updateCallback {
		cb(s, newVersion)
	}
	if oldversion != newVersion {
		for _, cb := range s.changeCallback {
			cb(s, newVersion)
		}
	}
}

func (s *SoftwareVersion) LastUpdated() time.Time {
	return s.Ts
}

func (s *SoftwareVersion) Decode(pkt interface{}) error {
	v, err := s.decodeFunc(pkt)
	if err != nil {
		return err
	}
	s.SetValue(v)
	return nil
}

func newSoftwareVersion(name string, f func(pkt interface{}) (string, error)) SoftwareVersion {
	id := IDcnt
	IDcnt++
	return SoftwareVersion{
		N:              name,
		id:             id,
		S:              "unknown",
		decodeFunc:     f,
		changeCallback: []func(Sensor, string){},
		updateCallback: []func(Sensor, string){},
	}
}

var ControlUnitSoftwareVersion = newSoftwareVersion("Control",
	func(pkt interface{}) (string, error) {
		if p, ok := pkt.(Packet13Resp); ok {
			return fmt.Sprintf("ID%04X", p.ControlSoftwareVersion), nil
		}
		return "unknown", fmt.Errorf("Wrong message")
	})

var HeatPumpSoftwareVersion = newSoftwareVersion("Heatpump",
	func(pkt interface{}) (string, error) {
		if p, ok := pkt.(Packet13Resp); ok {
			return fmt.Sprintf("ID%04X", p.HeatPumpSoftwareVersion), nil
		}
		return "unknown", fmt.Errorf("Wrong message")
	})

func init() {
	Packet13RespRegisterCallback(func(p Packet13Resp) { _ = ControlUnitSoftwareVersion.Decode(p) })
	Packet13RespRegisterCallback(func(p Packet13Resp) { _ = HeatPumpSoftwareVersion.Decode(p) })
}
