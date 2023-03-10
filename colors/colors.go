package colors

import (
	"fmt"
	"os"
	"strings"
)

type Color string

var (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Purple  Color = "\033[35m"
	Cyan    Color = "\033[36m"
	White   Color = "\033[37m"
	Default Color = "\033[39m"
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

// Printf formats a colored string according to a format specifier. It's
// an explicit definition for creating color formatted output strings.
func Printf(format string, a ...any) error {
	_, err := fmt.Printf(format, a...)
	if err != nil {
		return err
	}
	return nil
}
