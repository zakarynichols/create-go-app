package colors

import (
	"os"
	"strings"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Default = "\033[39m"
)

// CheckTerminal will check the TERM environment variable to determine the type of terminal being used
// and adjust the ANSI codes accordingly
func CheckTerminal() {
	term := os.Getenv("TERM")
	if strings.Contains(term, "xterm") || strings.Contains(term, "rxvt") || strings.Contains(term, "linux") {
		// ANSI codes will work as-is on these terminal types
		return
	} else if strings.Contains(term, "screen") {
		// Screen-256color terminal type requires different ANSI codes
		Reset = "\033[38;5;15m"
		Red = "\033[38;5;9m"
		Green = "\033[38;5;10m"
		Yellow = "\033[38;5;11m"
		Blue = "\033[38;5;12m"
		Purple = "\033[38;5;13m"
		Cyan = "\033[38;5;14m"
		White = "\033[38;5;7m"
		Default = "\033[38;5;15m"
	} else {
		// If the terminal type is not recognized, set all colors to the default
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		White = ""
		Default = ""
	}
}
