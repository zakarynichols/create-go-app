package colors

import (
	"os"
	"testing"
)

func TestCheckTerminal(t *testing.T) {
	tests := []struct {
		term         string
		reset        string
		red          string
		green        string
		yellow       string
		blue         string
		purple       string
		cyan         string
		white        string
		defaultColor string
	}{
		{
			term:         "xterm",
			reset:        "\033[0m",
			red:          "\033[31m",
			green:        "\033[32m",
			yellow:       "\033[33m",
			blue:         "\033[34m",
			purple:       "\033[35m",
			cyan:         "\033[36m",
			white:        "\033[37m",
			defaultColor: "\033[39m",
		},
		{
			term:         "screen-256color",
			reset:        "\033[38;5;15m",
			red:          "\033[38;5;9m",
			green:        "\033[38;5;10m",
			yellow:       "\033[38;5;11m",
			blue:         "\033[38;5;12m",
			purple:       "\033[38;5;13m",
			cyan:         "\033[38;5;14m",
			white:        "\033[38;5;7m",
			defaultColor: "\033[38;5;15m",
		},
		{
			term:         "unknown-terminal",
			reset:        "",
			red:          "",
			green:        "",
			yellow:       "",
			blue:         "",
			purple:       "",
			cyan:         "",
			white:        "",
			defaultColor: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.term, func(t *testing.T) {
			os.Setenv("TERM", tt.term)
			CheckTerminal(tt.term)

			if Reset != tt.reset {
				t.Errorf("expected Reset: %v, got: %v", tt.reset, Reset)
			}
			if Red != tt.red {
				t.Errorf("expected Red: %v, got: %v", tt.red, Red)
			}
			if Green != tt.green {
				t.Errorf("expected Green: %v, got: %v", tt.green, Green)
			}
			if Yellow != tt.yellow {
				t.Errorf("expected Yellow: %v, got: %v", tt.yellow, Yellow)
			}
			if Blue != tt.blue {
				t.Errorf("expected Blue: %v, got: %v", tt.blue, Blue)
			}
			if Purple != tt.purple {
				t.Errorf("expected Purple: %v, got: %v", tt.purple, Purple)
			}
			if Cyan != tt.cyan {
				t.Errorf("expected Cyan: %v, got: %v", tt.cyan, Cyan)
			}
			if White != tt.white {
				t.Errorf("expected White: %v, got: %v", tt.white, White)
			}
			if Default != tt.defaultColor {
				t.Errorf("expected Default: %v, got: %v", tt.defaultColor, Default)
			}
		})
	}
}
