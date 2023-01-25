package code

// package for a basic cli program
// go run main.go -string=Hello -bool=true -int=10 -custom=myCustomFlag
const CliMain string = `
package main

import (
	"flag"
	"fmt"
)

// custom flag type that satisfies the Value interface
type customFlag struct {
	value string
}

func (f *customFlag) String() string {
	return f.value
}

func (f *customFlag) Set(value string) error {
	f.value = value
	return nil
}

func main() {
	var (
		stringFlag = flag.String("string", "default", "a string flag")
		boolFlag   = flag.Bool("bool", false, "a boolean flag")
		intFlag    = flag.Int("int", 0, "an integer flag")
	)

	// define a custom flag and couple it to a flag variable
	var custom customFlag
	flag.Var(&custom, "custom", "a custom flag")

	flag.Parse()

	fmt.Printf("string flag value: %s\n", *stringFlag)
	fmt.Printf("bool flag value: %t\n", *boolFlag)
	fmt.Printf("int flag value: %d\n", *intFlag)
	fmt.Printf("custom flag value: %s\n", custom.String())
}
`
