package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/zakarynichols/create-go-app/colors"
	"github.com/zakarynichols/create-go-app/print"
)

var (
	ErrNonNameFlag = errors.New(colors.Red + "create-go-app: only one non-named flag argument allowed" + colors.Reset)
	ErrNameFlag    = errors.New(colors.Red + "create-go-app: only a single flag can be used to init a package. e.g. cli or http or module" + colors.Reset)
	ErrDirExists   = errors.New(colors.Red + "create-go-app: directory already exists" + colors.Reset)
	ErrFailMkdir   = errors.New(colors.Red + "create-go-app: failed to create directory" + colors.Reset)
)

func main() {
	start := time.Now()

	var err error

	col := colors.New()

	cli := flag.Bool("cli", false, "set the cli")
	http := flag.Bool("http", false, "set the http")
	module := flag.Bool("module", false, "set the module")

	flag.Parse()

	flags := flag.Args()
	if len(flags) != 1 {
		print.FatalError(ErrNonNameFlag, 1)
	}

	namedFlags := []bool{*cli, *http, *module}
	var flagsLen int // If a flag is provided, increment this variable.
	for i := 0; i < len(namedFlags); i++ {
		if namedFlags[i] {
			flagsLen++
		}
	}

	// Only one type of flag is allowed: cli, http, or module.
	if flagsLen != 1 {
		print.FatalError(ErrNameFlag, 1)
	}

	pkgName := flags[0] // Will make this 'smarter' with help message and such.

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Creating a new %sGo%s app in %s%s/%s%s\n%s", col.Cyan, col.Default, col.Green, dir, pkgName, col.Cyan, col.Default)

	fmt.Print("\n")

	fmt.Printf("Checking if %s./%s%s already exists...%s\n", col.Green, pkgName, col.White, col.Reset)
	_, err = os.Open("./" + pkgName)
	if err == nil {
		print.FatalError(ErrDirExists, 1)
	}

	fmt.Print("\n")

	fmt.Printf("%sMaking new dir %s./%s%s\n", col.White, col.Green, pkgName, col.Reset)
	err = os.Mkdir(pkgName, 0750)
	if err != nil {
		fmt.Printf("\n")
		print.Colorf(col, ErrFailMkdir.Error(), "\n")
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sWriting %smain.go%s file...%s\n", col.White, col.Cyan, col.White, col.Reset)
	err = os.WriteFile(pkgName+"/main.go", []byte(mainTemplate), 0660)
	fmt.Print("\n")
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to write files\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", col.White, col.Cyan, col.Green, pkgName, col.Reset)
	err = os.Chdir("./" + pkgName)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to change directory\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sInitializing a module: %sgo mod init %s%s\n", col.White, col.Cyan, pkgName, col.Reset)
	cmd := exec.Command("go", "mod", "init", pkgName) // don't hardcode module name. only for testing
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to initialize a module\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sFormatting code: %sgo fmt ./...%s\n", col.White, col.Cyan, col.Reset)
	cmd = exec.Command("go", "fmt", "./...")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to format code\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	elapsed := time.Since(start)

	fmt.Print("\n")

	fmt.Printf("%sSucceeded in %f seconds\n%s", col.Green, elapsed.Seconds(), col.Reset)
}

// Put this in a 'code' package along with the other types of templates. cli, http server, module...
var mainTemplate = `
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const Port = 9999

func main() {
	foo := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("path: %s\n", r.URL.Path))
	}

	bar := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, fmt.Sprintf("path: %s\n", r.URL.Path))
	}

	fmt.Printf("Listening on port %d\n", Port)

	http.HandleFunc("/foo", foo)
	http.HandleFunc("/bar", bar)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
`
