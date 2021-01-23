package cpu

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/log"
)

type CPU struct {
	regs           Registers
	ime            bool
	bus            *bus.Bus
	instructionSet map[uint16]instruction
}

func New() CPU {
	return CPU{
		regs:           newRegisters(),
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
		return 0, fmt.Errorf("Unknown opcode 0x%x (PC: 0x%04x)", opcode, c.regs.PC)
	}

	log.Debugf("cpu (PC: 0x%04x): %s\n", c.regs.PC, instruction.mnemonic)

	cycles, err := instruction.handler(c)
	if err != nil {
		return 0, err
	}

	return cycles, nil
}

func (c *CPU) fetch() (uint16, error) {
	opcode, err := c.bus.ReadByte(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC++
	return (uint16)(opcode), nil
}
