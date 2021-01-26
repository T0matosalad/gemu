package bus

import "testing"

func TestIsOverlapped(t *testing.T) {
	a := NewAddressRange(0, 10)
	b := NewAddressRange(5, 15)
	if !a.IsOverlapped(b) {
		t.Fatalf("%s and %s is overlapped, but IsOverlapped() returns false", a.String(), b.String())
	}

	a = NewAddressRange(5, 15)
	b = NewAddressRange(0, 10)
	if !a.IsOverlapped(b) {
		t.Fatalf("%s and %s is overlapped, but IsOverlapped() returns false", a.String(), b.String())
	}

	a = NewAddressRange(0, 10)
	b = NewAddressRange(15, 30)
	if a.IsOverlapped(b) {
		t.Fatalf("%s and %s is not overlapped, but IsOverlapped() returns true", a.String(), b.String())
	}

	a = NewAddressRange(15, 30)
	b = NewAddressRange(0, 10)
	if a.IsOverlapped(b) {
		t.Fatalf("%s and %s is not overlapped, but IsOverlapped() returns true", a.String(), b.String())
	}
}
