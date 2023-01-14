package cmdFlags

import (
	"errors"
)

var (
	ErrInvalidNonNamed = errors.New("create-go-app: only one non-named flag argument allowed")
	ErrInvalidNamed    = errors.New("create-go-app: only a single named flag can be used to init a package. e.g. --cli, --http, or --module")
)

// ValidateNonNamed function validates the non-named flags passed to the application
func ValidateNonNamed(flags []string) (string, error) {
	if len(flags) != 1 {
		return "", ErrInvalidNonNamed
	}
	return flags[0], nil
}

// ValidateNamed function validates the named flags passed to the application
func ValidateNamed(flagType string) error {
	if flagType != "cli" && flagType != "http" && flagType != "lib" {
		return ErrInvalidNamed
	}
	return nil
}
