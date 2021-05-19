package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Mode int

const (
	VerboseMode Mode = iota + 1
	DebugMode
	WarnMode
	ErrorMode
	FatalMode

	ansiRed    string = "31"
	ansiGreen  string = "32"
	ansiYellow string = "33"
)

var mode Mode = DebugMode

func SetMode(m Mode) {
	mode = m
}

func ModeToString(m Mode) string {
	return map[Mode]string{
		VerboseMode: "verbose",
		DebugMode:   "debug",
		WarnMode:    "warn",
		ErrorMode:   "error",
		FatalMode:   "fatal",
	}[m]
}

func StringToMode(s string) (Mode, error) {
	m, ok := map[string]Mode{
		"verbose": VerboseMode,
		"debug":   DebugMode,
		"warn":    WarnMode,
		"error":   ErrorMode,
		"fatal":   FatalMode,
	}[s]
	if !ok {
		return 0, fmt.Errorf("No corresponding logging mode for %s", s)
	}
	return m, nil
}

func Verbosef(format string, args ...interface{}) {
	log(VerboseMode, ansiGreen, format, args)
}

func Debugf(format string, args ...interface{}) {
	log(DebugMode, ansiGreen, format, args)
}

func Warnf(format string, args ...interface{}) {
	log(WarnMode, ansiYellow, format, args)
}

func Errorf(format string, args ...interface{}) {
	log(ErrorMode, ansiRed, format, args)
}

func Fatalf(format string, args ...interface{}) {
	log(FatalMode, ansiRed, format, args)
	os.Exit(1)
}

func Fprintf(stream io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[%s]%s", time.Now().Format("2006-01-02 15:04:05"), format)
	fmt.Fprintf(stream, format, args...)
}

func log(baseMode Mode, color, format string, args ...interface{}) {
	if mode > baseMode {
		return
	}
	format = fmt.Sprintf("\x1b[%sm[%s]\x1b[0m %s", ModeToString(mode), color, format)
	Fprintf(os.Stderr, format, args...)
}
