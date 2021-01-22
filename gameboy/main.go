package gameboy

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/d2verb/gemu/gameboy/bus"
	"github.com/d2verb/gemu/gameboy/cpu"
	"github.com/d2verb/gemu/gameboy/rom"
)

func Start(romPath string) error {
	a := app.New()
	w := a.NewWindow("Gemu")
	// c := w.Canvas()

	cpu := cpu.New()
	rom := rom.New(romPath)
	bus := bus.New()

	cpu.ConnectToBus(&bus)
	rom.ConnectToBus(&bus)

	log.Printf("Starting game... (%s)\n", rom.Title())

	w.Resize(fyne.NewSize(160, 144))
	w.ShowAndRun()

	return nil
}
