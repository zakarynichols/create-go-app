package gotools

import (
	"os/exec"
)

func FormatCode() ([]byte, error) {
	cmd := exec.Command("go", "fmt", "./...")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return b, nil
}
