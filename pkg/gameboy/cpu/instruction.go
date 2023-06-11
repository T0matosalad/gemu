package cpu

type instruction struct {
	mnemonic string
	handler  func(cpu *CPU) int
}

func newInstruction(mnemonic string, handler func(cpu *CPU) int) instruction {
	return instruction{
		mnemonic: mnemonic,
		handler:  handler,
	}
}

func newInstructionSet() map[uint16]instruction {
	return map[uint16]instruction{
		0x00: newInstruction("nop", func(cpu *CPU) int {
			return 4
		}),
		0x01: newInstruction("ld BC, d16", func(cpu *CPU) int {
			cpu.regs.SetBC(cpu.operand16())
			return 12
		}),
		0x04: newInstruction("inc B", func(cpu *CPU) int {
			cpu.regs.B = cpu.add8(cpu.regs.B, 1, false)
			return 4
		}),
		0x05: newInstruction("dec B", func(cpu *CPU) int {
			cpu.regs.B = cpu.sub8(cpu.regs.B, 1, false)
			return 4
		}),
		0x06: newInstruction("ld B, d8", func(cpu *CPU) int {
			cpu.regs.B = cpu.operand8()
			return 8
		}),
		0x0b: newInstruction("dec BC", func(cpu *CPU) int {
			cpu.regs.SetBC(cpu.regs.BC() - 1)
			return 8
		}),
		0x0c: newInstruction("inc C", func(cpu *CPU) int {
			cpu.regs.C = cpu.add8(cpu.regs.C, 1, false)
			return 4
		}),
		0x0d: newInstruction("dec C", func(cpu *CPU) int {
			cpu.regs.C = cpu.sub8(cpu.regs.C, 1, false)
			return 4
		}),
		0x0e: newInstruction("ld C, d8", func(cpu *CPU) int {
			cpu.regs.C = cpu.operand8()
			return 12
		}),
		0x11: newInstruction("ld DE, d16", func(cpu *CPU) int {
			cpu.regs.SetDE(cpu.operand16())
			return 12
		}),
		0x12: newInstruction("ld (DE), A", func(cpu *CPU) int {
			cpu.bus.Write8(cpu.regs.DE(), cpu.regs.A)
			return 8
		}),
		0x13: newInstruction("inc DE", func(cpu *CPU) int {
			cpu.regs.SetDE(cpu.regs.DE() + 1)
			return 8
		}),
		0x18: newInstruction("jr r8", func(cpu *CPU) int {
			cpu.regs.PC += signExtU8ToU16(cpu.operand8())
			return 12
		}),
		0x1a: newInstruction("ld A, (DE)", func(cpu *CPU) int {
			cpu.regs.A = cpu.bus.Read8(cpu.regs.DE())
			return 8
		}),
		0x20: newInstruction("jr nz, r8", func(cpu *CPU) int {
			data := cpu.operand8()
			if cpu.regs.Flag(ZFlag) {
				return 8
			}
			cpu.regs.PC += signExtU8ToU16(data)
			return 12
		}),
		0x21: newInstruction("ld HL, d16", func(cpu *CPU) int {
			cpu.regs.SetHL(cpu.operand16())
			return 12
		}),
		0x22: newInstruction("ld (HL+), A", func(cpu *CPU) int {
			cpu.bus.Write8(cpu.regs.HL(), cpu.regs.A)
			cpu.regs.SetHL(cpu.regs.HL() + 1)
			return 8
		}),
		0x28: newInstruction("jr Z, r8", func(cpu *CPU) int {
			data := cpu.operand8()
			if cpu.regs.Flag(ZFlag) {
				cpu.regs.PC += signExtU8ToU16(data)
				return 12
			}
			return 8
		}),
		0x2a: newInstruction("ld A, (HL+)", func(cpu *CPU) int {
			cpu.regs.A = cpu.bus.Read8(cpu.regs.HL())
			cpu.regs.SetHL(cpu.regs.HL() + 1)
			return 8
		}),
		0x31: newInstruction("ld SP, d16", func(cpu *CPU) int {
			cpu.regs.SP = cpu.operand16()
			return 12
		}),
		0x38: newInstruction("jr C, r8", func(cpu *CPU) int {
			data := cpu.operand8()
			if cpu.regs.Flag(CFlag) {
				cpu.regs.PC += signExtU8ToU16(data)
				return 12
			}
			return 8
		}),
		0x3e: newInstruction("ld A, d8", func(cpu *CPU) int {
			cpu.regs.A = cpu.operand8()
			return 8
		}),
		0x76: newInstruction("halt", func(cpu *CPU) int {
			cpu.halt = true
			return 4
		}),
		0x78: newInstruction("ld A, B", func(cpu *CPU) int {
			cpu.regs.A = cpu.regs.B
			return 4
		}),
		0xaf: newInstruction("xor A", func(cpu *CPU) int {
			cpu.regs.A = 0
			cpu.regs.SetFlag(CFlag, false)
			cpu.regs.SetFlag(HFlag, false)
			cpu.regs.SetFlag(NFlag, false)
			cpu.regs.SetFlag(ZFlag, true)
			return 4
		}),
		0xa7: newInstruction("and A", func(cpu *CPU) int {
			cpu.regs.A = cpu.and8(cpu.regs.A, cpu.regs.A)
			return 4
		}),
		0xb1: newInstruction("or C", func(cpu *CPU) int {
			cpu.regs.A = cpu.or8(cpu.regs.A, cpu.regs.C)
			return 4
		}),
		0xc3: newInstruction("jp a16", func(cpu *CPU) int {
			cpu.regs.PC = cpu.operand16()
			return 16
		}),
		0xc9: newInstruction("ret", func(cpu *CPU) int {
			cpu.regs.PC = cpu.bus.Read16(cpu.regs.SP)
			cpu.regs.SP += 2
			return 16
		}),
		0xcd: newInstruction("call d16", func(cpu *CPU) int {
			cpu.regs.SP -= 2
			address := cpu.operand16()
			cpu.bus.Write16(cpu.regs.SP, cpu.regs.PC)
			cpu.regs.PC = address
			return 24
		}),
		0xe0: newInstruction("ldh (a8), A", func(cpu *CPU) int {
			offset := cpu.operand8()
			cpu.bus.Write8(0xff00+(uint16)(offset), cpu.regs.A)
			return 12
		}),
		0xee: newInstruction("xor d8", func(cpu *CPU) int {
			cpu.regs.A = cpu.xor8(cpu.regs.A, cpu.operand8())
			return 8
		}),
		0xf0: newInstruction("ldh A, (a8)", func(cpu *CPU) int {
			offset := cpu.operand8()
			cpu.regs.A = cpu.bus.Read8(0xff00 + (uint16)(offset))
			return 12
		}),
		0xfe: newInstruction("cp d8", func(cpu *CPU) int {
			cpu.sub8(cpu.regs.A, cpu.operand8(), true)
			return 8
		}),
		0xf3: newInstruction("di", func(cpu *CPU) int {
			cpu.ime = false
			return 4
		}),
		0xfb: newInstruction("ei", func(cpu *CPU) int {
			cpu.ime = true
			return 4
		}),
	}
}

func (c *CPU) operand8() uint8 {
	data := c.bus.Read8(c.regs.PC)
	c.regs.PC++
	return data
}

func (c *CPU) operand16() uint16 {
	data := c.bus.Read16(c.regs.PC)
	c.regs.PC += 2
	return data
}

func (c *CPU) add8(a uint8, b uint8, updateCFlag bool) uint8 {
	result := a + b

	c.regs.SetFlag(NFlag, false)
	c.regs.SetFlag(ZFlag, result == 0)
	c.regs.SetFlag(HFlag, (a&0xf)+(b&0xf) > 0xf)

	if updateCFlag {
		c.regs.SetFlag(CFlag, result <= a)
	}

	return result
}

func (c *CPU) sub8(a uint8, b uint8, updateCFlag bool) uint8 {
	result := a - b

	c.regs.SetFlag(NFlag, true)
	c.regs.SetFlag(ZFlag, result == 0)
	c.regs.SetFlag(HFlag, (a&0xf) < (b&0xf))

	if updateCFlag {
		c.regs.SetFlag(CFlag, a < b)
	}

	return result
}

func (c *CPU) and8(a uint8, b uint8) uint8 {
	result := a & b

	c.regs.SetFlag(NFlag, false)
	c.regs.SetFlag(ZFlag, result == 0)
	c.regs.SetFlag(HFlag, true)
	c.regs.SetFlag(CFlag, false)

	return result
}

func (c *CPU) or8(a uint8, b uint8) uint8 {
	result := a | b

	c.regs.SetFlag(NFlag, false)
	c.regs.SetFlag(ZFlag, result == 0)
	c.regs.SetFlag(HFlag, false)
	c.regs.SetFlag(CFlag, false)

	return result
}

func (c *CPU) xor8(a uint8, b uint8) uint8 {
	result := a ^ b

	c.regs.SetFlag(NFlag, false)
	c.regs.SetFlag(ZFlag, result == 0)
	c.regs.SetFlag(HFlag, false)
	c.regs.SetFlag(CFlag, false)

	return result
}

// Example: 0b1111_1010 to 0b1111_1111_1111_1010
func signExtU8ToU16(from uint8) uint16 {
	to := (uint16)(from)
	if (from & 0b1000_0000) != 0 {
		to |= 0xff00
	}
	return to
}
