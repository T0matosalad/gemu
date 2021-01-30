package bus

import (
	"fmt"
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

func (b *Bus) ReadByte(address uint16) (uint8, error) {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		return 0, err
	}
	return device.ReadByte(address)
}

func (b *Bus) ReadWord(address uint16) (uint16, error) {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		return 0, err
	}
	return device.ReadWord(address)
}

func (b *Bus) WriteByte(address uint16, data uint8) error {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		return err
	}
	return device.WriteByte(address, data)
}

func (b *Bus) WriteWord(address uint16, data uint16) error {
	device, err := b.findDeviceFromAddress(address)
	if err != nil {
		return err
	}
	return device.WriteWord(address, data)
}

func (b *Bus) findDeviceFromAddress(address uint16) (Addressable, error) {
	for _range, device := range b.addressSpace {
		if _range.Contains(address) {
			return device, nil
		}
	}
	return nil, fmt.Errorf("No addressable device at 0x%04x", address)
}
