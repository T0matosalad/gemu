package rom

import (
	"log"

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

func (m *MBC0) Read8(address uint16) uint8 {
	if m.bankRange.Contains(address) {
		offset := address - m.bankRange.Start
		return m.data[offset]
	} else {
		log.Fatalf("ROM cannot be accessed at 0x%04x", address)
	}
	return 0
}

func (m *MBC0) Write8(address uint16, data uint8) {
	if m.bankRange.Contains(address) {
		offset := address - m.bankRange.Start
		m.data[offset] = data
	} else {
		log.Fatalf("ROM cannot be accessed at 0x%04x", address)
	}
}
