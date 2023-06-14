package ppu

type Registers [12]uint8

const (
	BGDisplayBit = 0
	OBJEnableBit = 1
	OBJSizeBit   = 2
	BGMapBit     = 3
	BGTileBit    = 4
	WinEnableBit = 5
	WinMapBit    = 6
)

func (r *Registers) LCDC(bitPos int) bool {
	return (r[0] & (1 << bitPos)) != 0
}

func (r *Registers) SCY() uint8 {
	return r[2]
}

func (r *Registers) SetSCY(data uint8) {
	r[2] = data
}

func (r *Registers) SCX() uint8 {
	return r[3]
}

func (r *Registers) SetSCX(data uint8) {
	r[3] = data
}

func (r *Registers) LY() uint8 {
	return r[4]
}

func (r *Registers) SetLY(data uint8) {
	r[4] = data
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
func (r *Registers) BGP() uint8 {
	return r[7]
}

func (r *Registers) SetBGP(data uint8) {
	r[7] = data
}
