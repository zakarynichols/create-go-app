package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"create-go-app.dev/colors"
	"create-go-app.dev/gotools"
	"create-go-app.dev/timer"
)

// TODO: Ensure there are no hard-coded path separaters. Use "path/filepath".
// TODO: Cleanup app if an error occurs by removing the newly created directory.

//go:embed all:emit
var emitted embed.FS

const exampleRepoURL = "github.com/username/repo"

type app struct {
	appName  string
	fullPath string
}

var strFlag = flag.String("type", "http", "'http' or 'cli'")

func main() {
	a := new(app)

	// This is to start cleanup when the user tries to exit. When prompted to
	// enter the module name, if the user signals an interrupt, it doesn't
	// stop execution the app will wait for the user to enter a newline '\n'.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Wait for the app to complete or the app will wait for an interrupt signal.
	done := make(chan error, 1)

	go func() {
		err := run(a)
		done <- err
	}()

	// Wait for either an interrupt or the app logic to complete.
	select {
	case <-sigChan:
		fmt.Println("\nInterrupt received. Initiating cleanup...")
		err := cleanup(a.fullPath)
		if err != nil {
			fmt.Println("cleanup failed:", err)
			os.Exit(1)
		}
		fmt.Println("Cleanup complete, exiting.")
	case err := <-done:
		if err != nil {
			if errors.Is(ErrDirExists, err) {
				fmt.Printf("create-go-app: directory '%s' already exists\n", a.fullPath)
				// TODO: Add (y/n) overwrite prompt.
			} else {
				fmt.Printf("fatal error: %v\n", err)
			}
		} else {
			fmt.Println("App logic completed successfully.")
		}
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
	fmt.Printf("flags: %s\n", *strFlag)

	// Get all non-flag arguments passed to the program.
	nonFlagArgs := flag.Args()

	// Validate the non-named flags. There should only be one.
	if len(nonFlagArgs) != 1 {
		return ErrNonFlagArgs
	}

	// Assign the last and only non-flag argument to the apps's root directory name.
	a.appName = nonFlagArgs[0]

	// Get the working directory to show the user where the app is being created.
	wkdir, err := os.Getwd()
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrWkdir, colors.Default)
		return err
	}

	a.fullPath = filepath.Join(wkdir, a.appName)

	if info, err := os.Stat(a.fullPath); err == nil && info.IsDir() {
		return ErrDirExists
	}

	colors.Printf("Creating a new %sGo%s app in %s%s\n%s", colors.Cyan, colors.Default, colors.Green, a.fullPath, colors.Default)

	moduleName, err := gotools.EnterModuleName()
	if err != nil {
		return err
	}

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

	// Now that the directory is created from via fs.WalkDir, walk the newly created dir
	// and update all import paths with the user's provided module string.
	err = filepath.WalkDir(a.appName, func(path string, d fs.DirEntry, err error) error {
		return changeGoImports(path, d, exampleRepoURL, moduleName)
	})

	if err != nil {
		return err
	}

	p := filepath.Join(a.appName, "go")

	err = os.Chdir(p)
	if err != nil {
		return err
	}

	_, err = gotools.InitializeModule(moduleName)
	if err != nil {
		return err
	}

	colors.Printf("%sFetching dependencies: %sgo get ./...%s\n", colors.White, colors.Cyan, colors.Default)
	err = gotools.GetAllDeps()
	if err != nil {
		return err
	}

	_, err = gotools.FormatCode()
	if err != nil {
		colors.Printf("%s%v%s\n", colors.Red, ErrFmt, colors.Default)
		return err
	}
	colors.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)

	// Get the time it took for the program to complete.
	elapsed := start.Elapsed()

	colors.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)

	return nil
}

func usage() {
	fmt.Printf("  To create an http server with the name 'my-app' run:\n")
	fmt.Printf("  go run create-go-app.dev@latest -http my-app\n")
	fmt.Printf("  The last argument must be the name. e.g. 'my-app'\n")
	fmt.Printf("  Available named flag arguments: cli or http\n")
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

func cleanup(path string) error {
	var err error

	colors.Printf("%sExecuting cleanup...%s\n", colors.Cyan, colors.Default)

	_, err = os.Stat(path)
	if err != nil {
		fmt.Printf("Skipping cleanup. Directory '%s' does not exist.\n", path)
		return err
	}

	colors.Printf("Attempting to delete directory %s%s%s (y/n): ", colors.Red, path, colors.Default)

	var input string
	fmt.Scan(&input)

	if input == "y" {
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Failed to cleanup directory '%s'\n", path)
			return err
		}
		colors.Printf("%sCleanup successful.%s Removed directory %s%s%s.\n", colors.Green, colors.Default, colors.Yellow, path, colors.Default)

		return nil
	}

	if input == "n" {
		fmt.Printf("Skipping deletion of directory '%s'\n", path)

		return nil
	}

	fmt.Printf("Skipping deletion of directory '%s'\n", path)

	return nil
}
