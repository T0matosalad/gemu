package bus

import (
	"fmt"

	"github.com/d2verb/gemu/pkg/log"
)

type Bus struct {
	addressSpace map[AddressRange]Addressable
}

func New() *Bus {
	return &Bus{
		addressSpace: map[AddressRange]Addressable{},
	}
}

func (b *Bus) Map(targetRange AddressRange, device Addressable) error {
	for _range := range b.addressSpace {
		if _range.IsOverlapped(targetRange) {
			return fmt.Errorf("Address range %s and %s is overlapped", _range.String(), targetRange.String())
		}
	}
	b.addressSpace[targetRange] = device
	return nil
}

func (b *Bus) Read8(address uint16) uint8 {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return device.Read8(address)
}

func (b *Bus) Read16(address uint16) uint16 {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	return device.Read16(address)
}

func (b *Bus) Write8(address uint16, data uint8) {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	device.Write8(address, data)
}

func (b *Bus) Write16(address uint16, data uint16) {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	device.Write16(address, data)
}

func (b *Bus) findDeviceFromAddress(address uint16) (Addressable, error) {
	for _range, device := range b.addressSpace {
		if _range.Contains(address) {
			return device, nil
		}
	}
	return nil, fmt.Errorf("No addressable device at 0x%04x", address)
}
