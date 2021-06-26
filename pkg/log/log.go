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

func Fprintf(stream io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), format)
	fmt.Fprintf(stream, format, args...)
}

func Verbosef(format string, args ...interface{}) {
	if mode > VerboseMode {
		return
	}
	format = fmt.Sprintf("\x1b[32m[verbose]\x1b[0m %s", format)
	Fprintf(os.Stdout, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if mode > DebugMode {
		return
	}
	format = fmt.Sprintf("\x1b[32m[debug]\x1b[0m %s", format)
	Fprintf(os.Stdout, format, args...)
}

func Warnf(format string, args ...interface{}) {
	if mode > WarnMode {
		return
	}
	format = fmt.Sprintf("\x1b[33m[warn]\x1b[0m %s", format)
	Fprintf(os.Stdout, format, args...)
}

func Errorf(format string, args ...interface{}) {
	if mode > ErrorMode {
		return
	}
	format = fmt.Sprintf("\x1b[31m[error]\x1b[0m %s", format)
	Fprintf(os.Stderr, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	if mode > FatalMode {
		return
	}
	format = fmt.Sprintf("\x1b[31m[fatal]\x1b[0m %s", format)
	Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
