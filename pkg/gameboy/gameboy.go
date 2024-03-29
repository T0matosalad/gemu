package gameboy

import (
	"context"

	"github.com/d2verb/gemu/pkg/debug/pb"
	"github.com/d2verb/gemu/pkg/gameboy/apu"
	"github.com/d2verb/gemu/pkg/gameboy/bus"
	"github.com/d2verb/gemu/pkg/gameboy/cpu"
	"github.com/d2verb/gemu/pkg/gameboy/lcd"
	"github.com/d2verb/gemu/pkg/gameboy/ppu"
	"github.com/d2verb/gemu/pkg/gameboy/ram"
	"github.com/d2verb/gemu/pkg/gameboy/rom"
	"github.com/d2verb/gemu/pkg/log"
)

type GameBoy struct {
	c         *cpu.CPU
	r         *rom.ROM
	a         *ram.RAM
	l         *lcd.LCD
	p         *ppu.PPU
	s         *apu.APU
	b         *bus.Bus
	ch        chan any
	debugMode bool
}

func NewGameBoy(romContent []uint8, ch chan any, debugMode bool) (*GameBoy, error) {
	l := lcd.New()

	r, err := rom.New(romContent)
	if err != nil {
		return nil, err
	}

	g := GameBoy{
		c:         cpu.New(),
		r:         r,
		a:         ram.New(),
		l:         l,
		p:         ppu.New(l),
		s:         apu.New(),
		b:         bus.New(),
		ch:        ch,
		debugMode: debugMode,
	}
	g.c.ConnectToBus(g.b)
	g.r.ConnectToBus(g.b)
	g.a.ConnectToBus(g.b)
	g.p.ConnectToBus(g.b)
	g.s.ConnectToBus(g.b)

	return &g, nil
}

func (g *GameBoy) LCD() *lcd.LCD {
	return g.l
}

func (g *GameBoy) Start(ctx context.Context, cancel context.CancelFunc) {
	log.Debugf("Starting game... (%s)\n", g.r.String())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if g.debugMode {
				runNextEmulatorStep := g.debuggerStep()
				if !runNextEmulatorStep {
					continue
				}
			}

			cycles := g.c.Step()
			g.p.Step(cycles)
		}
	}
}

func (g *GameBoy) debuggerStep() (runNextEmulatorStep bool) {
	req := <-g.ch

	switch req.(type) {
	case *pb.NextRequest:
		g.ch <- pb.NextReply{}
		runNextEmulatorStep = true
	default:
		log.Errorf("Unknown debug request: %T\n", req)
		g.ch <- nil
		runNextEmulatorStep = false
	}

	return
}
