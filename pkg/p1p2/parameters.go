package p1p2

import "fmt"

// Parameter support. Only possible using an "external controller".
// An external controller answers the packet "00f030...""
// and thus requests more communication.

var Parameter35 map[int]uint8 = map[int]uint8{}

func init() {
	PacketF035ReqRegisterCallback(func(p PacketF035Req) {
		for i := range p.Parameters {
			Parameter35[int(p.Parameters[i].Offset)] = p.Parameters[i].Value
		}
		for i := 0; i < 0xffff; i++ {
			val, ok := Parameter35[i]
			if ok {
				fmt.Printf("%x = %x\n", i, val)
			}
		}
		for i := 0x162; i < 0x16e; i += 1 {
			val, ok := Parameter35[i]
			if ok {
				fmt.Printf("%c", val)
			}
		}
		fmt.Printf("\n")

	})
}
