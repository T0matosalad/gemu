package gameboy

import (
	"context"
	"io/ioutil"
)

func Start(romPath string, ratio int) error {
	romContent, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emu, err := newGameBoy(romContent)
	if err != nil {
		return err
	}

	gui := newGUI("Gemu", emu.l, ratio)

	go emu.start(ctx, cancel)
	gui.start(ctx, cancel)

	return nil
}
