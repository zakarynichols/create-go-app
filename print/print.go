package print

import (
	"fmt"
	"os"

	"github.com/zakarynichols/create-go-app/colors"
)

// Colorf prints a color formatted string. After writing its
// output, the terminal color is reset to the operating system default.
func Colorf(col colors.Colors, strs ...string) {
	for i := 0; i < len(strs); i++ {
		fmt.Printf("%s%s", strs[i], col.Default)
	}
}

// FatalError prints a formatted error with new
// lines above and below the written output.
// After writing the output it will exit
// with a provided status code.
func FatalError(err error, code int) {
	fmt.Print("\n")
	fmt.Printf("%s\n", err.Error())
	fmt.Print("\n")
	os.Exit(code)
}
