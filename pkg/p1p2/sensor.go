package p1p2

import "time"

type sensorID int

var IDcnt sensorID

type Sensor interface {
	Unit() string
	Type() string
	Value() interface{} // Usually the direct value. Can be a sum if Divisor is > 1
	Name() string
	Description() string
	LastUpdated() time.Time
}
