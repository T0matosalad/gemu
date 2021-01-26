package rom

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
)

type ROM struct {
	data       []uint8
	bank0Range bus.AddressRange
}

func New(data []uint8) ROM {
	return ROM{
		data:       data,
		bank0Range: bus.NewAddressRange(0x0000, 0x3fff),
	}
}

func (r *ROM) Title() []uint8 {
	start, end := 0x134, 0x134
	for ; end < 0x144 && r.data[end] != 0; end++ {
	}
	return r.data[start:end]
}

func (r *ROM) MBCType() uint8 {
	return r.data[0x147]
}

func (r *ROM) ConnectToBus(b *bus.Bus) error {
	return b.Map(r.bank0Range, r)
}

func (r *ROM) ReadByte(address uint16) (uint8, error) {
	if r.bank0Range.Contains(address) {
		offset := address - r.bank0Range.Start
		return r.data[offset], nil
	}

	return 0, fmt.Errorf("ROM cannot be accessed at 0x%04x", address)
}

func (r *ROM) ReadWord(address uint16) (uint16, error) {
	loByte, err := r.ReadByte(address)
	if err != nil {
		return 0, err
	}

	hiByte, err := r.ReadByte(address + 1)
	if err != nil {
		return 0, err
	}

	return ((uint16)(hiByte)<<8 | (uint16)(loByte)), nil
}

func (r *ROM) WriteByte(address uint16, data uint8) error {
	if r.bank0Range.Contains(address) {
		offset := address - r.bank0Range.Start
		r.data[offset] = data
		return nil
	}

	return fmt.Errorf("ROM cannot be accessed at 0x%04x", address)
}

func (r *ROM) WriteWord(address uint16, data uint16) error {
	hiByte := (uint8)((data >> 8) & 0xff)
	loByte := (uint8)(data & 0xff)

	if err := r.WriteByte(address, loByte); err != nil {
		return err
	}

	if err := r.WriteByte(address+1, hiByte); err != nil {
		return err
	}

	return nil
}
