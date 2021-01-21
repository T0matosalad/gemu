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
