package gameboy

import "context"

func Start(romPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emu := newGameBoy(ctx, romPath)
	gui := newGUI(ctx, "Gemu", &emu.l)

	go emu.start()

	gui.start()

	return nil
}
