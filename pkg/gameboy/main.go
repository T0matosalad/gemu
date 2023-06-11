package gameboy

import (
	"context"
	"io/ioutil"

	"github.com/d2verb/gemu/pkg/debug"
)

func Start(romPath string, ratio int, debugMode bool) error {
	romContent, err := ioutil.ReadFile(romPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan any)

	emu, err := newGameBoy(romContent, ch, debugMode)
	if err != nil {
		return err
	}

	gui := newGUI("Gemu", emu.l, ratio)
	dbg := debug.NewDebugServer(9000, ch, debugMode)

	go emu.start(ctx, cancel)
	go dbg.Start(ctx, cancel)
	gui.start(ctx, cancel)

	return nil
}
