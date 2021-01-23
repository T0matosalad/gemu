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
		0x31: newInstruction("ld SP, d16", 12, func(cpu *CPU) error {
			data, err := cpu.bus.ReadWord(cpu.regs.PC)
			if err != nil {
				return err
			}
			cpu.regs.PC += 2
			cpu.regs.SP = data
			return nil
		}),
		0xc3: newInstruction("jp a16", 16, func(cpu *CPU) error {
			address, err := cpu.bus.ReadWord(cpu.regs.PC)
			if err != nil {
				return err
			}
			cpu.regs.PC = address
			return nil
		}),
		0xf0: newInstruction("ldh A, (a8)", 12, func(cpu *CPU) error {
			offset, err := cpu.bus.ReadByte(cpu.regs.PC)
			if err != nil {
				return err
			}
			cpu.regs.PC++

			address := 0xff00 + (uint16)(offset)
			data, err := cpu.bus.ReadByte(address)
			if err != nil {
				return err
			}

			cpu.regs.A = data
			return nil
		}),
		0xf3: newInstruction("di", 4, func(cpu *CPU) error {
			cpu.ime = false
			return nil
		}),
	}
}
