package rom

import "github.com/d2verb/gemu/pkg/gameboy/bus"

type MBC interface {
	AddressRanges() []bus.AddressRange
	Data() []uint8
	ReadByte(uint16) (uint8, error)
	WriteByte(uint16, uint8) error
}
