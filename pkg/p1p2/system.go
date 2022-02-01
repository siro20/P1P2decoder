package p1p2

type System struct {
	Status       []*State       `json:"state"`
	Valves       []*State       `json:"valves"`
	Temperatures []*Temperature `json:"temperatures"`
	Pumps        []*State       `json:"pumps"`
}

var Sys = System{
	Temperatures: []*Temperature{
		&TempLeavingWater,
		&TempDomesticHotWater,
		&TempDomesticHotWaterTarget,
		&TempMainZoneTarget,
		&TempOutside,
		&TempReturnWater,
		&TempGasBoiler,
		&TempRefrigerant,
		&TempActualRoom,
		&TempExternalSensor,
		&TempDeltaT,
	},
	Valves: []*State{
		&ValveDomesticHotWater,
		&ValveHeating,
		&ValveCooling,
		&ValveMainZone,
		&ValveAdditionalZone,
		&ValveThreeWay,
	},
	Status: []*State{
		&StatePower,
		&StateQuietMode,
		&StateDHWBooster,
		&StateDHWEnable,
		&StateDHW,
		&StateGas,
		&StateCompressor,
	},
	Pumps: []*State{
		&PumpMain,
		&PumpDHWCirculation,
	},
}

var Sensors = []Sensor{
	&TempLeavingWater,
	&TempDomesticHotWater,
	&TempDomesticHotWaterTarget,
	&TempMainZoneTarget,
	&TempOutside,
	&TempReturnWater,
	&TempGasBoiler,
	&TempRefrigerant,
	&TempActualRoom,
	&TempExternalSensor,
	&TempDeltaT,
	&ValveDomesticHotWater,
	&ValveHeating,
	&ValveCooling,
	&ValveMainZone,
	&ValveAdditionalZone,
	&ValveThreeWay,
	&StatePower,
	&StateQuietMode,
	&StateDHWBooster,
	&StateDHWEnable,
	&StateDHW,
	&StateGas,
	&StateCompressor,
	&PumpMain,
	&PumpDHWCirculation,
}
