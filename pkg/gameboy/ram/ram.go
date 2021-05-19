package ram

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
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

func (r *RAM) ReadUInt8(address uint16) (uint8, error) {
	if r.vramRange.Contains(address) {
		offset := address - r.vramRange.Start
		return r.ram[offset], nil
	}

	if r.wramRange.Contains(address) {
		offset := address - r.wramRange.Start + 0x2000
		return r.ram[offset], nil
	}

	if r.eramRange.Contains(address) {
		offset := address - r.eramRange.Start
		return r.ram[offset], nil
	}

	if r.hramRange.Contains(address) {
		offset := address - r.hramRange.Start + 0x4000
		return r.ram[offset], nil
	}

	return 0, fmt.Errorf("RAM cannot be accessed at 0x%04x", address)
}

func (r *RAM) ReadUInt16(address uint16) (uint16, error) {
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

func (r *RAM) WriteUInt8(address uint16, data uint8) error {
	if r.vramRange.Contains(address) {
		offset := address - r.vramRange.Start
		r.ram[offset] = data
		return nil
	}

	if r.wramRange.Contains(address) {
		offset := address - r.wramRange.Start + 0x2000
		r.ram[offset] = data
		return nil
	}

	if r.eramRange.Contains(address) {
		offset := address - r.eramRange.Start
		r.ram[offset] = data
		return nil
	}

	if r.hramRange.Contains(address) {
		offset := address - r.hramRange.Start + 0x4000
		r.ram[offset] = data
		return nil
	}

	return fmt.Errorf("RAM cannot be accessed at 0x%04x", address)
}

func (r *RAM) WriteUInt16(address uint16, data uint16) error {
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
