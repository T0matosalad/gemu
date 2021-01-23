package ppu

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
)

const (
	CyclesPerScanLine = 456
	ScreenHeight      = 144
	ScreenWidth       = 160
	VBlankLines       = 10
)

type PPU struct {
	ioRegs  [12]uint8
	ioRange bus.AddressRange
	cycles  int
}

func New() PPU {
	return PPU{
		ioRange: bus.NewAddressRange(0xff40, 0xff4b),
	}
}

func (p *PPU) Step(cycles int) error {
	p.cycles += cycles
	if p.cycles < CyclesPerScanLine {
		return nil
	}
	p.SetLY((p.LY() + 1) % (ScreenHeight + VBlankLines))
	p.cycles -= CyclesPerScanLine
	return nil
}

func (p *PPU) ConnectToBus(b *bus.Bus) error {
	return b.Map(p.ioRange, p)
}

func (p *PPU) ReadByte(address uint16) (uint8, error) {
	if p.ioRange.Contains(address) {
		offset := address - p.ioRange.Start
		return p.ioRegs[offset], nil
	}

	return 0, fmt.Errorf("PPU cannot be accessed at 0x%04x", address)
}

func (p *PPU) ReadWord(address uint16) (uint16, error) {
	loByte, err := p.ReadByte(address)
	if err != nil {
		return 0, err
	}

	hiByte, err := p.ReadByte(address + 1)
	if err != nil {
		return 0, err
	}

	return ((uint16)(hiByte)<<8 | (uint16)(loByte)), nil
}

func (p *PPU) WriteByte(address uint16, data uint8) error {
	if p.ioRange.Contains(address) {
		offset := address - p.ioRange.Start
		p.ioRegs[offset] = data
		return nil
	}

	return fmt.Errorf("PPU cannot be accessed at 0x%04x", address)
}

func (p *PPU) WriteWord(address uint16, data uint16) error {
	hiByte := (uint8)((data >> 8) & 0xff)
	loByte := (uint8)(data & 0xff)

	if err := p.WriteByte(address, loByte); err != nil {
		return err
	}

	if err := p.WriteByte(address+1, hiByte); err != nil {
		return err
	}

	return nil
}

func (p *PPU) LY() uint8 {
	return p.ioRegs[4]
}

func (p *PPU) SetLY(data uint8) {
	p.ioRegs[4] = data
}
