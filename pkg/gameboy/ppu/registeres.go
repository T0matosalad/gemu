package ppu

type Registers [12]uint8

// LCDC bit flags
const (
	BGDisplayFlag = 0b1
	OBJEnableFlag = 0b10
	OBJSizeFlag   = 0b100
	BGMapFlag     = 0b1000
	BGTileFlag    = 0b10000
	WinEnableFlag = 0b100000
	WinMapFlag    = 0b1000000
	LCDEnableFlag = 0b10000000
)

// LCD status bit flags
const (
	ModeFlag = 0b11
)

// LCD mode types
const (
	HBlankMode        = 0
	VBlankMode        = 1
	OAMSearchMode     = 2
	PixelTransferMode = 3
)

func (r *Registers) LCDC(mask uint8) uint8 {
	return r[0] & mask
}

func (r *Registers) STAT(mask uint8) uint8 {
	return r[1] & mask
}

func (r *Registers) SetSTAT(mask uint8, on bool) {
	if on {
		r[1] |= mask
	} else {
		r[1] &= ^mask
	}
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
