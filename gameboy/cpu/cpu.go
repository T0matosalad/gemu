package cpu

type CPU struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	F  uint8
	PC uint16
	SP uint16
}

func New() CPU {
	return CPU{
		A:  0,
		B:  0,
		C:  0,
		D:  0,
		E:  0,
		H:  0,
		L:  0,
		F:  0,
		PC: 0x100,
		SP: 0,
	}
}
