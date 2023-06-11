package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/d2verb/gemu"
)

func main() {
	config, err := gemu.SetUp()
	if err != nil {
		if err == flag.ErrHelp {
			flag.Usage()
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
	gemu.Run(config)
}
