package rom

import "github.com/d2verb/gemu/pkg/gameboy/bus"

type MBC interface {
	AddressRanges() []bus.AddressRange
	Data() []uint8
	ReadUInt8(uint16) (uint8, error)
	WriteUInt8(uint16, uint8) error
}
