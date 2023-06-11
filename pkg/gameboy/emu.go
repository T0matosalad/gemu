package gameboy

import (
	"context"
	"time"

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

func newGameBoy(romContent []uint8, ch chan any, debugMode bool) (*GameBoy, error) {
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

func (g *GameBoy) start(ctx context.Context, cancel context.CancelFunc) {
	log.Debugf("Starting game... (%s)\n", g.r.String())

	startTime := NowInMillisecond()
	accumulatedCycles := 0

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

			// Ensure that the CPU only runs cpu.Hz cycles per second
			accumulatedCycles += cycles
			if accumulatedCycles >= cpu.Hz {
				elapsedTime := NowInMillisecond() - startTime
				if elapsedTime < 1000 {
					duration := time.Duration(1000 - elapsedTime)
					time.Sleep(duration * time.Millisecond)
				}
				accumulatedCycles -= cpu.Hz
				startTime = NowInMillisecond()
			}
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
		g.ch <- struct{}{}
		runNextEmulatorStep = false
	}

	return
}

func NowInMillisecond() int64 {
	return time.Now().Unix()*1000 + time.Now().UnixNano()/int64(time.Millisecond)
}
