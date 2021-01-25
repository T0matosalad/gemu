package gameboy

import (
	"context"
	"time"

	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/gameboy/cpu"
	"github.com/d2verb/gemu/pkg/gameboy/lcd"
	"github.com/d2verb/gemu/pkg/gameboy/ppu"
	"github.com/d2verb/gemu/pkg/gameboy/rom"
	"github.com/d2verb/gemu/pkg/log"
)

type GameBoy struct {
	c   cpu.CPU
	r   rom.ROM
	l   lcd.LCD
	p   ppu.PPU
	b   bus.Bus
	ctx context.Context
}

func newGameBoy(ctx context.Context, romPath string) GameBoy {
	l := lcd.New()
	g := GameBoy{
		c:   cpu.New(),
		r:   rom.New(romPath),
		l:   l,
		p:   ppu.New(&l),
		b:   bus.New(),
		ctx: ctx,
	}
	g.c.ConnectToBus(&g.b)
	g.r.ConnectToBus(&g.b)
	g.p.ConnectToBus(&g.b)
	return g
}

func (g *GameBoy) start() {
	log.Debugf("Starting game... (%s)\n", g.r.Title())

	startTime := Now()
	accumulatedCycles := 0

	for {
		select {
		case <-g.ctx.Done():
			return
		default:
			cycles, err := g.c.Step()
			if err != nil {
				log.Errorf("%s\n", err.Error())
				return
			}

			err = g.p.Step(cycles)
			if err != nil {
				log.Errorf("%s\n", err.Error())
				return
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
	}
}

func Now() int64 {
	return time.Now().Unix()*1000 + time.Now().UnixNano()/int64(time.Millisecond)
}
