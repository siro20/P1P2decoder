package p1p2

import (
	"testing"
)

var pkt10req = []byte{0x00, 0x00, 0x10, 0x01, 0x81, 0x01, 0x00, 0x00, 0x00, 0x00, 0x15, 0x00, 0x40, 0x00, 0x00, 0x08, 0x00, 0x00, 0x18, 0x00, 0x40, 0x30, 0x00, 0x6d}
var pkt10resp = []byte{0x40, 0x00, 0x10, 0x01, 0x80, 0xa1, 0x01, 0x30, 0x00, 0x18, 0x00, 0x15, 0x00, 0x5a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x08, 0x02, 0x8b}
var pkt11req = []byte{0x00, 0x00, 0x11, 0x0e, 0xe6, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc9}
var pkt11resp = []byte{0x40, 0x00, 0x11, 0x34, 0xca, 0x2f, 0x1a, 0x07, 0x00, 0x34, 0xa4, 0x34, 0xa4, 0x1c, 0x44, 0x0e, 0xe6, 0x07, 0x01, 0x00, 0x00, 0x00, 0x00, 0xe4}
var pkt12req = []byte{0x00, 0x00, 0x12, 0xc0, 0x02, 0x13, 0x38, 0x16, 0x02, 0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x41, 0x34, 0x03, 0xe8}
var pkt13resp = []byte{0x40, 0x00, 0x13, 0x30, 0x00, 0x03, 0x52, 0x00, 0x00, 0x15, 0x00, 0x00, 0x00, 0x75, 0xb2, 0xf6, 0xc1, 0x00, 0x00, 0x4b}
var pkt14req = []byte{0x00, 0x00, 0x14, 0x2d, 0x00, 0x12, 0x00, 0x2d, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1f}
var pkt14resp = []byte{0x40, 0x00, 0x14, 0x2d, 0x00, 0x12, 0x00, 0x2d, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23, 0x00, 0x29, 0x00, 0x63}
var pkt16resp = []byte{0x40, 0x00, 0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xde, 0xb4}

func Test_CRC(t *testing.T) {
	// Test short packages
	pkt := []byte{0x00, 0x00, 0x10}
	_, err := calcCRC(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{0x00, 0x00}
	_, err = calcCRC(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{0x00}
	_, err = calcCRC(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{}
	_, err = calcCRC(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
}

func Test_shortpacket(t *testing.T) {
	// Test short packages
	pkt := []byte{0x00, 0x00, 0x10}
	_, err := Decode(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{0x00, 0x00}
	_, err = Decode(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{0x00}
	_, err = Decode(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}
	pkt = []byte{}
	_, err = Decode(pkt)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p1 := make([]byte, len(pkt10req)-2)
	copy(p1, pkt10req)
	_, err = Decode(p1)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p2 := make([]byte, len(pkt11req)-2)
	copy(p2, pkt11req)
	_, err = Decode(p2)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p3 := make([]byte, len(pkt14req)-2)
	copy(p3, pkt14req)
	_, err = Decode(p3)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p4 := make([]byte, len(pkt10resp)-2)
	copy(p4, pkt10resp)
	_, err = Decode(p4)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p5 := make([]byte, len(pkt11resp)-2)
	copy(p5, pkt11resp)
	_, err = Decode(p5)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p6 := make([]byte, len(pkt14resp)-2)
	copy(p6, pkt14resp)
	_, err = Decode(p6)
	if err == nil {
		t.Errorf("Expected decode error")
	}

	p7 := make([]byte, len(pkt16resp)-2)
	copy(p7, pkt16resp)
	_, err = Decode(p7)
	if err == nil {
		t.Errorf("Expected decode error")
	}
}

func Test_unsupportedpacket(t *testing.T) {
	pkt10 := make([]byte, len(pkt10req))
	copy(pkt10, pkt10req)

	pkt10[2] = 0xff
	_, err := Decode(pkt10)
	if err == nil {
		t.Errorf("Expected decode error")
	}
}

func Test_decodepkt10(t *testing.T) {
	// Test good case
	p, err := Decode(pkt10req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	} else {
		dec, ok := p.(*Packet10Req)
		if !ok {
			t.Errorf("Unexpected returned packet type")
		}

		if dec.Heating != 0x01 {
			t.Errorf("Unexpected field value")
		}
		if dec.OperationMode != 0x81 {
			t.Errorf("Unexpected field value")
		}
		if dec.DHWTank != 0x01 {
			t.Errorf("Unexpected field value")
		}
		if dec.DWHTankMode != 0x40 {
			t.Errorf("Unexpected field value")
		}
		if dec.DHWTankTargetTemperature.Decode() != 48.0 {
			t.Errorf("Unexpected field value")
		}
	}

	pkt10 := make([]byte, len(pkt10req))
	copy(pkt10, pkt10req)
	// Test invalid CRC
	pkt10[len(pkt10)-1] = 0x6e
	p, err = Decode(pkt10)
	if err == nil {
		t.Errorf("Expected CRC error")
	}

	pkt10[len(pkt10)-1] = 0x6d
	pkt10[len(pkt10)-2] = 0xff
	p, err = Decode(pkt10)
	if err == nil {
		t.Errorf("Expected CRC error")
	}

	// Test missing CRC
	pkt10 = []byte{0x00, 0x00, 0x10, 0x01, 0x81, 0x01, 0x00, 0x00, 0x00, 0x00, 0x15, 0x00, 0x40, 0x00, 0x00, 0x08, 0x00, 0x00, 0x18, 0x00, 0x40, 0x30, 0x00}
	_, err = Decode(pkt10)
	if err == nil {
		t.Errorf("Expected error")
	}
}

func Test_decodeCBpkt10req(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet10Req{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt10resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}
	_, err = Decode(pkt10req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt10resp(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet10Resp{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt10req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}

	_, err = Decode(pkt10resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt11req(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet11Req{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt11resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}
	_, err = Decode(pkt11req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt11resp(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet11Resp{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt11req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}

	_, err = Decode(pkt11resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt12req(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet12Req{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt11resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}
	_, err = Decode(pkt12req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt13resp(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet13Resp{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt11resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}
	_, err = Decode(pkt13resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt14req(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet14Req{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt14resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}
	_, err = Decode(pkt14req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt14resp(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet14Resp{}, func(p interface{}) error {
		gotcha = true
		return nil
	})
	_, err := Decode(pkt14req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if gotcha {
		t.Errorf("Callback called, but shouldn't")
	}

	_, err = Decode(pkt14resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}

func Test_decodeCBpkt16resp(t *testing.T) {
	gotcha := false
	RegisterPacketCallback(&Packet16Resp{}, func(p interface{}) error {
		gotcha = true
		return nil
	})

	_, err := Decode(pkt16resp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !gotcha {
		t.Errorf("Callback not called")
	}
}
