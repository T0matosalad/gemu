package cpu

import (
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

func (c *CPU) Step() int {
	if c.halt {
		return 4
	}

	instAddr := c.regs.PC

	// Fetch opcode
	opcode := c.fetch()
	instruction, ok := c.instructionSet[opcode]
	if !ok {
		log.Fatalf("Unknown opcode 0x%x (PC: 0x%04x)", opcode, instAddr)
	}

	log.Verbosef("(cpu) [0x%04x]: %s\n", instAddr, instruction.mnemonic)

	// Execute instruction
	cycles := instruction.handler(c)

	// Handle interrupts
	if c.ime {
		cycles += c.handleInterrupts()
	}

	return cycles
}

func (c *CPU) fetch() uint16 {
	opcode := c.bus.Read8(c.regs.PC)
	c.regs.PC++
	return (uint16)(opcode)
}

func (c *CPU) Read8(address uint16) uint8 {
	switch address {
	case 0xff0f:
		return c._if
	case 0xffff:
		return c.ie
	default:
		log.Fatalf("CPU cannot be accessed at 0x%04x", address)
	}
	return 0
}

func (c *CPU) Write8(address uint16, data uint8) {
	switch address {
	case 0xff0f:
		c._if = data
	case 0xffff:
		c.ie = data
	default:
		log.Fatalf("CPU cannot be accessed at 0x%04x", address)
	}
}

func (c *CPU) Read16(address uint16) uint16 {
	log.Fatalf("CPU cannot be accessed at 0x%04x", address+1)
	return 0
}

func (c *CPU) Write16(address uint16, data uint16) {
	log.Fatalf("CPU cannot be accessed at 0x%04x", address+1)
}

func (c *CPU) handleInterrupts() int {
	filteredFlags := c.ie & c._if
	accumulatedCycles := 0
	for i := 0; i < intSentinel; i++ {
		if filteredFlags&(1<<i) == 0 {
			continue
		}
		accumulatedCycles += c.handleInterrupt(1 << i)
	}
	return accumulatedCycles
}

func (c *CPU) handleInterrupt(flag uint8) int {
	c.bus.ClearIF(flag)

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

func (c *CPU) callISR(address uint16) int {
	c.ime = false

	// Save current PC to the stack as a return address
	c.regs.SP = c.regs.SP - 2
	c.bus.Write16(c.regs.SP, c.regs.PC)

	// Jump to the ISR
	c.regs.PC = address

	return 12
}

func (c *CPU) retISR() int {
	c.ime = true

	// Get return address from stack adn set it to PC
	c.regs.PC = c.bus.Read16(c.regs.SP)
	c.regs.SP = c.regs.SP + 2

	return 12
}
