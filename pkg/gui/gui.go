package gui

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/d2verb/gemu/pkg/gameboy/lcd"
)

type GUI struct {
	app        fyne.App
	win        fyne.Window
	l          *lcd.LCD
	screenHash string
	ratio      int
}

func NewGUI(winTitle string, l *lcd.LCD, ratio int) GUI {
	a := app.New()
	return GUI{
		app:   a,
		win:   a.NewWindow(winTitle),
		l:     l,
		ratio: ratio,
	}
}

func (g *GUI) Start(ctx context.Context, cancel context.CancelFunc) {
	// Start a goroutine to update the screen content
	go func() {
		for {
			select {
			case <-g.l.Updated:
				// Copy the LCD screen buffer
				// Lock() prevents the LCD screen buffer from overwriting
				// by the goroutine of emulator while this copy process
				g.l.Lock()
				screen := make([][]uint8, len(g.l.Screen))
				for i := range g.l.Screen {
					screen[i] = make([]uint8, len(g.l.Screen[i]))
					copy(screen[i], g.l.Screen[i][:])
				}
				g.l.Unlock()

				// If screen content is the same, skip gui updating
				screenHash := calcScreenHash(screen)
				if screenHash == g.screenHash {
					continue
				}
				g.screenHash = screenHash

				g.win.SetContent(canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
					actualX := x * lcd.ScreenWidth / w
					actualY := y * lcd.ScreenHeight / h
					dot := screen[actualY][actualX]
					return color.RGBA{dot, dot, dot, 0xff}
				}))
			case <-ctx.Done():
				g.app.Quit()
				return
			default:
				continue
			}
		}
	}()

	g.win.Resize(fyne.NewSize(float32(lcd.ScreenWidth*g.ratio), float32(lcd.ScreenHeight*g.ratio)))
	g.win.SetFixedSize(true)
	g.win.ShowAndRun()
}

func calcScreenHash(screen [][]uint8) string {
	h := sha256.New()
	for i := range screen {
		h.Write(screen[i])
	}
	return hex.EncodeToString(h.Sum(nil))
}
