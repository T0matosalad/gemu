package gameboy

import "context"

func Start(romPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emu := newGameBoy(romPath)
	gui := newGUI("Gemu", &emu.l)

	go emu.start(ctx, cancel)
	gui.start(ctx, cancel)

	return nil
}
