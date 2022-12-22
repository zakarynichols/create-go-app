package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zakarynichols/create-go-app/colors"
)

type Name struct {
	text string
}

func (c *Name) String() string {
	return fmt.Sprintf("String() - %s", strings.ToTitle(c.text))
}

func (c *Name) Set(str string) error {
	c.text = str
	return nil
}

func main() {
	var err error

	color := colors.New()

	pkgName := new(Name)

	flag.Var(pkgName, "name", "The name of the package")

	flag.Parse()

	fmt.Printf("%sCreating a new Go app in %s'./%s'%s %s\n", color.Cyan, color.Green, pkgName.text, color.Cyan, color.Reset)

	fmt.Print("\n")

	fmt.Printf("%sChecking if %s'./%s'%s already exists...%s\n", color.Cyan, color.Green, pkgName.text, color.Cyan, color.Reset)
	_, err = os.Open("./" + pkgName.text)
	if err == nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: directory "+"'./"+"%s"+"'"+" already exists%s\n", color.Red, pkgName.text, color.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("%sMaking new dir %s'./%s'%s\n", color.Cyan, color.Green, pkgName.text, color.Reset)
	err = os.Mkdir(pkgName.text, 0750)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to create directory\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("Writing ./main.go file...\n")
	err = os.WriteFile(pkgName.text+"/main.go", []byte(mainTemplate), 0660)
	fmt.Print("\n")
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to write files\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Printf("Changing to dir ./%s\n", pkgName.text)
	err = os.Chdir("./" + pkgName.text)
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to change directory\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("Executing go mod init...\n")
	cmd := exec.Command("go", "mod", "init", "example/my-module") // don't hardcode module name. only for testing
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to initialize a module\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}

	fmt.Print("\n")

	fmt.Printf("Executing go fmt...\n")
	cmd = exec.Command("go", "fmt", "./...")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("\n")
		fmt.Printf("%serror: failed to format code\n%s", colors.Red, colors.Reset)
		fmt.Print("\n")
		os.Exit(1)
	}
}

var mainTemplate = `
package main

import (
"fmt"
"net/http"
)

func main() {
fmt.Println("Create Go App")
http.ListenAndServe(":0101", nil)
}
`
