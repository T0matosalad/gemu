package bus

import "fmt"

type Addressable interface {
	ReadUInt8(uint16) (uint8, error)
	ReadUInt16(uint16) (uint16, error)
	WriteUInt8(uint16, uint8) error
	WriteUInt16(uint16, uint16) error
	ConnectToBus(bus *Bus) error
}

type AddressRange struct {
	Start uint16
	End   uint16
}

func NewAddressRange(start uint16, end uint16) AddressRange {
	return AddressRange{
		Start: start,
		End:   end,
	}
}

func (a *AddressRange) IsOverlapped(b AddressRange) bool {
	return !(b.End < a.Start || a.End < b.Start)
}

func (a *AddressRange) Contains(address uint16) bool {
	return a.Start <= address && address <= a.End
}

func (a *AddressRange) String() string {
	return fmt.Sprintf("[%d, %d]", a.Start, a.End)
}
