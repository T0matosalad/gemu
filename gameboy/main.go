package gameboy

import (
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/d2verb/gemu/gameboy/rom"
)

func Start(romPath string) error {
	a := app.New()
	w := a.NewWindow("Gemu")
	// c := w.Canvas()

	r := rom.New(romPath)
	log.Printf("Starting game... (%s)\n", r.Title())

	w.Resize(fyne.NewSize(160, 144))
	w.ShowAndRun()

	return nil
}
