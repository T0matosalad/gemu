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
			data, err := cpu.readOperandWord()
			if err != nil {
				return err
			}
			cpu.regs.SP = data
			return nil
		}),
		0xc3: newInstruction("jp a16", 16, func(cpu *CPU) error {
			address, err := cpu.readOperandWord()
			if err != nil {
				return err
			}
			cpu.regs.PC = address
			return nil
		}),
		0xf0: newInstruction("ldh A, (a8)", 12, func(cpu *CPU) error {
			offset, err := cpu.readOperandByte()
			if err != nil {
				return err
			}

			data, err := cpu.bus.ReadByte(0xff00 + (uint16)(offset))
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

func (c *CPU) readOperandByte() (uint8, error) {
	data, err := c.bus.ReadByte(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC++
	return data, nil
}

func (c *CPU) readOperandWord() (uint16, error) {
	data, err := c.bus.ReadWord(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC += 2
	return data, nil
}
