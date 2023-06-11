package lcd

import "sync"

const (
	ScreenWidth  = 160
	ScreenHeight = 144
)

type LCD struct {
	sync.Mutex
	Updated chan any
	Screen  [ScreenHeight][ScreenWidth]uint8
}

func New() *LCD {
	return &LCD{
		Updated: make(chan any),
	}
}
