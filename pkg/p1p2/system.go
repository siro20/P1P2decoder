package p1p2

type System struct {
	Status       []*State           `json:"state"`
	Valves       []*State           `json:"valves"`
	Temperatures []*Temperature     `json:"temperatures"`
	Pumps        []*State           `json:"pumps"`
	Software     []*SoftwareVersion `json:"software_version"`
	Flow         []*Flow            `json:"flow"`
	Energy       []*Energy          `json:"energy"`
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
	Software: []*SoftwareVersion{
		&ControlUnitSoftwareVersion,
		&HeatPumpSoftwareVersion,
	},
	Flow: []*Flow{
		&MainPumpFlow,
	},
	Energy: []*Energy{
		&EnergyConsumedBackUpHeaterForHeating,
		&EnergyConsumedBackUpHeaterForDHW,
		&EnergyConsumedCompressorForHeating,
		&EnergyConsumedCompressorForCooling,
		&EnergyConsumedCompressorForDHW,
		&EnergyConsumedTotal,
		&EnergyProducedForHeating,
		&EnergyProducedForCooling,
		&EnergyProducedForDHW,
		&EnergyProducedTotal,
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
	&ControlUnitSoftwareVersion,
	&HeatPumpSoftwareVersion,
	&MainPumpFlow,
	&EnergyConsumedBackUpHeaterForHeating,
	&EnergyConsumedBackUpHeaterForDHW,
	&EnergyConsumedCompressorForHeating,
	&EnergyConsumedCompressorForCooling,
	&EnergyConsumedCompressorForDHW,
	&EnergyConsumedTotal,
	&EnergyProducedForHeating,
	&EnergyProducedForCooling,
	&EnergyProducedForDHW,
	&EnergyProducedTotal,
}
