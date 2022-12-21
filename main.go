package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Name struct {
	text string
}

func (c *Name) String() string {
	return fmt.Sprintf("String() - %s", strings.ToTitle(c.text))
}

// the string coming in is the input from cmd line
func (c *Name) Set(str string) error {
	c.text = str
	return nil
}

func main() {
	c := new(Name)
	// c.Set("You ran without a name flag")
	flag.Var(c, "name", "help message for name")
	flag.Parse()

	mkdir(c.text)
	writefile(c.text)

	cmd := exec.Command("go", "fmt", "./"+c.text+"/./...")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir("./" + c.text)
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("go", "mod", "init", "example/my-module")

	err = cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func mkdir(dirname string) {
	err := os.Mkdir(dirname, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func writefile(dirname string) {
	err := os.WriteFile(
		dirname+"/main.go",
		[]byte(`
package main

import (
"fmt"
"net/http"
)

func main() {
fmt.Println("Create Go App")
http.ListenAndServe(":1337", nil)
}
	`),
		0660,
	)
	if err != nil {
		log.Fatal(err)
	}
}
