package p1p2

import "time"

type SensorID int

var IDcnt SensorID

type Sensor interface {
	ID() SensorID
	Unit() string
	// Type can be one of "temperature", "gauge", "state", "valve", "software", "pump"
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
