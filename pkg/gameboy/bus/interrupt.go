package bus

func (b *Bus) SetIE(flag uint8) error {
	data, err := b.ReadUInt8((0xffff))
	if err != nil {
		return err
	}
	return b.WriteUInt8(0xffff, data|flag)
}

func (b *Bus) ClearIE(flag uint8) error {
	data, err := b.ReadUInt8((0xffff))
	if err != nil {
		return err
	}
	return b.WriteUInt8(0xffff, data&(^flag))
}

func (b *Bus) SetIF(flag uint8) error {
	data, err := b.ReadUInt8((0xff0f))
	if err != nil {
		return err
	}
	return b.WriteUInt8(0xff0f, data|flag)
}

func (b *Bus) ClearIF(flag uint8) error {
	data, err := b.ReadUInt8((0xff0f))
	if err != nil {
		return err
	}
	return b.WriteUInt8(0xff0f, data&(^flag))
}
