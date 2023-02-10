package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"create-go-app.com/cmdFlags"
	"create-go-app.com/directories"
	"create-go-app.com/gotools"
	"create-go-app.com/pkg/colors"
	"create-go-app.com/pkg/timer"
)

//go:embed emit
var content embed.FS

const BaseRepo = "github.com/username/repo"

type App struct {
	dirname string
	flag    string // http, cli, or lib.
}

func main() {
	// Make sure ANSI codes are supported by this terminal.
	colors.CheckTerminal()

	sw := timer.Start()

	// Allocate zero-value app.
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

	err = fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Println(d.Name())
			return nil
		}
		f, err := content.Open(path)
		if err != nil {
			return err
		}
		s, err := f.Stat()
		if err != nil {
			return err
		}
		fmt.Println(s.Name())
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		log.Printf("%q", b)
		// return nil

		if err != nil {
			return err
		}

		dstPath := filepath.Join(app.dirname, strings.TrimPrefix(path, app.dirname))

		if d.IsDir() {
			err := os.MkdirAll(dstPath, os.FileMode(0777))
			if err != nil {
				return err
			}
		} else {
			_, err := cpFile(path, dstPath) // return the written bytes?
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return

	// err = Cp("emit", app.dirname)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = filepath.WalkDir(app.dirname, func(path string, d fs.DirEntry, err error) error {
		return changeGoImports(path, d, BaseRepo, app.dirname)
	})

	if err != nil {
		log.Fatal(err)
	}

	err = directories.Change(app.dirname + "/go")
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrChdir, colors.Default)
		os.Exit(1)
	}

	err = gotools.ChangeModuleName(app.dirname)

	if err != nil {
		log.Fatal(err)
	}

	return

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

	// err = directories.CreateMainFile(app.dirname, app.flag)
	// mainFile, err := os.Create(app.dirname + "/" + emitted.Combined[0].Filename + ".go")
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrCreateFile, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sWriting %smain.go%s file...%s\n", colors.White, colors.Cyan, colors.White, colors.Default)

	fmt.Printf("\n")

	// mainTemp := template.Must(template.New(emitted.Combined[0].Filename).Parse(emitted.Combined[0].Text))
	// err = mainTemp.Execute(mainFile, struct{ Port int }{Port: 9999})

	// Create go package
	// err = os.Mkdir(app.dirname+"/go", os.FileMode(0777))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// goFile, err := os.Create(app.dirname + "/go/" + emittedgo.GoTemplate.Filename + ".go")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// goTemp := template.Must(template.New(emittedgo.GoTemplate.Filename).Parse(emittedgo.GoTemplate.Text))
	// err = goTemp.Execute(goFile, nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = directories.Change(app.dirname)
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrChdir, colors.Default)
		os.Exit(1)
	}
	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", colors.White, colors.Cyan, colors.Green, app.dirname, colors.Default)

	fmt.Printf("\n")

	err = gotools.ValidateModule()
	if err != nil {
		fmt.Printf("%s%v%s\n", colors.Red, ErrInvalidModule, colors.Default)
		os.Exit(1)
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
	fmt.Printf("  Available types: cli or http\n")
	fmt.Printf("\n")
	flag.PrintDefaults()
}

// Used to copy _emitted into the users new directory
func Cp(src, dst string) error {
	return filepath.WalkDir(src, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, strings.TrimPrefix(path, src))

		if info.IsDir() {
			err := os.MkdirAll(dstPath, info.Type())
			if err != nil {
				return err
			}
		} else {
			_, err := cpFile(path, dstPath) // return the written bytes?
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func cpFile(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()
	written, err := io.Copy(dstFile, srcFile)
	return written, err
}

// Change go import paths in all files.
// Will make this smarter to walk other paths in node or go.mod
func changeGoImports(path string, d fs.DirEntry, prev string, new string) error {
	if d.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", d.Name())

	if err != nil {
		return err
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), prev, new, -1)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}

	}

	return nil
}
