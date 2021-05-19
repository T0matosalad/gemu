package rom

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
)

type ROM struct {
	m MBC
}

func New(data []uint8) (*ROM, error) {
	mbcType := data[0x147]

	switch mbcType {
	case 0:
		return &ROM{
			m: NewMBC0(data),
		}, nil
	default:
		return nil, fmt.Errorf("MBC type %d is not supported", mbcType)
	}
}

func (r *ROM) String() string {
	return fmt.Sprintf("Title: %s, MBCType: %d", r.Title(), r.MBCType())
}

func (r *ROM) Title() []uint8 {
	start, end := 0x134, 0x134
	for ; end < 0x144 && r.m.Data()[end] != 0; end++ {
	}
	return r.m.Data()[start:end]
}

func (r *ROM) MBCType() uint8 {
	return r.m.Data()[0x147]
}

func (r *ROM) ConnectToBus(b *bus.Bus) error {
	for _, _range := range r.m.AddressRanges() {
		if err := b.Map(_range, r); err != nil {
			return err
		}
	}
	return nil
}

func (r *ROM) ReadUInt8(address uint16) (uint8, error) {
	return r.m.ReadUInt8(address)
}

func (r *ROM) ReadUInt16(address uint16) (uint16, error) {
	loByte, err := r.ReadUInt8(address)
	if err != nil {
		return 0, err
	}

	hiByte, err := r.ReadUInt8(address + 1)
	if err != nil {
		return 0, err
	}

	return ((uint16)(hiByte)<<8 | (uint16)(loByte)), nil
}

func (r *ROM) WriteUInt8(address uint16, data uint8) error {
	return r.m.WriteUInt8(address, data)
}

func (r *ROM) WriteUInt16(address uint16, data uint16) error {
	hiByte := (uint8)((data >> 8) & 0xff)
	loByte := (uint8)(data & 0xff)

	if err := r.WriteUInt8(address, loByte); err != nil {
		return err
	}

	if err := r.WriteUInt8(address+1, hiByte); err != nil {
		return err
	}

	return nil
}
