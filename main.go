package main

import (
	"flag"
	"fmt"
	"os"

	"create-go-app.com/cmdFlags"
	"create-go-app.com/colors"
	"create-go-app.com/directories"
	"create-go-app.com/formatter"
	"create-go-app.com/modules"
	"create-go-app.com/stopwatch"
)

type App struct {
	dirname string
	flag    string // HTTP, CLI, or library.
}

func main() {
	// Make sure ANSI codes are supported by this terminal.
	colors.CheckTerminal()

	sw := stopwatch.Start()

	// Init a new program as a pointer.
	app := new(App)

	// Setup flags state and usage handler.
	flagType := flag.String("type", "", "Type of project to create. Options are: cli, http, lib")

	flag.Usage = usage

	flag.Parse()

	app.flag = *flagType

	// Get all non-named flags passed to the program.
	nonNamedFlags := flag.Args()

	// Validate the non-named flags. There should only be one.
	validNonNamedFlag, err := cmdFlags.ValidateNonNamed(nonNamedFlags)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrNonNamedFlag, colors.Default)
		os.Exit(1)
	}

	// Assign the last and only non-named flag to the apps's root directory name.
	app.dirname = validNonNamedFlag

	err = cmdFlags.ValidateNamed(app.flag)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrNamedFlag, colors.Default)
		os.Exit(1)
	}

	wkdir, err := directories.GetWorkingDirectory()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrWkdir, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("Creating a new %sGo%s app in %s%s/%s\n%s", colors.Cyan, colors.Default, colors.Green, wkdir, app.dirname, colors.Default)

	fmt.Printf("\n")

	err = directories.Exists(app.dirname)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrDirExists, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("Checking if %s./%s%s already exists...%s\n", colors.Green, app.dirname, colors.White, colors.Default)

	fmt.Printf("\n")

	err = directories.Create(app.dirname)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrCreateDir, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sMaking new dir %s./%s%s\n", colors.White, colors.Green, app.dirname, colors.Default)

	fmt.Printf("\n")

	err = directories.CreateFile(app.dirname, app.flag)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrCreateFile, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sWriting %smain.go%s file...%s\n", colors.White, colors.Cyan, colors.White, colors.Default)

	fmt.Printf("\n")

	err = directories.Change(app.dirname)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrChdir, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", colors.White, colors.Cyan, colors.Green, app.dirname, colors.Default)

	fmt.Printf("\n")

	err = modules.ValidateModule()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrInvalidModule, colors.Default)
		os.Exit(1)
	}

	fmt.Printf("\n")

	err = formatter.FormatCode()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrFmt, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)

	fmt.Printf("\n")

	// Get the time it took for the program to complete.
	elapsed := sw.Elapsed()

	fmt.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)
}

func usage() {
	fmt.Printf("  To create an http server with the name 'my-app' run:\n")
	fmt.Printf("\n")
	fmt.Printf("  go run create-go-app.com@latest -type http my-app\n")
	fmt.Printf("\n")
	fmt.Printf("  The last argument must be the name. e.g. 'my-app'\n")
	fmt.Printf("\n")
	fmt.Printf("  Available types: cli, http, lib\n")
	fmt.Printf("\n")
	flag.PrintDefaults()
}
