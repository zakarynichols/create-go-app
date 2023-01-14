package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/zakarynichols/create-go-app/code"
	"github.com/zakarynichols/create-go-app/colors"
	"github.com/zakarynichols/create-go-app/perm"
	"github.com/zakarynichols/create-go-app/print"
	"github.com/zakarynichols/create-go-app/timer"
)

// Errors exposed to the user. Stack traces and more detailed
// errors for debugging will be written to a log file.
var (
	ErrNonNameFlag   = errors.New("create-go-app: only one non-named flag argument allowed")
	ErrNameFlag      = errors.New("create-go-app: only a single named flag can be used to init a package. e.g. --cli, --http, or --module")
	ErrDirExists     = errors.New("create-go-app: directory already exists")
	ErrMkdir         = errors.New("create-go-app: failed to create directory")
	ErrChdir         = errors.New("create-go-app: failed to change directory")
	ErrWkdir         = errors.New("create-go-app: failed to get working directory")
	ErrInitMod       = errors.New("create-go-app: failed to init a module")
	ErrFmt           = errors.New("create-go-app: failed to format code")
	ErrWriteFiles    = errors.New("create-go-app: failed to write files")
	ErrReadModule    = errors.New("create-go-app: failed to read module name")
	ErrEmptyModule   = errors.New("create-go-app: module name cannot be empty")
	ErrLongModule    = errors.New("create-go-app: module name is too long")
	ErrInvalidModule = errors.New("create-go-app: invalid module name")
)

// program is the structure describing the initialized program.
// It should have a directory name and a module name.
type program struct {
	dirname string
	module  string
}

func main() {
	// Make sure ANSI codes are supported by this terminal.
	colors.CheckTerminal()
	// Construct a new timer.
	t := timer.New()
	// Init a new program as a pointer.
	prog := new(program)
	// Setup flags state and usage handler.
	namedFlags := namedFlags() // If using pointers, must be declared before flag.Parse().
	flag.Usage = usage
	flag.Parse()
	// Init and validate flags.
	nonNamedFlags := flag.Args()
	errNonNamed(nonNamedFlags)
	prog.dirname = nonNamedFlags[0]
	errNamed(namedFlags)
	// Core functions that interface with the user directing the flow of the CLI program.
	err := getCurrentWorkingDirectory(prog.dirname)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	checkIfDirectoryExists(prog.dirname)
	createDirectory(prog.dirname)
	createFile(prog.dirname)
	changeDirectory(prog.dirname)
	validateModule(prog)
	formatCode()
	// Get the time it took for the program to complete.
	elapsed := t.Since(t.Start)
	fmt.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)
}

func formatCode() error {
	fmt.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)
	cmd := exec.Command("go", "fmt", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command finished with error: %v, output: %s", err, output)
		return fmt.Errorf("%w: %v", ErrFmt, err)
	}
	fmt.Print("\n")
	return nil
}

func initializeModule(module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, out)
	}
	return nil
}

func validateModule(prog *program) error {
	fmt.Print("Enter the name of the module: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("%w: %v", ErrReadModule, err)
	}
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return fmt.Errorf("%w", ErrEmptyModule)
	}
	if len(input) > 100 {
		return fmt.Errorf("%w", ErrLongModule)
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", input)
	if !match {
		return fmt.Errorf("%w", ErrInvalidModule)
	}
	prog.module = input
	err = initializeModule(prog.module)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInitMod, err)
	}

	fmt.Print("\n")
	return nil
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

func changeDirectory(dirname string) {
	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", colors.White, colors.Cyan, colors.Green, dirname, colors.Default)
	err := os.Chdir("./" + dirname)
	if err != nil {
		fatal(ErrChdir)
	}
	fmt.Print("\n")
}

// createFile may write multiple files depending on requirements for other types of apps.
func createFile(dirname string) {
	fmt.Printf("%sWriting %smain.go%s file...%s\n", colors.White, colors.Cyan, colors.White, colors.Default)
	err := os.WriteFile(dirname+"/main.go", []byte(code.HTTP), perm.RW)
	if err != nil {
		fatal(ErrWriteFiles)
	}
	fmt.Print("\n")
}

// createDirectory makes a new directory with read, write, and execute permissions.
func createDirectory(dirname string) {
	fmt.Printf("%sMaking new dir %s./%s%s\n", colors.White, colors.Green, dirname, colors.Default)
	err := os.Mkdir(dirname, perm.RWX)
	if err != nil {
		fatal(ErrMkdir)
	}
	fmt.Print("\n")
}

// checkIfDirectoryExists will trying opening an existing directory. If a dir exists (there is no error) then fatally exit.
func checkIfDirectoryExists(dirname string) {
	fmt.Printf("Checking if %s./%s%s already exists...%s\n", colors.Green, dirname, colors.White, colors.Default)
	_, err := os.Open("./" + dirname)
	if err == nil {
		fatal(ErrDirExists)
	}
	fmt.Print("\n")
}

func getCurrentWorkingDirectory(dirname string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWkdir, err)
	}
	fmt.Printf("Creating a new %sGo%s app in %s%s/%s\n%s", colors.Cyan, colors.Default, colors.Green, dir, dirname, colors.Default)
	fmt.Print("\n")
	return nil
}

// errNamed checks all the named flags and errors if only one is not set.
func errNamed(flags []*bool) {
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

// errNonNamed checks all the non-named flags and errors if only one is not set.
func errNonNamed(flags []string) {
	if len(flags) != 1 {
		fatal(ErrNonNameFlag)
	}
}

// namedFlags prepares the pointer flags for consumption. The flags setup logic should live here.
func namedFlags() []*bool {
	cli := flag.Bool("cli", false, "Create a CLI app")
	http := flag.Bool("http", false, "Create an HTTP server")
	module := flag.Bool("lib", false, "Create a shareable library")

	return []*bool{cli, http, module}
}
