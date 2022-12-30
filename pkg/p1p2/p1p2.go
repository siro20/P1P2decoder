package p1p2

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type f8p8 int16

func (f f8p8) Decode() float32 {
	return float32(f) / 256.0
}

func (f f8p8) Encode(n float32) {
	f = f8p8(int16(n * 256))
}

type sabs4 uint8

func (s sabs4) Decode() float32 {
	return float32(s&0xF) - float32(s&0x10)
}

type uint24 [3]uint8

func (u uint24) Decode() int {
	return int(u[0])<<16 | int(u[1])<<8 | int(u[2])
}

type Header struct {
	RequestResponse uint8
	SlaveAddress    uint8
	Type            uint8
}

const (
	Request  = 0
	Response = 0x40
	Other    = 0x80
)

const (
	Heatpump            = 0
	ExternalController0 = 0xF0 + iota
	ExternalController1 = 0xF0 + iota
)

type Packet interface {
	Crc() bool
}

// Packet10Req doesn't change on state
type Packet10Req struct {
	Heating                  uint8
	OperationMode            uint8
	DHWTank                  uint8
	Reserved                 [4]uint8
	TargetRoomTemperature    uint8
	Reserved1                uint8
	Flags                    uint8
	QuietMode                uint8
	Reserved2                [6]uint8
	DWHTankMode              uint8
	DHWTankTargetTemperature f8p8
}

func (p Packet10Req) Crc() bool {
	return true
}

var packet10ReqCB []func(p Packet10Req) = []func(p Packet10Req){}

func Packet10ReqRegisterCallback(f func(p Packet10Req)) {
	packet10ReqCB = append(packet10ReqCB, f)
}

type Packet10Resp struct {
	Heating                  uint8
	Reserved                 uint8
	Valves                   uint8
	ThreeWayValve            uint8
	DHWTankTargetTemperature f8p8
	Reserved1                [5]uint8
	QuietMode                uint8
	Reserved2                [6]uint8
	PumpAndCompressorStatus  uint8
	DHWActive                uint8
}

func (p Packet10Resp) Crc() bool {
	return true
}

var packet10RespCB []func(p Packet10Resp) = []func(p Packet10Resp){}

func Packet10RespRegisterCallback(f func(p Packet10Resp)) {
	packet10RespCB = append(packet10RespCB, f)
}

type Packet11Req struct {
	ActualRoomtemperature f8p8
	Reserved              [6]uint8
}

func (p Packet11Req) Crc() bool {
	return true
}

var packet11ReqCB []func(p Packet11Req) = []func(p Packet11Req){}

func Packet11ReqRegisterCallback(f func(p Packet11Req)) {
	packet11ReqCB = append(packet11ReqCB, f)
}

type Packet11Resp struct {
	LWTtemperature            f8p8
	DHWtemperature            f8p8
	Outsidetemperature        f8p8
	RWT                       f8p8
	GasBoiler                 f8p8
	Refrigerant               f8p8
	ActualRoomtemperature     f8p8
	ExternalTemperatureSensor f8p8
	Reserved                  [4]uint8
}

func (p Packet11Resp) Crc() bool {
	return true
}

var packet11RespCB []func(p Packet11Resp) = []func(p Packet11Resp){}

func Packet11RespRegisterCallback(f func(p Packet11Resp)) {
	packet11RespCB = append(packet11RespCB, f)
}

type Packet12Req struct {
	NewHourIndicator uint8
	DayOfWeek        uint8
	TimeHours        uint8
	TimeMinutes      uint8
	DateYear         uint8
	DateMonth        uint8
	DateDayOfMonth   uint8
	Reserved         [5]uint8
	Flags            uint8
	Flags2           uint8
	Reserved2        uint8
}

func (p Packet12Req) Crc() bool {
	return true
}

var packet12ReqCB []func(p Packet12Req) = []func(p Packet12Req){}

func Packet12ReqRegisterCallback(f func(p Packet12Req)) {
	packet12ReqCB = append(packet12ReqCB, f)
}

type Packet13Resp struct {
	DHWTankTargetTemperature f8p8     // From Boiler?
	Flags                    [2]uint8 // Unknown
	Reserved                 [2]uint8
	Flags2                   uint8 // Lower bits change per second. Flow in L/min?
	Reserved2                uint8
	FlowDeciLiterPerMin      uint16
	ControlSoftwareVersion   uint16
	HeatPumpSoftwareVersion  uint16
	Reserved3                [2]uint8
}

func (p Packet13Resp) Crc() bool {
	return true
}

var packet13RespCB []func(p Packet13Resp) = []func(p Packet13Resp){}

func Packet13RespRegisterCallback(f func(p Packet13Resp)) {
	packet13RespCB = append(packet13RespCB, f)
}

type Packet14Req struct {
	Reserved  [8]uint8
	DeltaT    sabs4
	EcoMode   uint8
	Reserved2 [5]uint8
}

func (p Packet14Req) Crc() bool {
	return true
}

var packet14ReqCB []func(p Packet14Req) = []func(p Packet14Req){}

func Packet14ReqRegisterCallback(f func(p Packet14Req)) {
	packet14ReqCB = append(packet14ReqCB, f)
}

type Packet14Resp struct {
	Reserved                  [15]uint8
	MainZoneTargetTemperature f8p8
	AddZoneargetTemperature   f8p8
}

func (p Packet14Resp) Crc() bool {
	return true
}

var packet14RespCB []func(p Packet14Resp) = []func(p Packet14Resp){}

func Packet14RespRegisterCallback(f func(p Packet14Resp)) {
	packet14RespCB = append(packet14RespCB, f)
}

type Packet16Resp struct {
	UptimeInMinutes uint16 // Counter increasing every minute
	Reserved        [7]uint8
}

func (p Packet16Resp) Crc() bool {
	return true
}

var packet16RespCB []func(p Packet16Resp) = []func(p Packet16Resp){}

func Packet16RespRegisterCallback(f func(p Packet16Resp)) {
	packet16RespCB = append(packet16RespCB, f)
}

type PacketB8RespEnergyConsumed struct {
	BackUpHeaterForHeating uint24
	BackUpHeaterForDHW     uint24
	CompressorForHeating   uint24
	CompressorForCooling   uint24
	CompressorForDHW       uint24
	Total                  uint24
}

func (p PacketB8RespEnergyConsumed) Crc() bool {
	return false
}

var packetB8RespEnergyConsumed []func(p PacketB8RespEnergyConsumed) = []func(p PacketB8RespEnergyConsumed){}

func PacketB8RespEnergyConsumedRegisterCallback(f func(p PacketB8RespEnergyConsumed)) {
	packetB8RespEnergyConsumed = append(packetB8RespEnergyConsumed, f)
}

type PacketB8RespEnergyProduced struct {
	ForHeating uint24
	ForCooling uint24
	ForDHW     uint24
	Total      uint24
}

func (p PacketB8RespEnergyProduced) Crc() bool {
	return false
}

var packetB8RespEnergyProduced []func(p PacketB8RespEnergyProduced) = []func(p PacketB8RespEnergyProduced){}

func PacketB8RespEnergyProducedRegisterCallback(f func(p PacketB8RespEnergyProduced)) {
	packetB8RespEnergyProduced = append(packetB8RespEnergyProduced, f)
}

type PacketB8RespOperatingHours struct {
	Pump                 uint24
	CompressorForHeating uint24
	CompressorForCooling uint24
	CompressorForDHW     uint24
}

func (p PacketB8RespOperatingHours) Crc() bool {
	return false
}

var packetB8RespOperatingHours []func(p PacketB8RespOperatingHours) = []func(p PacketB8RespOperatingHours){}

func PacketB8RespOperatingHoursRegisterCallback(f func(p PacketB8RespOperatingHours)) {
	packetB8RespOperatingHours = append(packetB8RespOperatingHours, f)
}

type PacketB8RespOperatingHoursHeater struct {
	BackupHeater1ForHeating uint24
	BackupHeater1ForDHW     uint24
	BackupHeater2ForHeating uint24
	BackupHeater2ForDHW     uint24
}

func (p PacketB8RespOperatingHoursHeater) Crc() bool {
	return false
}

var packetB8RespOperatingHoursHeater []func(p PacketB8RespOperatingHoursHeater) = []func(p PacketB8RespOperatingHoursHeater){}

func PacketB8RespOperatingHoursHeaterRegisterCallback(f func(p PacketB8RespOperatingHoursHeater)) {
	packetB8RespOperatingHoursHeater = append(packetB8RespOperatingHoursHeater, f)
}

type PacketB8RespOperatingHoursCompressor struct {
	Reserved                 [9]uint8
	NumberOfCompressorStarts uint24
}

func (p PacketB8RespOperatingHoursCompressor) Crc() bool {
	return false
}

var packetB8RespOperatingHoursCompressor []func(p PacketB8RespOperatingHoursCompressor) = []func(p PacketB8RespOperatingHoursCompressor){}

func PacketB8RespOperatingHoursCompressorRegisterCallback(f func(p PacketB8RespOperatingHoursCompressor)) {
	packetB8RespOperatingHoursCompressor = append(packetB8RespOperatingHoursCompressor, f)
}

type PacketB8RespOperatingHoursGas struct {
	BoilerForHeating     uint24
	BoilerForDHW         uint24
	GasUsageForHeating   uint24
	GasUsageForDHW       uint24
	NumberOfBoilerStarts uint24
	GasUsageTotal        uint24
}

func (p PacketB8RespOperatingHoursGas) Crc() bool {
	return false
}

var packetB8RespOperatingHoursGas []func(p PacketB8RespOperatingHoursGas) = []func(p PacketB8RespOperatingHoursGas){}

func PacketB8RespOperatingHoursGasRegisterCallback(f func(p PacketB8RespOperatingHoursGas)) {
	packetB8RespOperatingHoursGas = append(packetB8RespOperatingHoursGas, f)
}

func calcCRC(b []byte) (crc byte, err error) {
	crc = 0

	if len(b) <= 3 {
		err = fmt.Errorf("Packet too small")
		return
	}

	for i := 0; i < len(b)-1; i++ {
		c := b[i]
		for j := 0; j < 8; j++ {
			if ((crc ^ c) & 0x01) > 0 {
				crc = (crc >> 1) ^ 0xd9
			} else {
				crc = crc >> 1
			}
			c >>= 1
		}
	}
	return
}

func decode(b []byte) (pkt interface{}, err error) {
	var hdr Header
	var crc byte

	r := bytes.NewReader(b)
	if err = binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		err = fmt.Errorf("Error reading header: %v\n", err.Error())
		return
	}

	// Decode package
	if hdr.RequestResponse == Response {
		if hdr.SlaveAddress == Heatpump {
			if hdr.Type == 0x10 {
				var p10r Packet10Resp
				if err = binary.Read(r, binary.BigEndian, &p10r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p10r
			} else if hdr.Type == 0x11 {
				var p11r Packet11Resp
				if err = binary.Read(r, binary.BigEndian, &p11r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p11r
			} else if hdr.Type == 0x13 {
				var p13r Packet13Resp
				if err = binary.Read(r, binary.BigEndian, &p13r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p13r
			} else if hdr.Type == 0x14 {
				var p14r Packet14Resp
				if err = binary.Read(r, binary.BigEndian, &p14r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p14r
			} else if hdr.Type == 0x16 {
				var p16r Packet16Resp
				if err = binary.Read(r, binary.BigEndian, &p16r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p16r
			} else if hdr.Type == 0xb8 {
				var DataType uint8
				var Consumed PacketB8RespEnergyConsumed
				var Produced PacketB8RespEnergyProduced
				var OperatingHours PacketB8RespOperatingHours
				var OperatingHoursHeater PacketB8RespOperatingHoursHeater
				var OperatingHoursCompressor PacketB8RespOperatingHoursCompressor
				var OperatingHoursGas PacketB8RespOperatingHoursGas

				if err = binary.Read(r, binary.BigEndian, &DataType); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				switch DataType {
				case 0:
					if err = binary.Read(r, binary.BigEndian, &Consumed); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = Consumed
				case 1:
					if err = binary.Read(r, binary.BigEndian, &Produced); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = Produced
				case 2:
					if err = binary.Read(r, binary.BigEndian, &OperatingHours); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = OperatingHours
				case 3:
					if err = binary.Read(r, binary.BigEndian, &OperatingHoursHeater); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = OperatingHoursHeater
				case 4:
					if err = binary.Read(r, binary.BigEndian, &OperatingHoursCompressor); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = OperatingHoursCompressor
				case 5:
					if err = binary.Read(r, binary.BigEndian, &OperatingHoursGas); err != nil {
						err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
						return
					}
					pkt = OperatingHoursGas
				}
			}
		}
	} else if hdr.RequestResponse == Request {
		if hdr.SlaveAddress == Heatpump {
			if hdr.Type == 0x10 {
				var p10r Packet10Req
				if err = binary.Read(r, binary.BigEndian, &p10r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p10r
			} else if hdr.Type == 0x11 {
				var p11r Packet11Req
				if err = binary.Read(r, binary.BigEndian, &p11r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p11r
			} else if hdr.Type == 0x12 {
				var p12r Packet12Req
				if err = binary.Read(r, binary.BigEndian, &p12r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p12r
			} else if hdr.Type == 0x14 {
				var p14r Packet14Req
				if err = binary.Read(r, binary.BigEndian, &p14r); err != nil {
					err = fmt.Errorf("Error reading packet payload: %v\n", err.Error())
					return
				}
				pkt = p14r
			}
		}
	}

	if pkt == nil {
		return nil, fmt.Errorf("unsupported message")
	}

	// Check CRC last
	if iface, ok := pkt.(Packet); ok && iface.Crc() {
		// Calculate CRC
		crc, err = calcCRC(b)
		if err != nil {
			return
		}

		// Read last byte of packet
		var b byte
		if err = binary.Read(r, binary.BigEndian, &b); err != nil {
			err = fmt.Errorf("Error reading crc: %v\n", err.Error())
			return
		}
		if b != crc {
			err = fmt.Errorf("CRC of packet doesn't match: %02x != %02x", b, crc)
			return
		}
	}
	return pkt, nil
}

func callbacks(pkt interface{}) {
	if p, ok := pkt.(Packet10Req); ok {
		for _, cb := range packet10ReqCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet10Resp); ok {
		for _, cb := range packet10RespCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet11Req); ok {
		for _, cb := range packet11ReqCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet11Resp); ok {
		for _, cb := range packet11RespCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet12Req); ok {
		for _, cb := range packet12ReqCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet13Resp); ok {
		for _, cb := range packet13RespCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet14Resp); ok {
		for _, cb := range packet14RespCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet14Req); ok {
		for _, cb := range packet14ReqCB {
			cb(p)
		}
	} else if p, ok := pkt.(Packet16Resp); ok {
		for _, cb := range packet16RespCB {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespEnergyConsumed); ok {
		for _, cb := range packetB8RespEnergyConsumed {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespEnergyProduced); ok {
		for _, cb := range packetB8RespEnergyProduced {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespOperatingHours); ok {
		for _, cb := range packetB8RespOperatingHours {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespOperatingHoursHeater); ok {
		for _, cb := range packetB8RespOperatingHoursHeater {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespOperatingHoursCompressor); ok {
		for _, cb := range packetB8RespOperatingHoursCompressor {
			cb(p)
		}
	} else if p, ok := pkt.(PacketB8RespOperatingHoursGas); ok {
		for _, cb := range packetB8RespOperatingHoursGas {
			cb(p)
		}
	}
}

func Decode(b []byte) (pkt interface{}, err error) {
	pkt, err = decode(b)
	if err != nil {
		return
	}
	callbacks(pkt)
	return
}
