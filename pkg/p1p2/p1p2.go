package p1p2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
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
	ExternalController0 = 0xF0
	ExternalController1 = 0xF1
)

type decoder struct {
	RequestResponse           uint8
	Address                   uint8
	Type                      uint8
	Index                     uint8
	FirstByteOfPayloadIsIndex bool
	Callback                  []func(interface{}) error
	Value                     interface{}
	Order                     binary.ByteOrder
}

var decoders []decoder = []decoder{
	{
		Request,
		Heatpump,
		0x10,
		0,
		false,
		[]func(interface{}) error{},
		&Packet10Req{},
		binary.BigEndian,
	},
	{
		Request,
		Heatpump,
		0x11,
		0,
		false,
		[]func(interface{}) error{},
		&Packet11Req{},
		binary.BigEndian,
	},
	{
		Request,
		Heatpump,
		0x12,
		0,
		false,
		[]func(interface{}) error{},
		&Packet12Req{},
		binary.BigEndian,
	},
	{
		Request,
		Heatpump,
		0x14,
		0,
		false,
		[]func(interface{}) error{},
		&Packet14Req{},
		binary.BigEndian,
	},
	{
		Request,
		ExternalController0,
		0x31,
		0,
		false,
		[]func(interface{}) error{},
		&PacketF031Req{},
		binary.BigEndian,
	},
	{
		Request,
		ExternalController0,
		0x35,
		0,
		false,
		[]func(interface{}) error{},
		&PacketF035Req{},
		binary.LittleEndian,
	},
	{
		Request,
		ExternalController1,
		0x31,
		0,
		false,
		[]func(interface{}) error{},
		&PacketF031Req{},
		binary.BigEndian,
	},
	{
		Request,
		ExternalController1,
		0x35,
		0,
		false,
		[]func(interface{}) error{},
		&PacketF035Req{},
		binary.LittleEndian,
	},
	{
		Response,
		Heatpump,
		0x10,
		0,
		false,
		[]func(interface{}) error{},
		&Packet10Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0x11,
		0,
		false,
		[]func(interface{}) error{},
		&Packet11Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0x12,
		0,
		false,
		[]func(interface{}) error{},
		&Packet12Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0x13,
		0,
		false,
		[]func(interface{}) error{},
		&Packet13Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0x14,
		0,
		false,
		[]func(interface{}) error{},
		&Packet14Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0x16,
		0,
		false,
		[]func(interface{}) error{},
		&Packet16Resp{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		0,
		true,
		[]func(interface{}) error{},
		&PacketB8RespEnergyConsumed{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		1,
		true,
		[]func(interface{}) error{},
		&PacketB8RespEnergyProduced{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		2,
		true,
		[]func(interface{}) error{},
		&PacketB8RespOperatingHours{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		3,
		true,
		[]func(interface{}) error{},
		&PacketB8RespOperatingHoursHeater{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		4,
		true,
		[]func(interface{}) error{},
		&PacketB8RespOperatingHoursCompressor{},
		binary.BigEndian,
	},
	{
		Response,
		Heatpump,
		0xb8,
		5,
		true,
		[]func(interface{}) error{},
		&PacketB8RespOperatingHoursGas{},
		binary.BigEndian,
	},
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

type Packet11Req struct {
	ActualRoomtemperature f8p8
	Reserved              [6]uint8
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

type Packet12Resp struct {
	Reserved  [12]uint8
	State     uint8
	Reserved1 [7]uint8
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

type Packet14Req struct {
	Reserved  [8]uint8
	DeltaT    sabs4
	EcoMode   uint8
	Reserved2 [5]uint8
}

type Packet14Resp struct {
	Reserved                  [15]uint8
	MainZoneTargetTemperature f8p8
	AddZoneargetTemperature   f8p8
}

type Packet16Resp struct {
	UptimeInMinutes uint16 // Counter increasing every minute
	Reserved        [7]uint8
}

// auxiliary controller ID, date, time
type PacketF031Req struct {
	Reserved       [3]uint8
	Status         [3]uint8
	DateYear       uint8
	DateMonth      uint8
	DateDayOfMonth uint8
	TimeHours      uint8
	TimeMinutes    uint8
	TimeSeconds    uint8
}

type Parameter8B struct {
	Offset uint16
	Value  uint8
}

// 8bit parameter exchange
type PacketF035Req struct {
	Parameters [6]Parameter8B
}

type PacketB8RespEnergyConsumed struct {
	BackUpHeaterForHeating uint24
	BackUpHeaterForDHW     uint24
	CompressorForHeating   uint24
	CompressorForCooling   uint24
	CompressorForDHW       uint24
	Total                  uint24
}

type PacketB8RespEnergyProduced struct {
	ForHeating uint24
	ForCooling uint24
	ForDHW     uint24
	Total      uint24
}

type PacketB8RespOperatingHours struct {
	Pump                 uint24
	CompressorForHeating uint24
	CompressorForCooling uint24
	CompressorForDHW     uint24
}

type PacketB8RespOperatingHoursHeater struct {
	BackupHeater1ForHeating uint24
	BackupHeater1ForDHW     uint24
	BackupHeater2ForHeating uint24
	BackupHeater2ForDHW     uint24
	Reserved                [6]uint8
}

type PacketB8RespOperatingHoursCompressor struct {
	Reserved                 [9]uint8
	NumberOfCompressorStarts uint24
}

type PacketB8RespOperatingHoursGas struct {
	BoilerForHeating     uint24
	BoilerForDHW         uint24
	GasUsageForHeating   uint24
	GasUsageForDHW       uint24
	NumberOfBoilerStarts uint24
	GasUsageTotal        uint24
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

func VerifyCRC(b []byte) (err error) {
	var crc byte
	if crc, err = calcCRC(b); err != nil {
		return
	}
	if crc != b[len(b)-1] {
		err = fmt.Errorf("CRC doesn't match")
	}
	return
}

func Decode(b []byte) (pkt interface{}, err error) {
	var hdr Header
	var crc byte
	var crcByte byte

	r := bytes.NewReader(b)
	if err = binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		err = fmt.Errorf("Error reading header: %v\n", err.Error())
		return
	}

	for _, d := range decoders {
		if d.Address != hdr.SlaveAddress {
			continue
		}
		if d.RequestResponse != hdr.RequestResponse {
			continue
		}
		if d.Type != hdr.Type {
			continue
		}
		if d.FirstByteOfPayloadIsIndex {
			var Index uint8
			if err = binary.Read(r, binary.LittleEndian, &Index); err != nil {
				err = fmt.Errorf("Error reading index: %v\n", err.Error())
				continue
			}
			if Index != d.Index {
				// reset reader
				r.UnreadByte()
				continue
			}
		}

		// Read data
		if err = binary.Read(r, d.Order, d.Value); err != nil {
			continue
		}
		pkt = d.Value

		// Calculate CRC
		crc, err = calcCRC(b)
		if err != nil {
			return
		}

		// Read last byte of packet
		if err = binary.Read(r, binary.BigEndian, &crcByte); err != nil {
			err = fmt.Errorf("Error reading crc: %v\n", err.Error())
			return
		}
		if crcByte != crc {
			err = fmt.Errorf("CRC of packet doesn't match: %02x != %02x", crcByte, crc)
			return
		}

		for _, cb := range d.Callback {
			err = cb(d.Value)
			if err != nil {
				fmt.Printf("Callback for packet %s failed with: %v\n", reflect.TypeOf(d.Value), err)
				err = nil
			}
		}
		return pkt, nil
	}
	if err == nil {
		err = fmt.Errorf("Paket not supported")
	}
	return
}

func RegisterPacketCallback(value interface{}, cb func(p interface{}) error) error {
	found := false
	for i := range decoders {
		if reflect.TypeOf(decoders[i].Value) == reflect.TypeOf(value) {
			decoders[i].Callback = append(decoders[i].Callback, cb)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("Unsupported packet type: %s", reflect.TypeOf(value))
	}
	return nil
}
