package rom

import (
	"io/ioutil"
	"log"
)

type ROM struct {
	data []uint8
}

func New(romPath string) ROM {
	content, err := ioutil.ReadFile(romPath)
	if err != nil {
		log.Fatal(err)
	}
	return ROM{
		data: content,
	}
}

func (r *ROM) Title() []uint8 {
	start, end := 0x134, 0x134
	for ; end < 0x144 && r.data[end] != 0; end++ {
	}
	return r.data[start:end]
}

func (r *ROM) MBCType() uint8 {
	return r.data[0x147]
}
