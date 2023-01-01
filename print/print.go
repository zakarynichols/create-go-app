package print

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/zakarynichols/create-go-app/colors"
	"github.com/zakarynichols/create-go-app/perm"
)

const logFile = "create-go-app-debug.txt"

// Colorf prints a color formatted string. After writing its
// output, the terminal color is reset to the operating system default.
func Colorf(strs ...string) {
	col := colors.New()
	for i := 0; i < len(strs); i++ {
		fmt.Printf("%s%s", strs[i], col.Default)
	}
}

// Errorln prints a formatted error with new
// lines above and below the written output.
// After writing the output it will exit
// with a provided status code.
func Errorln(err error) {
	fmt.Print("\n")
	fmt.Printf("%s%s\n%s", colors.Red, err.Error(), colors.Default)
	fmt.Print("\n")
}

// WriteDebugLog writes the error to a file. Might not need it and if it
// stays needs to collect more information useful for troubleshooting.
// Probably should have this be a non-default opt-in flag.
func WriteDebugLog(debugErr error) {
	err := os.Remove(logFile)
	// Ignore error if a file or dir doesn't exist. Report an error otherwise.
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, perm.RWX)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.WriteString(debugErr.Error() + "\n\n" + string(debug.Stack()))
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
