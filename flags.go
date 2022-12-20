package flags

import (
	"flag"
	"fmt"
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
	c.Set("You ran without a name flag")
	flag.Var(c, "name", "help message for name")
	flag.Parse()

	fmt.Println(c.text)
	fmt.Println(c.String())
}
