package flags

import "fmt"

const (
	HTTPkey string = "http"
	CLIkey  string = "cli"
)

type Flag struct {
	flagType string // http or cli
	isSet    bool   // true is flag set
}

type Flags map[int]Flag

func NewFlags(httpFlag, cliFlag bool) Flags {
	return Flags{
		0: Flag{
			flagType: HTTPkey,
			isSet:    httpFlag,
		},
		1: Flag{
			flagType: CLIkey,
			isSet:    cliFlag,
		},
	}
}

func (f Flags) Validate() (string, error) {
	var count int
	var flag string
	for _, v := range f {
		if v.isSet {
			flag = v.flagType
			count++
		}
	}
	if count != 1 {
		return "", fmt.Errorf("validate: there can only be a single named flag")
	}

	return flag, nil
}
