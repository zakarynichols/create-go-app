package flags

import (
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
