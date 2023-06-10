package gemu

import (
	"flag"
	"fmt"
	"os"

	"github.com/d2verb/gemu/pkg/gameboy"
	"github.com/d2verb/gemu/pkg/log"
)

const version = "0.0.1"

func Run() error {
	flag.Usage = flagUsage

	v := flag.Bool("v", false, "display version")
	r := flag.Int("r", 1, "magnification ratio of screen")
	l := flag.String("l", log.ModeToString(log.DebugMode), "log level")
	d := flag.Bool("d", false, "start debug server")
	flag.Parse()

	if *v {
		fmt.Printf("gemu v%s\n", version)
		return nil
	}

	if len(flag.Args()) != 1 {
		return flag.ErrHelp
	}

	mode, err := log.StringToMode(*l)
	if err != nil {
		return err
	}
	log.SetMode(mode)

	return gameboy.Start(flag.Arg(0), *r, *d)
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
