package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zakarynichols/create-go-app/colors"
	"github.com/zakarynichols/create-go-app/dir"
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

func main() {
	start := time.Now()

	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Only one cmd flag is permitted, the module name")
	}

	pkg := *new(pkg)

	pkg.set(flag.Args()[0])

	err := os.Mkdir(pkg.string(), 0750)
	if err != nil {
		log.Fatal(err)
	}

	err = dir.CpdirAll("template", pkg.string())
	if err != nil {
		log.Fatal(err)
	}

	since := time.Since(start)

	fmt.Print("\n")

	fmt.Printf("%sSucceeded in %f seconds\n", colors.Green, since.Seconds())
}

// func main() {
// 	start := time.Now()

// 	color := colors.New()

// 	pkgName := new(Name)

// 	flag.Parse()

// 	if len(flag.Args()) != 1 {
// 		fmt.Printf("%serror: only one non-named flag argument allowed.%s\n", colors.Red, colors.Reset)
// 		os.Exit(1)
// 	}

// 	name := flag.Args()[0]

// 	pkgName.text = name

// 	fmt.Printf("%sCreating a new Go app in %s'./%s'%s %s\n", color.Cyan, color.Green, pkgName.text, color.Cyan, color.Reset)

// 	fmt.Print("\n")

// 	fmt.Printf("%sChecking if %s'./%s'%s already exists...%s\n", color.Cyan, color.Green, pkgName.text, color.Cyan, color.Reset)
// 	_, err = os.Open("./" + pkgName.text)
// 	if err == nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: directory "+"'./"+"%s"+"'"+" already exists%s\n", color.Red, pkgName.text, color.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	fmt.Print("\n")

// 	fmt.Printf("%sMaking new dir %s'./%s'%s\n", color.Cyan, color.Green, pkgName.text, color.Reset)
// 	err = os.Mkdir(pkgName.text, 0750)
// 	if err != nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: failed to create directory\n%s", colors.Red, colors.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	fmt.Print("\n")

// 	fmt.Printf("Writing ./main.go file...\n")
// 	err = os.WriteFile(pkgName.text+"/main.go", []byte(mainTemplate), 0660)
// 	fmt.Print("\n")
// 	if err != nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: failed to write files\n%s", colors.Red, colors.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	fmt.Printf("Changing to dir ./%s\n", pkgName.text)
// 	err = os.Chdir("./" + pkgName.text)
// 	if err != nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: failed to change directory\n%s", colors.Red, colors.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	fmt.Print("\n")

// 	fmt.Printf("Executing go mod init...\n")
// 	cmd := exec.Command("go", "mod", "init", "example/my-module") // don't hardcode module name. only for testing
// 	err = cmd.Run()
// 	if err != nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: failed to initialize a module\n%s", colors.Red, colors.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	fmt.Print("\n")

// 	fmt.Printf("Executing go fmt...\n")
// 	cmd = exec.Command("go", "fmt", "./...")
// 	err = cmd.Run()
// 	if err != nil {
// 		fmt.Printf("\n")
// 		fmt.Printf("%serror: failed to format code\n%s", colors.Red, colors.Reset)
// 		fmt.Print("\n")
// 		os.Exit(1)
// 	}

// 	t := time.Now()
// 	elapsed := t.Sub(start)

// 	fmt.Print("\n")

// 	fmt.Printf("%sSucceeded in %f seconds\n", colors.Green, elapsed.Seconds())
// }
