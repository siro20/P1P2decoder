package p1p2

import (
	"fmt"
	"time"
)

type SoftwareVersion struct {
	id         sensorID
	N          string    `json:"name"`
	S          string    `json:"version"`
	Ts         time.Time `json:"last_updated"`
	decodeFunc func(pkt interface{}) (string, error)
}

func (s *SoftwareVersion) Version() string {
	return s.S
}

func (s *SoftwareVersion) Name() string {
	return s.N
}

func (s *SoftwareVersion) SetValue(newVersion string) {
	s.S = newVersion
	s.Ts = time.Now()

	cbs, ok := svCB[s.id]
	if ok {
		for _, cb := range cbs {
			cb(newVersion)
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

var svCB map[sensorID][]func(p string) = map[sensorID][]func(p string){}

func SoftwareVersionRegisterCallback(t *SoftwareVersion, f func(v string)) {
	svCB[t.id] = append(svCB[t.id], f)
}

func newSoftwareVersion(name string, f func(pkt interface{}) (string, error)) SoftwareVersion {
	id := IDcnt
	IDcnt++
	return SoftwareVersion{
		N:          name,
		id:         id,
		S:          "unknown",
		decodeFunc: f,
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
