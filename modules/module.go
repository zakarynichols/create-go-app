package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func initializeModule(module string) error {
	cmd := exec.Command("go", "mod", "init", module)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func ValidateModule() error {
	fmt.Print("Enter the name of the module: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return err
	}
	if len(input) > 100 {
		return err
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", input)
	if !match {
		return err
	}
	err = initializeModule(input)
	if err != nil {
		return err
	}
	return nil
}
