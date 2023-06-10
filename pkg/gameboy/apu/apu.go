package apu

import "github.com/d2verb/gemu/pkg/gameboy/bus"

type APU struct {
	bus *bus.Bus
}

func New() *APU {
	return &APU{}
}

func (a *APU) ConnectToBus(b *bus.Bus) error {
	if err := b.Map(bus.NewAddressRange(0xff10, 0xff26), a); err != nil {
		return err
	}
	a.bus = b
	return nil
}

func (a *APU) Read8(address uint16) uint8 {
	return 0
}

func (a *APU) Write8(address uint16, data uint8) {
}

func (a *APU) Read16(address uint16) uint16 {
	return 0
}

func (a *APU) Write16(address uint16, data uint16) {
}
