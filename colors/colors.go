package colors

import (
	"runtime"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

type Colors struct {
	Reset  string
	Red    string
	Green  string
	Yellow string
	Blue   string
	Purple string
	Cyan   string
	Gray   string
	White  string
}

func New() Colors {
	if runtime.GOOS == "window" {
		return Colors{}
	}

	return Colors{
		Reset,
		Red,
		Green,
		Yellow,
		Blue,
		Purple,
		Cyan,
		Gray,
		White,
	}
}
