package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/d2verb/gemu"
)

func main() {
	err := gemu.Run()
	if err != nil {
		if err == flag.ErrHelp {
			flag.Usage()
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
