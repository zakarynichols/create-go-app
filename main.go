package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"create-go-app.com/colors"
	"create-go-app.com/internal/flags"
	"create-go-app.com/internal/gotools"
	"create-go-app.com/pkg/timer"
)

//go:embed emit
var emitted embed.FS

const BaseRepo = "github.com/username/repo"

type App struct {
	dirname string
	flag    string // http or cli
}

var httpFlag = flag.Bool("http", false, "Create an http server")
var cliFlag = flag.Bool("cli", false, "Create a command line interface")

func main() {
	// Make sure ANSI codes are supported by this terminal.
	colors.CheckTerminal()

	// Start a timer.
	t := timer.Start()

	// Allocate zero-value app.
	app := new(App)

	// Assign our own custom usage handler.
	flag.Usage = usage

	// Parse the cmd line flags.
	flag.Parse()

	// Flags come before positional arguments.
	namedFlags := flags.NewFlags(*httpFlag, *cliFlag)

	// Validate 'http' or 'cli' named argument.
	flagType, err := namedFlags.Validate()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrNamedFlag, colors.Default)
		os.Exit(1)
	}

	// Set the created app type.
	app.flag = flagType

	// Get all positional arguments passed to the program.
	posArgs := flag.Args()

	// Validate the non-named flags. There should only be one.
	if len(posArgs) != 1 {
		fmt.Printf("%s%v%s\n", colors.Red, ErrPosArgs, colors.Default)
		os.Exit(1)
	}

	// Assign the last and only positional argument to the apps's root directory name.
	app.dirname = posArgs[0]

	// Get the working directory to show the user where the app is being created.
	wkdir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrWkdir, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("Creating a new %sGo%s app in %s%s/%s\n%s", colors.Cyan, colors.Default, colors.Green, wkdir, app.dirname, colors.Default)

	fmt.Printf("\n")

	// Walk the whole 'emit' directory and dynamically create the directories and files.
	err = fs.WalkDir(emitted, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Remove the 'emit' string from the path.
		r := strings.Replace(path, "emit", "", -1)

		// Join the new app's directory name to the new string.
		dst := filepath.Join(app.dirname, strings.TrimPrefix(r, app.dirname))

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
		log.Fatal(err)
	}

	moduleName, err := gotools.EnterModuleName()
	if err != nil {
		log.Fatal(err)
	}

	// Now that the directory is created from via fs.WalkDir, walk the newly created dir
	// and update all import paths with the user's provided module string.
	err = filepath.WalkDir(app.dirname, func(path string, d fs.DirEntry, err error) error {
		return changeGoImports(path, d, BaseRepo, moduleName)
	})

	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(app.dirname + "/go")
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrChdir, colors.Default)
		os.Exit(1)
	}

	err = gotools.InitializeModule(moduleName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")

	_, err = gotools.FormatCode()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrFmt, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sFormatting code: %sgo fmt ./...%s\n", colors.White, colors.Cyan, colors.Default)

	fmt.Printf("\n")

	// Get the time it took for the program to complete.
	elapsed := t.Elapsed()

	fmt.Printf("%sSucceeded in %f seconds\n%s", colors.Green, elapsed.Seconds(), colors.Default)
}

func usage() {
	fmt.Printf("  To create an http server with the name 'my-app' run:\n")
	fmt.Printf("\n")
	fmt.Printf("  go run create-go-app.com@latest -http my-app\n")
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
