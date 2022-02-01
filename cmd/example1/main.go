package main

import (
	"fmt"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

func main() {
	p1p2.TemperatureRegisterCallback(p1p2.TempLeavingWater, func(v float32) {
		fmt.Printf("Temperature %-22s: %f\n", p1p2.TempLeavingWater.Name(), v)
	})
	p1p2.TemperatureRegisterCallback(p1p2.TempExternalSensor, func(v float32) {
		fmt.Printf("Temperature %-22s: %f\n", p1p2.TempExternalSensor.Name(), v)
	})

	pkt11resp := []byte{0x40, 0x00, 0x11, 0x34, 0xca, 0x2f, 0x1a, 0x07, 0x00, 0x34, 0xa4, 0x34, 0xa4, 0x1c, 0x44, 0x0e, 0xe6, 0x07, 0x01, 0x00, 0x00, 0x00, 0x00, 0xe4}

	if _, err := p1p2.Decode(pkt11resp); err != nil {
		fmt.Printf("Error decoding packet: %v\n", err)
	}
}
