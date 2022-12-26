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

func (pkg *pkg) string() string {
	return pkg.name
}

func (pkg *pkg) set(str string) error {
	pkg.name = str
	return nil
}

// func main() {
// 	start := time.Now()

// 	flag.Parse()

// 	if len(flag.Args()) != 1 {
// 		log.Fatal("Only one cmd flag is permitted, the module name")
// 	}

// 	pkg := new(pkg)

// 	pkg.set(flag.Args()[0])

// 	err := os.Mkdir(pkg.string(), 0750)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	since := time.Since(start)

// 	fmt.Print("\n")

// 	fmt.Printf("%sSucceeded in %f seconds\n", colors.Green, since.Seconds())
// }

func main() {
	var err error

	start := time.Now()

	color := colors.New()

	pkg := new(pkg)

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("%serror: only one non-named flag argument allowed.%s\n", colors.Red, colors.Reset)
		os.Exit(1)
	}

	name := flag.Args()[0]

	err = pkg.set(name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%sCreating a new Go app %s%s%s %s\n", color.Cyan, color.Green, pkg.name, color.Cyan, color.Reset)

	fmt.Print("\n")

	fmt.Printf("%sChecking if %s./%s%s already exists...%s\n", color.Cyan, color.Green, pkg.name, color.Cyan, color.Reset)
	_, err = os.Open("./" + pkg.name)
	if err == nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: directory "+"'./"+"%s"+"'"+" already exists%s\n", color.Red, pkg.name, color.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sMaking new dir %s./%s%s\n", color.Cyan, color.Green, pkg.name, color.Reset)
	err = os.Mkdir(pkg.name, 0750)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to create directory\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sWriting %smain.go%s file...%s\n", color.Cyan, color.Green, color.Cyan, color.Reset)
	err = os.WriteFile(pkg.name+"/main.go", []byte(mainTemplate), 0660)
	fmt.Print("\n")
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to write files\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Printf("%sChanging to dir ./%s%s\n", color.Cyan, pkg.name, color.Reset)
	err = os.Chdir("./" + pkg.name)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to change directory\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sExecuting go mod init...%s\n", color.Cyan, color.Reset)
	cmd := exec.Command("go", "mod", "init", "example/my-module") // don't hardcode module name. only for testing
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to initialize a module\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sExecuting go fmt...%s\n", color.Cyan, color.Reset)
	cmd = exec.Command("go", "fmt", "./...")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to format code\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Print("\n")

	fmt.Printf("%sSucceeded in %f seconds\n", colors.Green, elapsed.Seconds())
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
