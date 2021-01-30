package cpu

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/log"
)

const Hz = 4194304

type CPU struct {
	regs           Registers
	ime            bool  // Interrupt Master Enable Flag
	ie             uint8 // Interrupt Enable
	_if            uint8 // Interrupt Flag
	bus            *bus.Bus
	instructionSet map[uint16]instruction
}

func New() *CPU {
	return &CPU{
		regs:           newRegisters(),
		instructionSet: newInstructionSet(),
	}
}

func (c *CPU) ConnectToBus(b *bus.Bus) error {
	if err := b.Map(bus.NewAddressRange(0xff0f, 0xff0f), c); err != nil {
		return err
	}
	if err := b.Map(bus.NewAddressRange(0xffff, 0xffff), c); err != nil {
		return err
	}
	c.bus = b
	return nil
}

func (c *CPU) Step() (int, error) {
	instAddr := c.regs.PC

	opcode, err := c.fetch()
	if err != nil {
		return 0, err
	}

	instruction, ok := c.instructionSet[opcode]
	if !ok {
		return 0, fmt.Errorf("Unknown opcode 0x%x (PC: 0x%04x)", opcode, instAddr)
	}

	log.Verbosef("(cpu) [0x%04x]: %s\n", instAddr, instruction.mnemonic)

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

func (c *CPU) ReadByte(address uint16) (uint8, error) {
	switch address {
	case 0xff0f:
		return c._if, nil
	case 0xffff:
		return c.ie, nil
	default:
		return 0, fmt.Errorf("CPU cannot be accessed at 0x%04x", address)
	}
}

func (c *CPU) WriteByte(address uint16, data uint8) error {
	switch address {
	case 0xff0f:
		c._if = data
		return nil
	case 0xffff:
		c.ie = data
		return nil
	default:
		return fmt.Errorf("CPU cannot be accessed at 0x%04x", address)
	}
}

func (c *CPU) ReadWord(address uint16) (uint16, error) {
	return 0, fmt.Errorf("CPU cannot be accessed at 0x%04x", address+1)
}

func (c *CPU) WriteWord(address uint16, data uint16) error {
	return fmt.Errorf("CPU cannot be accessed at 0x%04x", address+1)
}
