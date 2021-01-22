package cpu

import (
	"fmt"

	"github.com/d2verb/gemu/gameboy/bus"
)

type CPU struct {
	Regs           Registers
	bus            *bus.Bus
	instructionSet map[uint16]instruction
}

func New() CPU {
	return CPU{
		Regs:           newRegisters(),
		instructionSet: newInstructionSet(),
	}
}

func (c *CPU) ConnectToBus(bus *bus.Bus) error {
	c.bus = bus
	return nil
}

func (c *CPU) Tick() (int, error) {
	opcode, err := c.fetch()
	if err != nil {
		return 0, err
	}

	instruction, ok := c.instructionSet[opcode]
	if !ok {
		return 0, fmt.Errorf("Unknown opcode 0x%x (PC: 0x%x)", opcode, c.Regs.PC)
	}

	if err := instruction.handler(c); err != nil {
		return 0, err
	}

	return instruction.cycles, nil
}

func (c *CPU) fetch() (uint16, error) {
	opcode, err := c.bus.ReadByte(c.Regs.PC)
	if err != nil {
		return 0, err
	}
	c.Regs.PC++
	return (uint16)(opcode), nil
}
