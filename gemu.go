package gemu

import (
	"flag"
	"fmt"
	"os"

	"github.com/d2verb/gemu/gameboy"
)

const version = "0.0.1"

func Run() error {
	flag.Usage = flagUsage

	v := flag.Bool("v", false, "display version")
	flag.Parse()

	if *v {
		fmt.Printf("gemu v%s\n", version)
		return nil
	}

	if len(flag.Args()) != 1 {
		return flag.ErrHelp
	}

	return gameboy.Start(flag.Arg(0))
}

func flagUsage() {
	usageText := `Usage of gemu:

gemu [-v] ROM
    -v    display version`

	fmt.Fprintf(os.Stderr, "%s\n", usageText)
}
