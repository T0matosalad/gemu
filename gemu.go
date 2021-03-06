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
	d := flag.String("d", log.ModeToString(log.DebugMode), "debug level")
	flag.Parse()

	if *v {
		fmt.Printf("gemu v%s\n", version)
		return nil
	}

	if len(flag.Args()) != 1 {
		return flag.ErrHelp
	}

	mode, err := log.StringToMode(*d)
	if err != nil {
		return err
	}
	log.SetMode(mode)

	return gameboy.Start(flag.Arg(0), *r)
}

func flagUsage() {
	usageText := `Usage of gemu:

gemu [-vrd] ROM
    -v         display version
    -r int     magnification ratio of screen (default: 1)
    -d string  debug level {debug, warn, error, fatal} (default: debug)`

	fmt.Fprintf(os.Stderr, "%s\n", usageText)
}
