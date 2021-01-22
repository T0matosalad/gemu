package cpu

import (
	"github.com/d2verb/gemu/gameboy/bus"
)

type CPU struct {
	Regs Registers
	bus  *bus.Bus
}

func New() CPU {
	return CPU{
		Regs: newRegisters(),
	}
}

func (c *CPU) ConnectToBus(bus *bus.Bus) error {
	c.bus = bus
	return nil
}
