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
		A:  0,
		F:  0,
		B:  0,
		C:  0,
		D:  0,
		E:  0,
		H:  0,
		L:  0,
		SP: 0,
		PC: 0x100, // Skip Boot ROM and execute Game ROM directly
	}
}

const (
	ZFlag = 3
	NFlag = 2
	HFlag = 1
	CFlag = 0
)

func (r *Registers) Flag(bitPos int) bool {
	return ((r.F >> bitPos) & 1) == 1
}

func (r *Registers) SetFlag(bitPos int) {
	r.F |= (1 << bitPos)
}

func (r *Registers) UnsetFlag(bitPos int) {
	r.F &= ^(uint8)(1 << bitPos)
}
