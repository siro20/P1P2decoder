package main

import (
	"fmt"

	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

func main() {

	var pkt10req = []byte{0x00, 0x00, 0x10, 0x01, 0x81, 0x01, 0x00, 0x00, 0x00, 0x00, 0x15, 0x00, 0x40, 0x00, 0x00, 0x08, 0x00, 0x00, 0x18, 0x00, 0x40, 0x30, 0x00, 0x6d}

	msg, err := p1p2.Decode(pkt10req)
	if err != nil {
		fmt.Printf("Error decoding packet: %v\n", err)
	}
	if p, ok := msg.(p1p2.Packet10Req); ok {
		fmt.Printf("Received message was Packet 10 Request: %+v\n", p)

	} else if p, ok := msg.(p1p2.Packet10Resp); ok {
		fmt.Printf("Received message was Packet 10 Response: %+v\n", p)
	}
}
