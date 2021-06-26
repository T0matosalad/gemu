package cpu

type instruction struct {
	mnemonic string
	handler  func(cpu *CPU) (int, error)
}

func newInstruction(mnemonic string, handler func(cpu *CPU) (int, error)) instruction {
	return instruction{
		mnemonic: mnemonic,
		handler:  handler,
	}
}

func newInstructionSet() map[uint16]instruction {
	return map[uint16]instruction{
		0x00: newInstruction("nop", func(cpu *CPU) (int, error) {
			return 4, nil
		}),
		0x01: newInstruction("ld BC, d16", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			cpu.regs.SetBC(data)
			return 12, nil
		}),
		0x04: newInstruction("inc B", func(cpu *CPU) (int, error) {
			cpu.regs.B = cpu.addByte(cpu.regs.B, 1, false)
			return 4, nil
		}),
		0x05: newInstruction("dec B", func(cpu *CPU) (int, error) {
			cpu.regs.B = cpu.subByte(cpu.regs.B, 1, false)
			return 4, nil
		}),
		0x06: newInstruction("ld B, d8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, nil
			}
			cpu.regs.B = data
			return 8, nil
		}),
		0x0c: newInstruction("inc C", func(cpu *CPU) (int, error) {
			cpu.regs.C = cpu.addByte(cpu.regs.C, 1, false)
			return 4, nil
		}),
		0x0d: newInstruction("dec C", func(cpu *CPU) (int, error) {
			cpu.regs.C = cpu.subByte(cpu.regs.C, 1, false)
			return 4, nil
		}),
		0x0e: newInstruction("ld C, d8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}

			cpu.regs.C = data
			return 12, nil
		}),
		0x11: newInstruction("ld DE, d16", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			cpu.regs.SetDE(data)
			return 12, nil
		}),
		0x12: newInstruction("ld (DE), A", func(cpu *CPU) (int, error) {
			if err := cpu.bus.WriteUInt8(cpu.regs.DE(), cpu.regs.A); err != nil {
				return 0, err
			}
			return 8, nil
		}),
		0x13: newInstruction("inc DE", func(cpu *CPU) (int, error) {
			cpu.regs.SetDE(cpu.regs.DE() + 1)
			return 8, nil
		}),
		0x18: newInstruction("jr r8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}

			cpu.regs.PC += signExtU8ToU16(data)
			return 12, nil
		}),
		0x20: newInstruction("jr nz, r8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}

			if cpu.regs.Flag(ZFlag) {
				return 8, nil
			}

			cpu.regs.PC += signExtU8ToU16(data)
			return 12, nil
		}),
		0x21: newInstruction("ld HL, d16", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			cpu.regs.SetHL(data)
			return 12, nil
		}),
		0x22: newInstruction("ld (HL+), A", func(cpu *CPU) (int, error) {
			data, err := cpu.bus.ReadUInt8(cpu.regs.HL())
			if err != nil {
				return 0, err
			}
			cpu.regs.A = data
			cpu.regs.SetHL(cpu.regs.HL() + 1)
			return 8, nil
		}),
		0x2a: newInstruction("ld A, (HL+)", func(cpu *CPU) (int, error) {
			if err := cpu.bus.WriteUInt8(cpu.regs.HL(), cpu.regs.A); err != nil {
				return 0, err
			}
			cpu.regs.SetHL(cpu.regs.HL() + 1)
			return 8, nil
		}),
		0x31: newInstruction("ld SP, d16", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			cpu.regs.SP = data
			return 12, nil
		}),
		0x3e: newInstruction("ld A, d8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}
			cpu.regs.A = data
			return 8, nil
		}),
		0x76: newInstruction("halt", func(cpu *CPU) (int, error) {
			cpu.halt = true
			return 4, nil
		}),
		0xaf: newInstruction("xor A", func(cpu *CPU) (int, error) {
			cpu.regs.A = 0
			cpu.regs.UnsetFlag(CFlag)
			cpu.regs.UnsetFlag(HFlag)
			cpu.regs.UnsetFlag(NFlag)
			cpu.regs.SetFlag(ZFlag)
			return 4, nil
		}),
		0xc3: newInstruction("jp a16", func(cpu *CPU) (int, error) {
			address, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			cpu.regs.PC = address
			return 16, nil
		}),
		0xc9: newInstruction("ret", func(cpu *CPU) (int, error) {
			address, err := cpu.bus.ReadUInt16(cpu.regs.SP)
			if err != nil {
				return 0, err
			}
			cpu.regs.PC = address
			cpu.regs.SP += 2
			return 16, nil
		}),
		0xcd: newInstruction("call d16", func(cpu *CPU) (int, error) {
			cpu.regs.SP -= 2
			address, err := cpu.readOperandWord()
			if err != nil {
				return 0, err
			}
			if err := cpu.bus.WriteUInt16(cpu.regs.SP, cpu.regs.PC); err != nil {
				return 0, err
			}
			cpu.regs.PC = address
			return 24, nil
		}),
		0xe0: newInstruction("ldh (a8), A", func(cpu *CPU) (int, error) {
			offset, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}

			if err := cpu.bus.WriteUInt8(0xff00+(uint16)(offset), cpu.regs.A); err != nil {
				return 0, err
			}

			return 12, nil
		}),
		0xf0: newInstruction("ldh A, (a8)", func(cpu *CPU) (int, error) {
			offset, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}

			data, err := cpu.bus.ReadUInt8(0xff00 + (uint16)(offset))
			if err != nil {
				return 0, err
			}

			cpu.regs.A = data
			return 12, nil
		}),
		0xfe: newInstruction("cp d8", func(cpu *CPU) (int, error) {
			data, err := cpu.readOperandByte()
			if err != nil {
				return 0, err
			}
			cpu.subByte(cpu.regs.A, data, true)
			return 8, nil
		}),
		0xf3: newInstruction("di", func(cpu *CPU) (int, error) {
			cpu.ime = false
			return 4, nil
		}),
		0xfb: newInstruction("ei", func(cpu *CPU) (int, error) {
			cpu.ime = true
			return 4, nil
		}),
	}
}

func (c *CPU) readOperandByte() (uint8, error) {
	data, err := c.bus.ReadUInt8(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC++
	return data, nil
}

func (c *CPU) readOperandWord() (uint16, error) {
	data, err := c.bus.ReadUInt16(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC += 2
	return data, nil
}

// [TODO] update HFlag
func (c *CPU) addByte(a uint8, b uint8, updateCFlag bool) uint8 {
	result := a + b

	c.regs.UnsetFlag(NFlag)

	if result == 0 {
		c.regs.SetFlag(ZFlag)
	} else {
		c.regs.UnsetFlag(ZFlag)
	}

	if !updateCFlag {
		return result
	}

	if a+b <= a {
		c.regs.SetFlag(CFlag)
	} else {
		c.regs.UnsetFlag(CFlag)
	}

	return result
}

// [TODO] update HFlag
func (c *CPU) subByte(a uint8, b uint8, updateCFlag bool) uint8 {
	result := a - b

	c.regs.SetFlag(NFlag)

	if result == 0 {
		c.regs.SetFlag(ZFlag)
	} else {
		c.regs.UnsetFlag(ZFlag)
	}

	if !updateCFlag {
		return result
	}

	if a < b {
		c.regs.SetFlag(CFlag)
	} else {
		c.regs.UnsetFlag(CFlag)
	}

	return result
}

// Example: 0b1111_1010 to 0b1111_1111_1111_1010
func signExtU8ToU16(from uint8) uint16 {
	to := (uint16)(from)
	if (from & 0b10000000) != 0 {
		to |= 0xff00
	}
	return to
}
