package flags

import "fmt"

const (
	HttpFlag int = 1
	CliFlag  int = 0
)

type Flag struct {
	Value int
}

func New(httpFlag, cliFlag bool) *Flag {
	if httpFlag {
		return &Flag{Value: HttpFlag}
	} else if cliFlag {
		return &Flag{Value: CliFlag}
	}
	return nil
}

func (f *Flag) Validate() error {
	if f == nil {
		return fmt.Errorf("at least one flag must be set")
	}
	if f.Value != HttpFlag && f.Value != CliFlag {
		return fmt.Errorf("invalid flag value: %d", f.Value)
	}
	return nil
}
