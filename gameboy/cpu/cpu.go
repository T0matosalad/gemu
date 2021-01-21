package cpu

type CPU struct {
	Regs Registers
}

func New() CPU {
	return CPU{
		Regs: newRegisters(),
	}
}
