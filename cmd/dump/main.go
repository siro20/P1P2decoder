package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose   = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	dumpFile  = kingpin.Arg("dump", "Path to regular file to read data from").Required().String()
	newPacket = false
	packet    = ""
)

func PrintFloat(pref string, t p1p2.Temperature) {
	p1p2.TemperatureRegisterCallback(t, func(v float32) {
		if newPacket {
			fmt.Printf("\nDecoding packet '%s':\n", packet)
			newPacket = false
		}
		fmt.Printf("%-16s %-22s: %f\n", pref, t.Name(), v)
	})
}

func PrintState(pref string, s p1p2.State) {
	p1p2.StateRegisterCallback(s, func(v bool) {
		if newPacket {
			fmt.Printf("\nDecoding packet '%s':\n", packet)
			newPacket = false
		}
		fmt.Printf("%-16s %-22s: %t\n", pref, s.Name(), v)
	})
}

func main() {
	fmt.Printf("Decoding dump file with the following syntax:\n")
	fmt.Printf("[timestamp :][0x] 00 [,][ ][0x] 01 [,][ ][0x] 02 ...\n")
	fmt.Printf("Examples:\n")
	fmt.Printf(" 000: 0x01, 0x02, 0x03\n")
	fmt.Printf(" 0xFF, 0xDD, 0xEE\n")
	fmt.Printf(" de ad be ef\n")

	kingpin.Parse()

	file, err := os.Open(*dumpFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	PrintFloat("Temperature", p1p2.TempLeavingWater)
	PrintFloat("Temperature", p1p2.TempExternalSensor)
	PrintFloat("Temperature", p1p2.TempActualRoom)
	PrintFloat("Temperature", p1p2.TempRefrigerant)
	PrintFloat("Temperature", p1p2.TempGasBoiler)
	PrintFloat("Temperature", p1p2.TempReturnWater)
	PrintFloat("Temperature", p1p2.TempOutside)
	PrintFloat("Temperature", p1p2.TempAdditionalZoneTarget)
	PrintFloat("Temperature", p1p2.TempMainZoneTarget)
	PrintFloat("Temperature", p1p2.TempDomesticHotWater)
	PrintFloat("Temperature", p1p2.TempDomesticHotWaterTarget)
	PrintState("State", p1p2.StateCompressor)
	PrintState("State", p1p2.StateGas)
	PrintState("State", p1p2.StatePower)
	PrintState("State", p1p2.StateQuietMode)
	PrintState("State", p1p2.StateDHWEnable)
	PrintState("State", p1p2.StateDHW)
	PrintState("State", p1p2.StateDHWBooster)
	PrintState("Valve", p1p2.ValveThreeWay)
	PrintState("Valve", p1p2.ValveAdditionalZone)
	PrintState("Valve", p1p2.ValveMainZone)
	PrintState("Valve", p1p2.ValveCooling)
	PrintState("Valve", p1p2.ValveHeating)
	PrintState("Valve", p1p2.ValveDomesticHotWater)
	PrintState("Pump", p1p2.PumpMain)
	PrintState("Pump", p1p2.PumpDHWCirculation)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var s string
		if strings.Contains(scanner.Text(), ":") {
			s = strings.Split(scanner.Text(), ":")[1]
		}

		// Remove whitespace
		s = strings.ReplaceAll(s, ", 0x", "")
		s = strings.ReplaceAll(s, " 0x", "")
		s = strings.ReplaceAll(s, " ", "")

		buf, err := hex.DecodeString(s)
		if err != nil {
			if *verbose {
				fmt.Printf("Skipping invalid line in file: %s", scanner.Text())
			}
			continue
		}

		packet = s
		newPacket = true
		_, err = p1p2.Decode(buf)
		if err != nil {
			if *verbose {
				fmt.Printf("Error decoding packet '%s': %v\n", s, err)
			}
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
