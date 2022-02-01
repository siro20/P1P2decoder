# P1P2 decoder - Example1

The library let's you register callbacks for every temperature sensor:

-TempLeavingWater
-TempDomesticHotWater
-TempDomesticHotWaterTarget
-TempMainZoneTarget
-TempOutside
-TempReturnWater
-TempGasBoiler
-TempRefrigerant
-TempActualRoom
-TempExternalSensor
-TempDeltaT

Use the sensor to register the callback. It will be called whenever the
`p1p2.Decode()` functions successfully decodes a packet containing the sensor
data.

```
p1p2.TemperatureRegisterCallback(p1p2.TempLeavingWater, func(v float32) {
	fmt.Printf("Temperature %-22s: %f\n", p1p2.TempLeavingWater.Name(), v)
})
p1p2.TemperatureRegisterCallback(p1p2.TempExternalSensor, func(v float32) {
	fmt.Printf("Temperature %-22s: %f\n", p1p2.TempExternalSensor.Name(), v)
})

...

if _, err := p1p2.Decode(rawpacket); err != nil {
	fmt.Printf("Error decoding packet: %v\n", err)
}
```

## Output:

```
Temperature ExternalSensor        : 7.003906
Temperature LeavingWater          : 52.789062
```