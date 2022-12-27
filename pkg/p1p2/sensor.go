package p1p2

import "time"

type sensorID int

var IDcnt sensorID

type Sensor interface {
	Unit() string
	// Type can be one of "gauge", "state", "valve", "software"
	Type() string
	Value() interface{} // Usually the direct value. Can be a sum if Divisor is > 1
	Name() string
	Icon() string
	Description() string
	LastUpdated() time.Time
}
