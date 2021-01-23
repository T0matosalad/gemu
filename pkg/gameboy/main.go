package gameboy

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/gameboy/cpu"
	"github.com/d2verb/gemu/pkg/gameboy/ppu"
	"github.com/d2verb/gemu/pkg/gameboy/rom"
	"github.com/d2verb/gemu/pkg/log"
)

func Start(romPath string) error {
	a := app.New()
	w := a.NewWindow("Gemu")
	// c := w.Canvas()

	cpu := cpu.New()
	rom := rom.New(romPath)
	ppu := ppu.New()
	bus := bus.New()

	cpu.ConnectToBus(&bus)
	rom.ConnectToBus(&bus)
	ppu.ConnectToBus(&bus)

	log.Debugf("Starting game... (%s)\n", rom.Title())

	for {
		cycles, err := cpu.Step()
		if err != nil {
			return err
		}

		err = ppu.Step(cycles)
		if err != nil {
			return err
		}
	}

	w.Resize(fyne.NewSize(160, 144))
	w.ShowAndRun()

	return nil
}
