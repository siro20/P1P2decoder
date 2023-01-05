package p1p2

// Parameter support. Only possible using an "external controller".
// An external controller answers the packet "00f030...""
// and thus requests more communication.

var Parameter35 map[int]uint8 = map[int]uint8{}

func init() {
	RegisterPacketCallback(&PacketF035Req{}, func(value interface{}) error {
		p, ok := value.(*PacketF035Req)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			Parameter35[int(p.Parameters[i].Offset)] = p.Parameters[i].Value
		}

		return nil
	})
}
