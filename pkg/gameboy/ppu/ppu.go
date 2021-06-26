package ppu

import (
	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/gameboy/cpu"
	"github.com/d2verb/gemu/pkg/gameboy/lcd"
	"github.com/d2verb/gemu/pkg/log"
)

const (
	CyclesPerScanLine = 456
	VBlankLines       = 10
)

type PPU struct {
	ioRegs   [12]uint8
	oam      [160]uint8
	ioRange  bus.AddressRange
	oamRange bus.AddressRange
	cycles   int
	l        *lcd.LCD
	bus      *bus.Bus
}

func New(l *lcd.LCD) *PPU {
	return &PPU{
		ioRange:  bus.NewAddressRange(0xff40, 0xff4b),
		oamRange: bus.NewAddressRange(0xfe00, 0xfe9f),
		l:        l,
	}
}

func (p *PPU) Step(cycles int) {
	p.cycles += cycles

	if p.cycles < CyclesPerScanLine && p.LY() < lcd.ScreenHeight {
		if p.cycles < 80 {
			// OAM Search
		} else if p.cycles < 252 {
			// Pixel Transfer
			p.buildBackground()
		} else {
			// HBlank
		}
	}

	p.SetLY((p.LY() + 1) % (lcd.ScreenHeight + VBlankLines))
	if p.LY() == lcd.ScreenHeight {
		p.bus.SetIF(cpu.IntVBlank)
		p.l.Updated <- nil
	}

	if p.cycles >= CyclesPerScanLine {
		p.cycles -= CyclesPerScanLine
	}
}

func (p *PPU) ConnectToBus(b *bus.Bus) error {
	if err := b.Map(p.ioRange, p); err != nil {
		return err
	}
	if err := b.Map(p.oamRange, p); err != nil {
		return err
	}
	p.bus = b
	return nil
}

func (p *PPU) Read8(address uint16) uint8 {
	if p.ioRange.Contains(address) {
		offset := address - p.ioRange.Start
		return p.ioRegs[offset]
	} else if p.oamRange.Contains(address) {
		offset := address - p.oamRange.Start
		return p.oam[offset]
	} else {
		log.Fatalf("PPU cannot be accessed at 0x%04x", address)
	}
	return 0
}

func (p *PPU) Read16(address uint16) uint16 {
	loByte := p.Read8(address)
	hiByte := p.Read8(address + 1)
	return ((uint16)(hiByte)<<8 | (uint16)(loByte))
}

func (p *PPU) Write8(address uint16, data uint8) {
	if p.ioRange.Contains(address) {
		offset := address - p.ioRange.Start
		p.ioRegs[offset] = data
	} else if p.oamRange.Contains(address) {
		offset := address - p.oamRange.Start
		p.oam[offset] = data
	} else {
		log.Fatalf("PPU cannot be accessed at 0x%04x", address)
	}
}

func (p *PPU) Write16(address uint16, data uint16) {
	hiByte := (uint8)((data >> 8) & 0xff)
	loByte := (uint8)(data & 0xff)

	p.Write8(address, loByte)
	p.Write8(address+1, hiByte)
}

func (p *PPU) buildBackground() {
	var x uint8 = 0
	for ; x < lcd.ScreenWidth; x++ {
		tileX := uint16((x + p.SCX()) / 8)
		tileY := uint16((p.LY() + p.SCY()) / 8)

		tileID := p.BGMap(tileY*32 + tileX)
		rawTileData := p.BGTiles(tileID)

		tile := p.constructTile(rawTileData)
		offsetX := x - x/8*8
		offsetY := p.LY() - p.LY()/8*8

		p.l.Lock()
		p.l.Screen[p.LY()][x] = tile[offsetY][offsetX] * 60
		p.l.Unlock()
	}
}

func (p *PPU) constructTile(rawTileData [16]uint8) [8][8]uint8 {
	tile := [8][8]uint8{}

	for i := 0; i < 16; i += 2 {
		loLine := rawTileData[i]
		hiLine := rawTileData[i+1]

		for j := 0; j < 8; j++ {
			paletteID := (hiLine >> j) & 1
			paletteID = (paletteID << 1) | ((loLine >> j) & 1)

			x := 7 - j
			y := i / 2

			tile[y][x] = p.BGColor(paletteID)
		}
	}

	return tile
}

func (p *PPU) BGColor(paletteID uint8) uint8 {
	var mask uint8 = 0b11 << (paletteID * 2)
	return (p.BGP() & mask >> (paletteID * 2)) & 0b11
}

func (p *PPU) BGTiles(tileID uint8) [16]uint8 {
	var baseAddress uint16 = 0x8800
	if (p.LCDC() & (1 << 4)) != 0 {
		baseAddress = 0x8000
	}
	baseAddress += uint16(tileID) * 16

	rawTile := [16]uint8{}
	for i := 0; i < 16; i++ {
		rawTile[i] = p.bus.Read8(baseAddress + uint16(i))
	}

	return rawTile
}

func (p *PPU) BGMap(offset uint16) uint8 {
	var baseAddress uint16 = 0x9800
	if (p.LCDC() & (1 << 3)) != 0 {
		baseAddress = 0x9c00
	}
	return p.bus.Read8(baseAddress + offset)
}

func (p *PPU) LY() uint8 {
	return p.ioRegs[4]
}

func (p *PPU) SetLY(data uint8) {
	p.ioRegs[4] = data
}

func (p *PPU) SCY() uint8 {
	return p.ioRegs[2]
}

func (p *PPU) SetSCY(data uint8) {
	p.ioRegs[2] = data
}

func (p *PPU) SCX() uint8 {
	return p.ioRegs[3]
}

func (p *PPU) SetSCX(data uint8) {
	p.ioRegs[3] = data
}

// BGP - BG Palette Data
//
// Bit 7-6 - Shade for Color Number 3
// Bit 5-4 - Shade for Color Number 2
// Bit 3-2 - Shade for Color Number 1
// Bit 1-0 - Shade for Color Number 0
//
// The four possible gray shades are:
//
// 0 White
// 1 Light Gray
// 2 Dark Gray
// 3 Black
func (p *PPU) BGP() uint8 {
	return p.ioRegs[7]
}

func (p *PPU) SetBGP(data uint8) {
	p.ioRegs[7] = data
}

func (p *PPU) LCDC() uint8 {
	return p.ioRegs[0]
}
