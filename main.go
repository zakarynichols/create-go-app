package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/zakarynichols/create-go-app/clock"
	"github.com/zakarynichols/create-go-app/code"
	"github.com/zakarynichols/create-go-app/colors"
	"github.com/zakarynichols/create-go-app/perm"
	"github.com/zakarynichols/create-go-app/print"
)

// Errors exposed to the user. Stack traces and more detailed
// errors for debugging will be written to a log file.
var (
	ErrNonNameFlag = errors.New("create-go-app: only one non-named flag argument allowed")
	ErrNameFlag    = errors.New("create-go-app: only a single named flag can be used to init a package. e.g. --cli, --http, or --module")
	ErrDirExists   = errors.New("create-go-app: directory already exists")
	ErrMkdir       = errors.New("create-go-app: failed to create directory")
	ErrChdir       = errors.New("create-go-app: failed to change directory")
	ErrWkdir       = errors.New("create-go-app: failed to get working directory")
	ErrInitMod     = errors.New("create-go-app: failed to init a module")
	ErrFmt         = errors.New("create-go-app: failed to format code")
	ErrWriteFiles  = errors.New("create-go-app: failed to write files")
)

// application is the structure describing the initialized application.
// It should have a directory name and a module name.
type application struct {
	dirname string
	module  string
}

func main() {
	app := new(application)

	timer := clock.Timer()

	namedFlagPtrs := setupFlags()

	flag.Usage = usage

	flag.Parse()

	nonNamedFlags := flag.Args()
	checkNonNamed(nonNamedFlags)
	app.dirname = nonNamedFlags[0]

	checkNamed(namedFlagPtrs)

	printWkdir(app.dirname)

	fmt.Print("\n")

	checkExists(app.dirname)

	fmt.Print("\n")

	mkdir(app.dirname)

	fmt.Print("\n")

	writeFile(app.dirname)

	chdir(app.dirname)

	fmt.Print("\n")

	app.modInit()

	fmt.Print("\n")

	fmtCode()

	elapsed := timer.Since(timer.Start)

	fmt.Print("\n")

	fmt.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)
}

func fmtCode() {
	fmt.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)
	cmd := exec.Command("go", "fmt", "./...")
	err := cmd.Run()
	if err != nil {
		fatal(ErrFmt)
	}
}

// modInit sets up the package with a custom module name read from stdin.
func (app *application) modInit() {
	fmt.Print("go mod init: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("failed to read string with err: ", err)
	}

	fmt.Print("\n")

	app.module = strings.Trim(input, "\r\n")

	fmt.Printf("%sInitializing a module: %sgo mod init %s%s\n", colors.White, colors.Cyan, app.module, colors.Default)
	cmd := exec.Command("go", "mod", "init", app.module)
	err = cmd.Run()
	if err != nil {
		fatal(ErrInitMod)
	}
}

// fatal is the main error handling mechanism. It prints an error message, writes a debug log, and exits with a status code 1.
func fatal(err error) {
	print.Errorln(err)
	print.WriteDebugLog(err)
	os.Exit(1)
}

// usage overrides the default `flag.Usage`.
func usage() {
	fmt.Printf("  To create an app with name 'my-app' run:\n")
	fmt.Printf("\n")
	fmt.Printf("  go run create-go-app.com --http my-app\n")
	fmt.Printf("\n")
	fmt.Printf("  There is only one non-named flag allowed for the name. e.g. 'my-app'\n")
	fmt.Printf("\n")
	flag.PrintDefaults()
}

func chdir(dirname string) {
	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", colors.White, colors.Cyan, colors.Green, dirname, colors.Default)
	err := os.Chdir("./" + dirname)
	if err != nil {
		fatal(ErrChdir)
	}
}

// writeFile may write multiple files depending on requirements for other types of apps.
func writeFile(dirname string) {
	fmt.Printf("%sWriting %smain.go%s file...%s\n", colors.White, colors.Cyan, colors.White, colors.Default)
	err := os.WriteFile(dirname+"/main.go", []byte(code.HTTP), perm.RW)
	fmt.Print("\n")
	if err != nil {
		fatal(ErrWriteFiles)
	}
}

// mkdir makes a new directory with read, write, and execute permissions.
func mkdir(dirname string) {
	fmt.Printf("%sMaking new dir %s./%s%s\n", colors.White, colors.Green, dirname, colors.Default)
	err := os.Mkdir(dirname, perm.RWX)
	if err != nil {
		fatal(ErrMkdir)
	}
}

// checkExists will trying opening an existing directory. If a dir exists (there is no error) then fatally exit.
func checkExists(dirname string) {
	fmt.Printf("Checking if %s./%s%s already exists...%s\n", colors.Green, dirname, colors.White, colors.Default)
	_, err := os.Open("./" + dirname)
	if err == nil {
		fatal(ErrDirExists)
	}
}

// printWkdir prints the new packages working directory.
func printWkdir(dirname string) {
	dir, err := os.Getwd()
	if err != nil {
		fatal(ErrWkdir)
	}

	fmt.Printf("Creating a new %sGo%s app in %s%s/%s\n%s", colors.Cyan, colors.Default, colors.Green, dir, dirname, colors.Default)
}

// checkNamed checks all the named flags and errors if only one is not set.
func checkNamed(flags []*bool) {
	var providedFlags int // If a flag is provided, increment this variable.
	for i := 0; i < len(flags); i++ {
		if *flags[i] {
			providedFlags++
		}
	}
	// Only one type of flag is allowed: cli, http, or module.
	if providedFlags != 1 {
		fatal(ErrNameFlag)
	}
}

// checkNonNamed checks all the non-named flags and errors if only one is not set.
func checkNonNamed(flags []string) {
	if len(flags) != 1 {
		fatal(ErrNonNameFlag)
	}
}

// setupFlags prepares the pointer flags for consumption. The flags setup logic should live here.
func setupFlags() []*bool {
	cli := flag.Bool("cli", false, "Create a CLI app")
	http := flag.Bool("http", false, "Create an HTTP server")
	module := flag.Bool("lib", false, "Create a shareable library")

	return []*bool{cli, http, module}
}

// debugCmdRun helps debug the output from cmd.Run method.
func debugCmdRun(cmd *exec.Cmd) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
}
