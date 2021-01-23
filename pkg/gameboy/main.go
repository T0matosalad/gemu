package gameboy

import (
	"time"

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

	c := cpu.New()
	r := rom.New(romPath)
	p := ppu.New()
	b := bus.New()

	c.ConnectToBus(&b)
	r.ConnectToBus(&b)
	p.ConnectToBus(&b)

	log.Debugf("Starting game... (%s)\n", r.Title())

	startTime := Now()
	accumulatedCycles := 0

	for {
		cycles, err := c.Step()
		if err != nil {
			return err
		}

		err = p.Step(cycles)
		if err != nil {
			return err
		}

		// Ensure that the CPU only runs cpu.Hz cycles per second
		accumulatedCycles += cycles
		if accumulatedCycles >= cpu.Hz {
			elapsedTime := Now() - startTime
			if elapsedTime < 1000 {
				duration := time.Duration(1000 - elapsedTime)
				time.Sleep(duration * time.Millisecond)
			}
			accumulatedCycles -= cpu.Hz
			startTime = Now()
		}
	}

	w.Resize(fyne.NewSize(160, 144))
	w.ShowAndRun()

	return nil
}

func Now() int64 {
	return time.Now().Unix()*1000 + time.Now().UnixNano()/int64(time.Millisecond)
}
