package cpu

type Registers struct {
	A  uint8
	F  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16
}

func newRegisters() Registers {
	return Registers{
		// A and F can be accessed as 16-bit register AF
		A: 0,
		F: 0,

		// B and C can be accessed as 16-bit register BC
		B: 0,
		C: 0,

		// D and E can be accessed as 16-bit register DE
		D: 0,
		E: 0,

		// H and L can be accessed as 16-bit register HL
		H: 0,
		L: 0,

		SP: 0,
		PC: 0x100, // Skip Boot ROM and execute Game ROM directly
	}
}

// Bit position of each flag in F register
const (
	CFlag = 0b1
	HFlag = 0b10
	NFlag = 0b100
	ZFlag = 0b1000
)

func (r *Registers) HL() uint16 {
	return ((uint16)(r.H) << 8) | (uint16)(r.L)
}

func (r *Registers) SetHL(data uint16) {
	r.H = (uint8)((data >> 8) & 0xff)
	r.L = (uint8)(data & 0xff)
}

func (r *Registers) BC() uint16 {
	return ((uint16)(r.B) << 8) | (uint16)(r.C)
}

func (r *Registers) SetBC(data uint16) {
	r.B = (uint8)((data >> 8) & 0xff)
	r.C = (uint8)(data & 0xff)
}

func (r *Registers) DE() uint16 {
	return ((uint16)(r.D) << 8) | (uint16)(r.E)
}

func (r *Registers) SetDE(data uint16) {
	r.D = (uint8)((data >> 8) & 0xff)
	r.E = (uint8)(data & 0xff)
}

func (r *Registers) Flag(mask uint8) uint8 {
	return r.F & mask
}

func (r *Registers) SetFlag(mask uint8, on bool) {
	if on {
		r.F |= mask
	} else {
		r.F &= ^mask
	}
}
