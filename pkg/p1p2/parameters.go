package p1p2

// Parameter support. Only possible using an "external controller".
// An external controller answers the packet "00f030...""
// and thus requests more communication.

var Parameter35 map[int]uint8 = map[int]uint8{}
var parameter35CB map[int][]func(int, uint8) = map[int][]func(int, uint8){}

var Parameter36 map[int]uint16 = map[int]uint16{}
var parameter36CB map[int][]func(int, uint16) = map[int][]func(int, uint16){}

var Parameter38 map[int]uint32 = map[int]uint32{}
var parameter38CB map[int][]func(int, uint32) = map[int][]func(int, uint32){}

var Parameter39 map[int]uint32 = map[int]uint32{}
var parameter39CB map[int][]func(int, uint32) = map[int][]func(int, uint32){}

var Parameter3b map[int]uint16 = map[int]uint16{}
var parameter3bCB map[int][]func(int, uint16) = map[int][]func(int, uint16){}

func RegisterParameter35ChangeCallback(offset int, cb func(offset int, value uint8)) {
	if s, ok := parameter35CB[offset]; ok {
		s = append(s, cb)
	} else {
		parameter35CB[offset] = []func(int, uint8){cb}
	}
}

func RegisterParameter36ChangeCallback(offset int, cb func(offset int, value uint16)) {
	if s, ok := parameter36CB[offset]; ok {
		s = append(s, cb)
	} else {
		parameter36CB[offset] = []func(int, uint16){cb}
	}
}

func RegisterParameter38ChangeCallback(offset int, cb func(offset int, value uint32)) {
	if s, ok := parameter38CB[offset]; ok {
		s = append(s, cb)
	} else {
		parameter38CB[offset] = []func(int, uint32){cb}
	}
}

func RegisterParameter39ChangeCallback(offset int, cb func(offset int, value uint32)) {
	if s, ok := parameter39CB[offset]; ok {
		s = append(s, cb)
	} else {
		parameter39CB[offset] = []func(int, uint32){cb}
	}
}

func RegisterParameter3bChangeCallback(offset int, cb func(offset int, value uint16)) {
	if s, ok := parameter3bCB[offset]; ok {
		s = append(s, cb)
	} else {
		parameter3bCB[offset] = []func(int, uint16){cb}
	}
}

func init() {
	RegisterPacketCallback(&PacketF035Req{}, func(value interface{}) error {
		p, ok := value.(*PacketF035Req)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			off := int(p.Parameters[i].Offset)
			oldVal, ok := Parameter35[off]
			if ok {
				if oldVal != p.Parameters[i].Value {
					for _, cb := range parameter35CB[off] {
						cb(off, p.Parameters[i].Value)
					}
				}
			} else {
				for _, cb := range parameter35CB[off] {
					cb(off, p.Parameters[i].Value)
				}
			}
			Parameter35[off] = p.Parameters[i].Value
		}

		return nil
	})

	RegisterPacketCallback(&PacketF038Req{}, func(value interface{}) error {
		p, ok := value.(*PacketF038Req)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			off := int(p.Parameters[i].Offset)

			oldVal, ok := Parameter38[off]
			if ok {
				if oldVal != p.Parameters[i].Value {
					for _, cb := range parameter38CB[off] {

						cb(off, p.Parameters[i].Value)
					}
				}
			} else {
				for _, cb := range parameter38CB[off] {
					cb(off, p.Parameters[i].Value)
				}
			}
			Parameter38[off] = p.Parameters[i].Value
		}

		return nil
	})

	RegisterPacketCallback(&PacketF039Req{}, func(value interface{}) error {
		p, ok := value.(*PacketF039Req)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			off := int(p.Parameters[i].Offset)

			oldVal, ok := Parameter39[off]
			if ok {
				if oldVal != p.Parameters[i].Value {
					for _, cb := range parameter39CB[off] {

						cb(off, p.Parameters[i].Value)
					}
				}
			} else {
				for _, cb := range parameter39CB[off] {
					cb(off, p.Parameters[i].Value)
				}
			}
			Parameter39[off] = p.Parameters[i].Value
		}

		return nil
	})

	RegisterPacketCallback(&PacketF036Req{}, func(value interface{}) error {
		p, ok := value.(*PacketF036Req)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			off := int(p.Parameters[i].Offset)
			oldVal, ok := Parameter36[off]
			if ok {
				if oldVal != p.Parameters[i].Value {
					for _, cb := range parameter36CB[off] {
						cb(off, p.Parameters[i].Value)
					}
				}
			} else {
				for _, cb := range parameter36CB[off] {
					cb(off, p.Parameters[i].Value)
				}
			}
			Parameter36[off] = p.Parameters[i].Value
		}

		return nil
	})
	RegisterPacketCallback(&PacketF03bReq{}, func(value interface{}) error {
		p, ok := value.(*PacketF03bReq)
		if !ok {
			return nil
		}
		for i := range p.Parameters {
			off := int(p.Parameters[i].Offset)
			oldVal, ok := Parameter3b[off]
			if ok {
				if oldVal != p.Parameters[i].Value {
					for _, cb := range parameter3bCB[off] {
						cb(off, p.Parameters[i].Value)
					}
				}
			} else {
				for _, cb := range parameter3bCB[off] {
					cb(off, p.Parameters[i].Value)
				}
			}
			Parameter3b[off] = p.Parameters[i].Value
		}

		return nil
	})
}
