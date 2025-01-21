package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"create-go-app.dev/colors"
	"create-go-app.dev/flags"
	"create-go-app.dev/gotools"
	"create-go-app.dev/timer"
)

// TODO: Cleanup app if an error occurs by removing the newly created directory.

//go:embed all:emit
var emitted embed.FS

const BaseRepo = "github.com/username/repo"

type app struct {
	appName  string
	fullPath string
}

var httpFlag = flag.Bool("http", false, "Create an http server")
var cliFlag = flag.Bool("cli", false, "Create a command line interface")

func main() {
	a := new(app)
	err := run(a)
	if err != nil {
		// fmt.Printf("\n")
		// fmt.Print(err)
		// fmt.Printf("\n")
		err = cleanup(a.fullPath)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	}
}

func run(a *app) error {
	// Make sure ANSI codes are supported by this terminal.
	colors.CheckTerminal()

	// Start a timer.
	start := timer.Start()

	// Assign our own custom usage handler.
	flag.Usage = usage

	// Parse the cmd line flags.
	flag.Parse()

	// Flags come before non-flag arguments.
	namedFlag := flags.New(*httpFlag, *cliFlag)

	// Validate 'http' or 'cli' named argument.
	err := namedFlag.Validate()
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrNamedFlag, colors.Default)
		return err
	}

	// Get all non-flag arguments passed to the program.
	nonFlagArgs := flag.Args()

	// Validate the non-named flags. There should only be one.
	if len(nonFlagArgs) != 1 {
		colors.Printf("%s%v%s\n", colors.Red, ErrPosArgs, colors.Default)
		return err
	}

	// Assign the last and only non-flag argument to the apps's root directory name.
	a.appName = nonFlagArgs[0]

	// Get the working directory to show the user where the app is being created.
	wkdir, err := os.Getwd()
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrWkdir, colors.Default)
		return err
	}

	fullPath := filepath.Join(wkdir, a.appName)

	a.fullPath = fullPath

	if info, err := os.Stat(a.fullPath); err == nil && info.IsDir() {
		return colors.Printf("%s Directory %s%s%s already exists%s\n", colors.Red, colors.Yellow, a.fullPath, colors.Red, colors.Default)
		// return errors.New("generic error")
	}

	colors.Printf("Creating a new %sGo%s app in %s%s\n%s", colors.Cyan, colors.Default, colors.Green, fullPath, colors.Default)

	fmt.Printf("\n")

	// Walk the whole 'emit' directory and dynamically create the directories and files.
	err = fs.WalkDir(emitted, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		// Remove the 'emit' string from the path.
		r := strings.Replace(path, "emit", "", -1)

		// Join the new app's directory name to the new string.
		dst := filepath.Join(a.appName, strings.TrimPrefix(r, a.appName))

		// Create directories if they don't exist.
		if d.IsDir() {
			err := os.Mkdir(dst, os.FileMode(0777))
			if os.IsExist(err) {
				return nil
			} else {
				return err
			}
		}

		f, err := emitted.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		err = os.WriteFile(dst, b, os.FileMode(0777))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	moduleName, err := gotools.EnterModuleName()
	if err != nil {
		return err
	}

	// Now that the directory is created from via fs.WalkDir, walk the newly created dir
	// and update all import paths with the user's provided module string.
	err = filepath.WalkDir(a.appName, func(path string, d fs.DirEntry, err error) error {
		return changeGoImports(path, d, BaseRepo, moduleName)
	})

	if err != nil {
		return err
	}

	err = os.Chdir(a.appName + "/go")
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrChdir, colors.Default)
		return err
	}

	_, err = gotools.InitializeModule(moduleName)
	if err != nil {
		return err
	}

	fmt.Printf("\n")

	colors.Printf("%sFetching dependencies: %sgo get ./...%s\n", colors.White, colors.Cyan, colors.Default)
	err = gotools.GetAllDeps()
	if err != nil {
		return err
	}

	fmt.Printf("\n")

	_, err = gotools.FormatCode()
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrFmt, colors.Default)
		return err
	}
	colors.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)

	fmt.Printf("\n")

	// Get the time it took for the program to complete.
	elapsed := start.Elapsed()

	colors.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)

	return nil
}

func usage() {
	fmt.Printf("  To create an http server with the name 'my-app' run:\n")
	fmt.Printf("\n")
	fmt.Printf("  go run create-go-app.dev@latest -http my-app\n")
	fmt.Printf("\n")
	fmt.Printf("  The last argument must be the name. e.g. 'my-app'\n")
	fmt.Printf("\n")
	fmt.Printf("  Available named flag arguments: cli or http\n")
	fmt.Printf("\n")
	flag.PrintDefaults()
}

// Change go import paths in all files.
func changeGoImports(path string, d fs.DirEntry, prev string, new string) error {
	if d.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", d.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), prev, new, -1)

		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

	}

	return nil
}

// TODO: Is this sufficient cleanup?
func cleanup(path string) error {
	var err error

	fmt.Printf("\n")

	colors.Printf("%sAn unexpected error occurred.%s\n", colors.Red, colors.Default)

	fmt.Printf("\n")

	colors.Printf("%sExecuting cleanup...%s\n", colors.Cyan, colors.Default)

	fmt.Printf("\n")

	_, err = os.Stat(path)
	if err != nil {
		fmt.Printf("No cleanup needed. Directory '%s' does not exist.\n", path)
		return err
	}

	colors.Printf("Attempting to delete directory %s%s%s (y/n): ", colors.Red, path, colors.Default)

	var input string
	fmt.Scan(&input)

	fmt.Printf("\n")

	if input == "y" {
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Failed to cleanup directory '%s'\n", path)
			return err
		}
		colors.Printf("%sCleanup successful.%s Removed directory %s%s%s.\n", colors.Green, colors.Default, colors.Yellow, path, colors.Default)

		fmt.Printf("\n")

		return nil
	}

	if input == "n" {
		fmt.Printf("Skipping deletion of directory '%s'\n", path)

		fmt.Printf("\n")

		return nil
	}

	fmt.Printf("Skipping deletion of directory '%s'\n", path)

	fmt.Printf("\n")

	return nil
}
