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

	DEPRECATED_colors "create-go-app.dev/colors"
	"create-go-app.dev/fsys"
	"create-go-app.dev/gotools"
	"create-go-app.dev/timer"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
)

// 'create-go-app/embed'
const EMBED_PATH = "embed"

//go:embed all:embed
var emb embed.FS

const exampleRepoURL = "github.com/username/repo"

type app struct {
	appName  string
	fullPath string
	embed    embedded
}

type embedded struct {
	fs embed.FS
}

func (e embedded) Open(name string) (fsys.FileReaderCloser, error) {
	file, err := e.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type fileService struct{}

func (f fileService) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f fileService) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (f fileService) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func (f fileService) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

type fileReaderWriter struct{}

func (f fileReaderWriter) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (f fileReaderWriter) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

var strFlag = flag.String("type", "http", "'http' or 'cli'")

func NewEmbedded(e embed.FS) embedded {
	return embedded{e}
}

func NewApp(embed embed.FS) app {
	e := NewEmbedded(embed)
	return app{embed: e}
}

func main() {
	// Inject embed path.
	fsys.EmbedPath = EMBED_PATH

	a := NewApp(emb)

	// This is to start cleanup when the user tries to exit. When prompted to
	// enter the module name, if the user signals an interrupt, it doesn't
	// stop execution. The app will wait for the user to enter a newline '\n'.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Wait for the app to complete or the app will wait for an interrupt signal.
	done := make(chan error, 1)

	go func() {
		err := run(&a)
		done <- err
	}()

	// Wait for either an interrupt or the app logic to complete.
	select {
	case <-sigChan:
		fmt.Println("\nInterrupt received. Initiating cleanup...")
		err := clean(a.fullPath)
		if err != nil {
			fmt.Println("cleanup failed: ", err)
			os.Exit(1)
		}
		fmt.Println("Cleanup complete, exiting.")
	case err := <-done:
		if err != nil {
			if errors.Is(ErrDirExists, err) {
				fmt.Printf("create-go-app: directory '%s' already exists\n", a.fullPath)
				// TODO: Add (y/n) overwrite prompt.
			} else {
				fmt.Printf("%v\n", err)
				// TODO: Call cleanup when done testing failed output.
			}
		} else {
			fmt.Println("App logic completed successfully.")
		}
	}
}

func run(a *app) error {
	env := os.Getenv("CREATE_GO_APP_ENV")
	fmt.Printf("CREATE_GO_APP_ENV: %s\n", env)

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
		return err
	}

	a.fullPath = filepath.Join(wkdir, a.appName)

	if info, err := os.Stat(a.fullPath); err == nil && info.IsDir() {
		return ErrDirExists
	}

	fmt.Fprintf(color.Output, "Creating a new %s app in %s\n", color.CyanString("Go"), a.fullPath)

	moduleName, err := gotools.EnterModuleName()
	if err != nil {
		return err
	}

	e := embedded{fs: emb}

	// Walk the whole 'emit' directory and dynamically create the directories and files.
	err = fs.WalkDir(e.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ops := fileService{}
		return fsys.Output(a.appName, path, d.IsDir(), e, ops)
	})

	if err != nil {
		return err
	}

	// Now that the directory is created from via fs.WalkDir, walk the newly created dir
	// and update all import paths with the user's provided module string.
	err = filepath.WalkDir(a.appName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ops := fileReaderWriter{}
		return fsys.ReplaceImports("*.go", path, exampleRepoURL, moduleName, d, ops)
	})

	if err != nil {
		return err
	}

	p := filepath.Join(a.appName, "go")

	_, err = os.Stat(p)

	if errors.Is(err, fs.ErrNotExist) {

		if env == "development" {
			return errors.Join(err, fmt.Errorf("create-go-app: did you remove embed/go/go.mod before running the app"))
		}
	}

	err = os.Chdir(p)
	if err != nil {
		return err
	}

	_, err = gotools.InitializeModule(moduleName)
	if err != nil {
		return err
	}

	DEPRECATED_colors.Printf("%sFetching dependencies: %sgo get ./...%s\n", DEPRECATED_colors.White, DEPRECATED_colors.Cyan, DEPRECATED_colors.Default)
	err = gotools.GetAllDeps()
	if err != nil {
		return err
	}

	_, err = gotools.FormatCode()
	if err != nil {
		DEPRECATED_colors.Printf("%s%v%s\n", DEPRECATED_colors.Red, ErrFmt, DEPRECATED_colors.Default)
		return err
	}
	DEPRECATED_colors.Printf("%sFormatting code: %sgo fmt ./...%s\n", DEPRECATED_colors.White, DEPRECATED_colors.Cyan, DEPRECATED_colors.Default)

	// Get the time it took for the program to complete.
	elapsed := start.Elapsed()

	DEPRECATED_colors.Printf("%sSucceeded in %f seconds\n%s", DEPRECATED_colors.Green, elapsed.Seconds(), DEPRECATED_colors.Default)

	return nil
}

func usage() {
	fmt.Printf("  To create an http server with the name 'my-app' run:\n")
	fmt.Printf("  go run create-go-app.dev@latest -http my-app\n")
	fmt.Printf("  The last argument must be the name. e.g. 'my-app'\n")
	fmt.Printf("  Available named flag arguments: cli or http\n")
	flag.PrintDefaults()
}

func clean(path string) error {
	var err error

	DEPRECATED_colors.Printf("%sExecuting cleanup...%s\n", DEPRECATED_colors.Cyan, DEPRECATED_colors.Default)

	_, err = os.Stat(path)
	if err != nil {
		fmt.Printf("Skipping cleanup. Directory '%s' does not exist.\n", path)
		return err
	}

	DEPRECATED_colors.Printf("Attempting to delete directory %s%s%s (y/n): ", DEPRECATED_colors.Red, path, DEPRECATED_colors.Default)

	var input string
	fmt.Scan(&input)

	if input == "y" {
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Printf("Failed to cleanup directory '%s'\n", path)
			return err
		}
		DEPRECATED_colors.Printf("%sCleanup successful.%s Removed directory %s%s%s.\n", DEPRECATED_colors.Green, DEPRECATED_colors.Default, DEPRECATED_colors.Yellow, path, DEPRECATED_colors.Default)

		return nil
	}

	if input == "n" {
		fmt.Printf("Skipping deletion of directory '%s'\n", path)

		return nil
	}

	fmt.Printf("Skipping deletion of directory '%s'\n", path)

	return nil
}
