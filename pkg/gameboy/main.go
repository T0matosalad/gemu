package gameboy

import (
	"context"
	"io/ioutil"
)

func Start(romPath string) error {
	romContent, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emu := newGameBoy(romContent)
	gui := newGUI("Gemu", &emu.l)

	go emu.start(ctx, cancel)
	gui.start(ctx, cancel)

	return nil
}
