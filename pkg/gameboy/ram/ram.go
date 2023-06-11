package ram

import (
	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/log"
)

type RAM struct {
	ram       [0x407f]uint8 // 8KB Video RAM and 8KB Work RAM
	vramRange bus.AddressRange
	wramRange bus.AddressRange
	eramRange bus.AddressRange
	hramRange bus.AddressRange
}

func New() *RAM {
	return &RAM{
		vramRange: bus.NewAddressRange(0x8000, 0x9fff),
		wramRange: bus.NewAddressRange(0xc000, 0xdfff),
		eramRange: bus.NewAddressRange(0xe000, 0xfdff),
		hramRange: bus.NewAddressRange(0xff80, 0xfffe),
	}
}

func (r *RAM) ConnectToBus(b *bus.Bus) error {
	if err := b.Map(r.vramRange, r); err != nil {
		return err
	}
	if err := b.Map(r.wramRange, r); err != nil {
		return err
	}
	if err := b.Map(r.eramRange, r); err != nil {
		return err
	}
	if err := b.Map(r.hramRange, r); err != nil {
		return err
	}
	return nil
}

func (r *RAM) Read8(address uint16) uint8 {
	if r.vramRange.Contains(address) {
		offset := address - r.vramRange.Start
		return r.ram[offset]
	} else if r.wramRange.Contains(address) {
		offset := address - r.wramRange.Start + 0x2000
		return r.ram[offset]
	} else if r.eramRange.Contains(address) {
		offset := address - r.eramRange.Start
		return r.ram[offset]
	} else if r.hramRange.Contains(address) {
		offset := address - r.hramRange.Start + 0x4000
		return r.ram[offset]
	} else {
		log.Fatalf("RAM cannot be accessed at 0x%04x", address)
	}
	return 0
}

func (r *RAM) Read16(address uint16) uint16 {
	loByte := r.Read8(address)
	hiByte := r.Read8(address + 1)
	return ((uint16)(hiByte)<<8 | (uint16)(loByte))
}

func (r *RAM) Write8(address uint16, data uint8) {
	if r.vramRange.Contains(address) {
		offset := address - r.vramRange.Start
		r.ram[offset] = data
	} else if r.wramRange.Contains(address) {
		offset := address - r.wramRange.Start + 0x2000
		r.ram[offset] = data
	} else if r.eramRange.Contains(address) {
		offset := address - r.eramRange.Start
		r.ram[offset] = data
	} else if r.hramRange.Contains(address) {
		offset := address - r.hramRange.Start + 0x4000
		r.ram[offset] = data
	} else {
		log.Fatalf("RAM cannot be accessed at 0x%04x", address)
	}
}

func (r *RAM) Write16(address uint16, data uint16) {
	hiByte := (uint8)(data >> 8)
	loByte := (uint8)(data & 0xff)

	r.Write8(address, loByte)
	r.Write8(address+1, hiByte)
}
