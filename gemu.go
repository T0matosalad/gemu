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
	flag.Parse()

	if *v {
		fmt.Printf("gemu v%s\n", version)
		return nil
	}

	if len(flag.Args()) != 1 {
		return flag.ErrHelp
	}

	log.SetMode(log.DebugMode)

	return gameboy.Start(flag.Arg(0), *r)
}

func flagUsage() {
	usageText := `Usage of gemu:

gemu [-v] ROM
    -v    display version`

	fmt.Fprintf(os.Stderr, "%s\n", usageText)
}
