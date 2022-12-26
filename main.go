package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/zakarynichols/create-go-app/colors"
)

type pkg struct {
	name string
}

func main() {
	var err error

	start := time.Now()

	col := colors.New()

	pkg := new(pkg)

	flag.Parse()

	flags := flag.Args()

	if len(flags) != 1 {
		fmt.Printf("%serror: only one non-named flag argument allowed.%s\n", col.Red, col.Reset)
		os.Exit(1)
	}

	pkg.name = flags[0] // Will make this 'smarter' with help message and such.

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%sCreating a new %sGo%s app in %s%s/%s%s %s\n", col.White, col.Cyan, col.White, col.Green, dir, pkg.name, col.Cyan, col.Reset)

	fmt.Print("\n")

	fmt.Printf("%sChecking if %s./%s%s already exists...%s\n", col.White, col.Green, pkg.name, col.White, col.Reset)
	_, err = os.Open("./" + pkg.name)
	if err == nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: directory "+"'./"+"%s"+"'"+" already exists%s\n", col.Red, pkg.name, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sMaking new dir %s./%s%s\n", col.White, col.Green, pkg.name, col.Reset)
	err = os.Mkdir(pkg.name, 0750)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to create directory\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sWriting %smain.go%s file...%s\n", col.White, col.Cyan, col.White, col.Reset)
	err = os.WriteFile(pkg.name+"/main.go", []byte(mainTemplate), 0660)
	fmt.Print("\n")
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to write files\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Printf("%sChanging to dir: %scd %s./%s%s\n", col.White, col.Cyan, col.Green, pkg.name, col.Reset)
	err = os.Chdir("./" + pkg.name)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to change directory\n%s", col.Red, col.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sInitializing a module: %sgo mod init %s%s\n", col.White, col.Cyan, pkg.name, col.Reset)
	cmd := exec.Command("go", "mod", "init", pkg.name) // don't hardcode module name. only for testing
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

	fmt.Printf("%sSucceeded in %f seconds\n", col.Green, elapsed.Seconds())
}

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
