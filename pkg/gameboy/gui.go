package gameboy

import (
	"context"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"github.com/d2verb/gemu/pkg/gameboy/lcd"
)

type GUI struct {
	app fyne.App
	win fyne.Window
	l   *lcd.LCD
	ctx context.Context
}

func newGUI(ctx context.Context, winTitle string, l *lcd.LCD) GUI {
	a := app.New()
	return GUI{
		app: a,
		win: a.NewWindow(winTitle),
		l:   l,
		ctx: ctx,
	}
}

func (g *GUI) start() {
	g.win.Resize(fyne.NewSize(lcd.ScreenWidth, lcd.ScreenHeight))
	g.win.ShowAndRun()
}

func (g *GUI) updateWindow() {
	for {
		select {
		case <-g.l.Updated:
			g.l.Lock()
			screen := make([][]uint8, len(g.l.Screen))
			for i := range g.l.Screen {
				screen[i] = make([]uint8, len(g.l.Screen[i]))
				copy(screen[i], g.l.Screen[i][:])
			}
			g.l.Unlock()

			g.win.SetContent(canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
				return color.Gray{Y: screen[y][x]}
			}))
		case <-g.ctx.Done():
			goto Done
		default:
			continue
		}
	}
Done:
	g.app.Quit()
}
