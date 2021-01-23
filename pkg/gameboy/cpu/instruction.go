package cpu

type instruction struct {
	mnemonic string
	cycles   int
	handler  func(cpu *CPU) error
}

func newInstruction(mnemonic string, cycles int, handler func(cpu *CPU) error) instruction {
	return instruction{
		mnemonic: mnemonic,
		cycles:   cycles,
		handler:  handler,
	}
}

func newInstructionSet() map[uint16]instruction {
	return map[uint16]instruction{
		0x00: newInstruction("nop", 4, func(cpu *CPU) error {
			return nil
		}),
		0xc3: newInstruction("jp nn", 16, func(cpu *CPU) error {
			address, err := cpu.bus.ReadWord(cpu.regs.PC)
			if err != nil {
				return err
			}
			cpu.regs.PC = address
			return nil
		}),
	}
}
