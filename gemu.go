package gemu

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/d2verb/gemu/pkg/debug"
	"github.com/d2verb/gemu/pkg/gameboy"
	"github.com/d2verb/gemu/pkg/gui"
	"github.com/d2verb/gemu/pkg/log"
)

const version = "0.0.1"

type Config struct {
	RomPath   string
	Ratio     int
	DebugMode bool
}

func SetUp() (*Config, error) {
	flag.Usage = flagUsage

	v := flag.Bool("v", false, "display version")
	r := flag.Int("r", 1, "magnification ratio of screen")
	l := flag.String("l", log.ModeToString(log.DebugMode), "log level")
	d := flag.Bool("d", false, "start debug server")
	flag.Parse()

	if *v {
		fmt.Printf("gemu v%s\n", version)
		return nil, nil
	}

	if len(flag.Args()) != 1 {
		return nil, flag.ErrHelp
	}

	mode, err := log.StringToMode(*l)
	if err != nil {
		return nil, err
	}
	log.SetMode(mode)

	return &Config{
		RomPath:   flag.Arg(0),
		Ratio:     *r,
		DebugMode: *d,
	}, nil
}

func Run(config *Config) error {
	romContent, err := ioutil.ReadFile(config.RomPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan any)

	gb, err := gameboy.NewGameBoy(romContent, ch, config.DebugMode)
	if err != nil {
		return err
	}

	gui := gui.NewGUI("Gemu", gb.LCD(), config.Ratio)
	dbg := debug.NewDebugServer(9000, ch, config.DebugMode)

	go gb.Start(ctx, cancel)
	go dbg.Start(ctx, cancel)
	gui.Start(ctx, cancel)

	return nil
}

func flagUsage() {
	usageText := `Usage of gemu:

gemu [-vrd] ROM
    -v         display version
    -r int     magnification ratio of screen (default: 1)
    -l string  log level {verbose, debug, warn, error, fatal} (default: debug)
    -d         start debug mode`

	fmt.Fprintf(os.Stderr, "%s\n", usageText)
}
