package bus

func (b *Bus) SetIE(flag uint8) {
	data := b.Read8(0xffff)
	b.Write8(0xffff, data|flag)
}

func (b *Bus) ClearIE(flag uint8) {
	data := b.Read8(0xffff)
	b.Write8(0xffff, data&(^flag))
}

func (b *Bus) SetIF(flag uint8) {
	data := b.Read8(0xff0f)
	b.Write8(0xff0f, data|flag)
}

func (b *Bus) ClearIF(flag uint8) {
	data := b.Read8(0xff0f)
	b.Write8(0xff0f, data&(^flag))
}
