package cpu

import "testing"

func TestSignExtU8ToU16(t *testing.T) {
	tests := []struct {
		in   uint8
		want uint16
	}{
		{0x7f, 0x007f},
		{0x80, 0xff80},
		{0xff, 0xffff},
	}

	for _, tt := range tests {
		got := signExtU8ToU16(tt.in)
		if got != tt.want {
			t.Errorf("signExtU8ToU16(0x%02x) = 0x%04x, want 0x%04x", tt.in, got, tt.want)
		}
	}
}

func TestAdd8(t *testing.T) {
	c := New()

	res := c.add8(0x0f, 0x01, true)
	if res != 0x10 {
		t.Errorf("expected 0x10, got 0x%02x", res)
	}
	if c.regs.Flag(NFlag) != 0 {
		t.Errorf("N flag should be cleared")
	}
	if c.regs.Flag(HFlag) == 0 {
		t.Errorf("H flag should be set")
	}
	if c.regs.Flag(CFlag) != 0 {
		t.Errorf("C flag should not be set")
	}
	if c.regs.Flag(ZFlag) != 0 {
		t.Errorf("Z flag should not be set")
	}

	c = New()
	res = c.add8(0xff, 0x01, true)
	if res != 0x00 {
		t.Errorf("expected 0x00, got 0x%02x", res)
	}
	if c.regs.Flag(ZFlag) == 0 {
		t.Errorf("Z flag should be set")
	}
	if c.regs.Flag(CFlag) == 0 {
		t.Errorf("C flag should be set")
	}
	if c.regs.Flag(HFlag) == 0 {
		t.Errorf("H flag should be set")
	}
}

func TestSub8(t *testing.T) {
	c := New()

	res := c.sub8(0x02, 0x01, true)
	if res != 0x01 {
		t.Errorf("expected 0x01, got 0x%02x", res)
	}
	if c.regs.Flag(NFlag) == 0 {
		t.Errorf("N flag should be set")
	}
	if c.regs.Flag(HFlag) != 0 {
		t.Errorf("H flag should not be set")
	}
	if c.regs.Flag(ZFlag) != 0 {
		t.Errorf("Z flag should not be set")
	}

	c = New()
	res = c.sub8(0x01, 0x02, true)
	if res != 0xff {
		t.Errorf("expected 0xff, got 0x%02x", res)
	}
	if c.regs.Flag(CFlag) == 0 {
		t.Errorf("C flag should be set")
	}
	if c.regs.Flag(HFlag) == 0 {
		t.Errorf("H flag should be set")
	}
}

func TestBit8(t *testing.T) {
	c := New()

	c.bit8(0x10, 4)
	if c.regs.Flag(ZFlag) != 0 {
		t.Errorf("Z flag should be cleared")
	}
	if c.regs.Flag(HFlag) == 0 {
		t.Errorf("H flag should be set")
	}
	if c.regs.Flag(NFlag) != 0 {
		t.Errorf("N flag should be cleared")
	}

	c.bit8(0x00, 3)
	if c.regs.Flag(ZFlag) == 0 {
		t.Errorf("Z flag should be set")
	}
}
