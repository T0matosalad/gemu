package rom

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
)

type MBC0 struct {
	data      []uint8
	bankRange bus.AddressRange
}

func NewMBC0(data []uint8) MBC {
	return &MBC0{
		data:      data,
		bankRange: bus.NewAddressRange(0x0000, 0x7fff),
	}
}

func (m *MBC0) AddressRanges() []bus.AddressRange {
	return []bus.AddressRange{
		m.bankRange,
	}
}

func (m *MBC0) Data() []uint8 {
	return m.data
}

func (m *MBC0) ReadByte(address uint16) (uint8, error) {
	if m.bankRange.Contains(address) {
		offset := address - m.bankRange.Start
		return m.data[offset], nil
	}

	return 0, fmt.Errorf("ROM cannot be accessed at 0x%04x", address)
}

func (m *MBC0) WriteByte(address uint16, data uint8) error {
	if m.bankRange.Contains(address) {
		offset := address - m.bankRange.Start
		m.data[offset] = data
		return nil
	}

	return fmt.Errorf("ROM cannot be accessed at 0x%04x", address)
}
