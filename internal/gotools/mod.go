package gotools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InitializeModule(module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func EnterModuleName() (string, error) {
	fmt.Print("Enter the name of the module: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "", err
	}
	// TODO: Restrict the length of the input?
	// if len(input) > 100 {
	// 	return "", err
	// }

	return input, nil
}

func ChangeModuleName(name string) error {
	cmd := exec.Command("go", "mod", "edit", "-module", name)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
