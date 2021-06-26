package cpu

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/log"
)

const Hz = 4194304

const (
	IntVBlank   = 1 << 0
	IntLCD      = 1 << 1
	IntTimer    = 1 << 2
	IntSerial   = 1 << 3
	IntJoypad   = 1 << 4
	intSentinel = 5 // Just a sentinel, not actual interrupt flag
)

type CPU struct {
	regs           Registers
	ime            bool  // Interrupt Master Enable Flag
	ie             uint8 // Interrupt Enable
	_if            uint8 // Interrupt Flag
	halt           bool
	bus            *bus.Bus
	instructionSet map[uint16]instruction
}

func New() *CPU {
	return &CPU{
		regs:           newRegisters(),
		halt:           false,
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
	if c.halt {
		return 4, nil
	}

	instAddr := c.regs.PC

	// Fetch opcode
	opcode, err := c.fetch()
	if err != nil {
		return 0, err
	}

	instruction, ok := c.instructionSet[opcode]
	if !ok {
		return 0, fmt.Errorf("Unknown opcode 0x%x (PC: 0x%04x)", opcode, instAddr)
	}

	log.Verbosef("(cpu) [0x%04x]: %s\n", instAddr, instruction.mnemonic)

	// Execute instruction
	cycles, err := instruction.handler(c)
	if err != nil {
		return 0, err
	}

	// Handle interrupts
	if c.ime {
		cyclesForInterrupts, err := c.handleInterrupts()
		if err != nil {
			return cycles, err
		}
		cycles += cyclesForInterrupts
	}

	return cycles, nil
}

func (c *CPU) fetch() (uint16, error) {
	opcode, err := c.bus.ReadUInt8(c.regs.PC)
	if err != nil {
		return 0, err
	}
	c.regs.PC++
	return (uint16)(opcode), nil
}

func (c *CPU) ReadUInt8(address uint16) (uint8, error) {
	switch address {
	case 0xff0f:
		return c._if, nil
	case 0xffff:
		return c.ie, nil
	default:
		return 0, fmt.Errorf("CPU cannot be accessed at 0x%04x", address)
	}
}

func (c *CPU) WriteUInt8(address uint16, data uint8) error {
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

func (c *CPU) ReadUInt16(address uint16) (uint16, error) {
	return 0, fmt.Errorf("CPU cannot be accessed at 0x%04x", address+1)
}

func (c *CPU) WriteUInt16(address uint16, data uint16) error {
	return fmt.Errorf("CPU cannot be accessed at 0x%04x", address+1)
}

func (c *CPU) handleInterrupts() (int, error) {
	filteredFlags := c.ie & c._if
	accumulatedCycles := 0
	for i := 0; i < intSentinel; i++ {
		if filteredFlags&(1<<i) == 0 {
			continue
		}
		cycles, err := c.handleInterrupt(1 << i)
		if err != nil {
			return accumulatedCycles, err
		}
		accumulatedCycles += cycles
	}
	return accumulatedCycles, nil
}

func (c *CPU) handleInterrupt(flag uint8) (int, error) {
	if err := c.bus.ClearIF(flag); err != nil {
		return 0, err
	}

	flagToAddress := map[uint8]uint16{
		IntVBlank: 0x40,
		IntLCD:    0x48,
		IntTimer:  0x50,
		IntSerial: 0x58,
		IntJoypad: 0x60,
	}

	address, _ := flagToAddress[flag]

	return c.callISR(address)
}

func (c *CPU) callISR(address uint16) (int, error) {
	c.ime = false

	// Save current PC to the stack as a return address
	c.regs.SP = c.regs.SP - 2
	if err := c.bus.WriteUInt16(c.regs.SP, c.regs.PC); err != nil {
		return 0, err
	}

	// Jump to the ISR
	c.regs.PC = address

	return 12, nil
}

func (c *CPU) retISR() (int, error) {
	c.ime = true

	// Get return address from stack adn set it to PC
	data, err := c.bus.ReadUInt16(c.regs.SP)
	if err != nil {
		return 0, err
	}
	c.regs.PC = data

	c.regs.SP = c.regs.SP + 2

	return 12, nil
}
