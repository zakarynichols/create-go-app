package gotools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InitializeModule(module string) ([]byte, error) {
	cmd := exec.Command("go", "mod", "init", module)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return b, nil
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
		return "", fmt.Errorf("create-go-app: you must provide a module name")
	}

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

func GetAllDeps() error {
	cmd := exec.Command("go", "get", "./...")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
