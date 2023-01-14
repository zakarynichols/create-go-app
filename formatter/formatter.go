package formatter

import (
	"os/exec"
)

func FormatCode() error {
	cmd := exec.Command("go", "fmt", "./...")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
