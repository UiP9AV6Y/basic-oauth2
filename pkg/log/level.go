package log

import (
	"fmt"
	"strings"
)

type Level int

var levelNames = []string{
	"Off",
	"Fatal",
	"Error",
	"Warn",
	"Info",
	"Debug",
	"Trace",
}

const (
	LevelOff Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func ParseLevelVerbosity(verbosity int) (level Level, err error) {
	for v := range levelNames {
		if v == verbosity {
			level = Level(v)
			return
		}
	}

	err = fmt.Errorf("invalid log level verbosity %d", verbosity)
	return
}

func ParseLevelName(name string) (level Level, err error) {
	if name == "" {
		level = LevelOff
		return
	}

	for l, n := range levelNames {
		if strings.EqualFold(n, name) {
			level = Level(l)
			return
		}
	}

	err = fmt.Errorf("invalid log level name %q", name)
	return
}

func (l Level) String() string {
	return levelNames[l]
}
